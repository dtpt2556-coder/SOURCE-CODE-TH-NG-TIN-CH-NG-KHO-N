// Package ingest điều phối content pipeline: khám phá, tải, bóc tách, khử
// trùng, gắn mã, phân loại, tóm tắt, kiểm chứng số liệu, sinh short link.
package ingest

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
	"sync"
	"time"

	"github.com/tenpoint/tenpoint-api/internal/domain"
	"github.com/tenpoint/tenpoint-api/internal/ingest/classifier"
	"github.com/tenpoint/tenpoint-api/internal/ingest/dedupe"
	"github.com/tenpoint/tenpoint-api/internal/ingest/extractor"
	"github.com/tenpoint/tenpoint-api/internal/ingest/fetcher"
	"github.com/tenpoint/tenpoint-api/internal/ingest/summarizer"
	"github.com/tenpoint/tenpoint-api/internal/ingest/tagger"
	"github.com/tenpoint/tenpoint-api/internal/shortlink"
	"github.com/tenpoint/tenpoint-api/internal/store"
	"github.com/tenpoint/tenpoint-api/internal/textutil"
)

// ErrAlreadyRunning được trả khi có pipeline khác đang chạy (advisory lock).
var ErrAlreadyRunning = errors.New("another pipeline run is in progress")

// DefaultWorkers — tối đa 4 worker song song theo nguồn.
const DefaultWorkers = 4

// ClusterWindow — cửa sổ tìm cụm trùng 72 giờ (ADR-009).
const ClusterWindow = 72 * time.Hour

// shortCodeAttempts — số lần thử sinh mã base62 trước khi bỏ cuộc (L-04).
const shortCodeAttempts = 5

// Config là cấu hình runtime của pipeline.
type Config struct {
	UserAgent            string
	FetchTimeout         time.Duration
	MaxArticlesPerSource int
	Workers              int
	Location             *time.Location
	IndustryTier         bool
}

// Result là tổng hợp một lần chạy.
type Result struct {
	RunID     int64
	Found     int
	New       int
	Updated   int
	Unchanged int
	Rejected  int
	Errors    int
	Duration  time.Duration
}

// Pipeline gom mọi dependency qua struct — không có biến global mutable.
type Pipeline struct {
	cfg        Config
	store      *store.Store
	fetcher    *fetcher.HTTPFetcher
	robots     *fetcher.RobotsCache
	extractor  *extractor.Extractor
	classifier *classifier.Classifier
	summarizer summarizer.Provider
	extractive summarizer.Provider
	validator  *summarizer.Validator
	log        *slog.Logger
	now        func() time.Time
}

// New dựng pipeline. `provider` quyết định chất lượng tóm tắt; `extractive` luôn
// có mặt làm phương án dự phòng khi LLM lỗi.
func New(cfg Config, st *store.Store, provider summarizer.Provider, log *slog.Logger) *Pipeline {
	if log == nil {
		log = slog.Default()
	}
	if cfg.Workers <= 0 {
		cfg.Workers = DefaultWorkers
	}
	if cfg.MaxArticlesPerSource <= 0 {
		cfg.MaxArticlesPerSource = 40
	}
	if cfg.FetchTimeout <= 0 {
		cfg.FetchTimeout = 20 * time.Second
	}
	if cfg.Location == nil {
		cfg.Location = time.UTC
	}

	httpFetcher := fetcher.NewHTTPFetcher(cfg.UserAgent, cfg.FetchTimeout, fetcher.NewRateLimiter(time.Second))
	return &Pipeline{
		cfg:        cfg,
		store:      st,
		fetcher:    httpFetcher,
		robots:     fetcher.NewRobotsCache(httpFetcher.GetString),
		extractor:  extractor.New(cfg.Location),
		classifier: classifier.New(),
		summarizer: provider,
		extractive: summarizer.NewExtractive(),
		validator:  summarizer.NewValidator(),
		log:        log,
		now:        time.Now,
	}
}

// Run chạy đủ một lượt pipeline. Chống chạy chồng bằng pg_try_advisory_lock.
func (p *Pipeline) Run(ctx context.Context, trigger string) (Result, error) {
	start := p.now()

	lock, acquired, err := p.store.TryAdvisoryLock(ctx, store.PipelineLockKey)
	if err != nil {
		return Result{}, fmt.Errorf("acquire advisory lock: %w", err)
	}
	if !acquired {
		return Result{}, ErrAlreadyRunning
	}
	defer func() {
		if err := lock.Release(context.WithoutCancel(ctx)); err != nil {
			p.log.Error("release advisory lock", "err", err)
		}
	}()

	runID, err := p.store.StartRun(ctx, trigger)
	if err != nil {
		return Result{}, err
	}
	p.log.Info("crawl run started", "run_id", runID, "trigger", trigger)

	sources, err := p.store.ListEnabledSources(ctx)
	if err != nil {
		return Result{RunID: runID}, err
	}
	tickers, err := p.store.ListTickersForTagger(ctx)
	if err != nil {
		return Result{RunID: runID}, err
	}
	domains, err := p.store.ListSourceDomains(ctx)
	if err != nil {
		return Result{RunID: runID}, err
	}

	newTagger := tagger.New
	if p.cfg.IndustryTier {
		newTagger = tagger.NewWithIndustryTier
	}
	deps := ingestDeps{
		tagger:    newTagger(toTaggerTickers(tickers)),
		allowlist: shortlink.NewAllowlist(domains),
	}

	run := Result{RunID: runID}
	var mu sync.Mutex
	var wg sync.WaitGroup
	slots := make(chan struct{}, p.cfg.Workers)

	for _, src := range sources {
		wg.Add(1)
		slots <- struct{}{}
		go func() {
			defer wg.Done()
			defer func() { <-slots }()
			stat := p.crawlSourceSafely(ctx, runID, src, deps)

			mu.Lock()
			run.Found += stat.found
			run.New += stat.created
			run.Updated += stat.updated
			run.Unchanged += stat.unchanged
			run.Rejected += stat.rejected
			if stat.errMessage != "" {
				run.Errors++
			}
			mu.Unlock()

			if err := p.store.RecordRunSource(ctx, domain.CrawlRunSource{
				RunID: runID, SourceID: src.ID, Found: stat.found,
				NewCount: stat.created, ErrorMessage: stat.errMessage,
				Truncated: stat.truncated, OffTopic: stat.offTopic,
			}); err != nil {
				p.log.Error("record crawl_run_sources", "source", src.Code, "err", err)
			}
		}()
	}
	wg.Wait()

	run.Duration = p.now().Sub(start)
	if err := p.store.FinishRun(ctx, domain.CrawlRun{
		ID: runID, ArticlesFound: run.Found, ArticlesNew: run.New,
		ArticlesRejected: run.Rejected, Errors: run.Errors,
	}); err != nil {
		return run, err
	}
	p.log.Info("crawl run finished", "run_id", runID, "found", run.Found,
		"new", run.New, "updated", run.Updated, "unchanged", run.Unchanged,
		"rejected", run.Rejected, "errors", run.Errors,
		"duration", run.Duration.String())
	return run, nil
}

type ingestDeps struct {
	tagger    *tagger.Tagger
	allowlist *shortlink.Allowlist
}

type sourceStat struct {
	found      int
	created    int
	updated    int
	unchanged  int
	rejected   int
	offTopic   int
	truncated  bool
	errMessage string
}

// crawlSourceSafely đảm bảo một nguồn chết không kéo theo cả run (I-01).
func (p *Pipeline) crawlSourceSafely(ctx context.Context, runID int64, src domain.Source, deps ingestDeps) (stat sourceStat) {
	defer func() {
		if r := recover(); r != nil {
			p.log.Error("panic while crawling source",
				"source", src.Code, "panic", r, "stack", string(debug.Stack()))
			stat.errMessage = fmt.Sprintf("panic: %v", r)
		}
	}()
	return p.crawlSource(ctx, src, deps)
}

func (p *Pipeline) crawlSource(ctx context.Context, src domain.Source, deps ingestDeps) sourceStat {
	var stat sourceStat
	log := p.log.With("source", src.Code)

	if src.RSSURL == "" {
		stat.errMessage = "source has no RSS feed configured"
		return stat
	}

	interval := time.Duration(src.RateLimitMS) * time.Millisecond
	allowed, crawlDelay, err := p.robots.Allowed(ctx, src.RSSURL, p.cfg.UserAgent)
	if err != nil {
		log.Warn("robots.txt unavailable, continuing with default rate limit", "err", err)
	}
	if !allowed {
		stat.errMessage = "skipped_by_robots: feed path disallowed"
		log.Warn("skipped by robots.txt", "url", src.RSSURL)
		return stat
	}
	if crawlDelay > interval {
		interval = crawlDelay
	}

	page, err := p.fetcher.Get(ctx, src.RSSURL, interval)
	if err != nil {
		stat.errMessage = fmt.Sprintf("fetch feed: %v", err)
		log.Error("fetch feed failed", "err", err)
		return stat
	}

	items, truncated, err := fetcher.ParseFeed(ctx, page.Body, p.now(), fetcher.DiscoverWindow, p.cfg.MaxArticlesPerSource)
	if err != nil {
		stat.errMessage = fmt.Sprintf("parse feed: %v", err)
		log.Error("parse feed failed", "err", err)
		return stat
	}
	stat.found = len(items)
	stat.truncated = truncated
	if truncated {
		log.Warn("feed truncated by MAX_ARTICLES_PER_SOURCE, older items skipped", "limit", p.cfg.MaxArticlesPerSource)
	}
	if streak, err := p.store.SourceHealth(ctx, src.ID, len(items)); err != nil {
		log.Warn("update source health", "err", err)
	} else if streak >= 2 {
		log.Error("feed returned zero items on consecutive runs, source layout may have changed",
			"consecutive_empty_runs", streak)
	}

	for _, item := range items {
		if ctx.Err() != nil {
			stat.errMessage = "run cancelled before all items were processed"
			return stat
		}
		outcome, err := p.ingestItem(ctx, src, item, deps, interval)
		if err != nil {
			log.Warn("item skipped", "url", item.URL, "err", err)
			continue
		}
		switch outcome {
		case outcomeCreated:
			stat.created++
		case outcomeUpdated:
			stat.updated++
		case outcomeUnchanged, outcomeRelocated:
			stat.unchanged++
		case outcomeRejected, outcomeHeld, outcomeGone:
			stat.rejected++
		case outcomeOffTopic:
			stat.offTopic++
			stat.rejected++
		}
	}
	return stat
}

type outcome int

const (
	outcomeSkipped outcome = iota
	outcomeCreated
	outcomeUpdated
	outcomeUnchanged
	outcomeRelocated
	outcomeRejected
	outcomeHeld
	outcomeGone
	outcomeOffTopic
)

// ingestItem xử lý một ứng viên từ feed: bài mới thì publish, bài đã có thì đối
// chiếu nội dung để quyết định giữ nguyên, cập nhật hay tạm giữ bản mới.
func (p *Pipeline) ingestItem(
	ctx context.Context,
	src domain.Source,
	item fetcher.FeedItem,
	deps ingestDeps,
	interval time.Duration,
) (outcome, error) {
	canonicalURL, err := dedupe.CanonicalURL(item.URL)
	if err != nil {
		return outcomeSkipped, fmt.Errorf("canonicalise feed URL: %w", err)
	}

	known, knownErr := p.store.FindArticleByURLHash(ctx, dedupe.URLHash(canonicalURL))
	isKnown := knownErr == nil
	if knownErr != nil && !errors.Is(knownErr, store.ErrNotFound) {
		return outcomeSkipped, knownErr
	}

	if err := deps.allowlist.Check(canonicalURL); err != nil {
		p.log.Warn("refused short link target outside allowlist", "url", canonicalURL, "err", err)
		return outcomeRejected, fmt.Errorf("allowlist: %w", err)
	}

	page, err := p.fetcher.Get(ctx, canonicalURL, interval)
	if err != nil {
		if isKnown && isGoneStatus(page.StatusCode) {
			if err := p.store.MarkSourceGone(ctx, known.ID); err != nil {
				return outcomeSkipped, err
			}
			p.log.Info("article removed at source, short link disabled",
				"article_id", known.ID, "status", page.StatusCode)
			return outcomeGone, nil
		}
		return outcomeSkipped, fmt.Errorf("fetch article: %w", err)
	}

	article, err := p.extractor.Extract(page.Body, canonicalURL)
	if err != nil {
		return outcomeRejected, fmt.Errorf("extract: %w", err)
	}

	// <link rel="canonical"> của trang đáng tin hơn URL trong feed, vốn hay kèm
	// tham số theo dõi hoặc trỏ tới bản AMP (E-01, F-12).
	if declared, err := dedupe.CanonicalURL(article.CanonicalURL); err == nil &&
		deps.allowlist.Allowed(declared) {
		canonicalURL = declared
	}
	urlHash := dedupe.URLHash(canonicalURL)
	contentHash := dedupe.ContentHash(article.Body)

	if !isKnown {
		moved, err := p.store.FindArticleByContentHash(ctx, src.ID, contentHash)
		switch {
		case err == nil && !sameArticle(moved.Title, article.Title):
			p.log.Error("two different titles produced identical body text, extraction is likely broken",
				"source", src.Code, "url", canonicalURL,
				"stored_title", moved.Title, "fetched_title", article.Title)
			return outcomeRejected, fmt.Errorf("body identical to article %d but titles differ", moved.ID)
		case err == nil:
			if err := p.store.RelocateArticle(ctx, moved.ID, canonicalURL, urlHash); err != nil {
				return outcomeSkipped, err
			}
			p.log.Info("article moved to a new URL, short link kept",
				"article_id", moved.ID, "code", moved.ShortCode, "url", canonicalURL)
			return outcomeRelocated, nil
		case err != nil && !errors.Is(err, store.ErrNotFound):
			return outcomeSkipped, err
		}
	}

	if isKnown {
		return p.reconcile(ctx, src, known, article, contentHash, deps)
	}
	return p.publish(ctx, src, item, article, canonicalURL, urlHash, contentHash, deps)
}

// reconcile xử lý bài đã có trong CSDL sau khi crawl lại (U-01 đến U-03, U-07).
func (p *Pipeline) reconcile(
	ctx context.Context,
	src domain.Source,
	known store.StoredArticle,
	article extractor.Article,
	contentHash string,
	deps ingestDeps,
) (outcome, error) {
	seenAt := p.now().UTC()

	if known.ContentHash == contentHash {
		return outcomeUnchanged, p.store.TouchArticle(ctx, known.ID, contentHash, seenAt)
	}
	if !dedupe.NeedsResummarize(known.Title, known.Body, article.Title, article.Body) {
		p.log.Info("article changed below rewrite threshold, summary kept",
			"article_id", known.ID, "ratio", dedupe.ChangeRatio(known.Body, article.Body))
		return outcomeUnchanged, p.store.TouchArticle(ctx, known.ID, contentHash, seenAt)
	}

	tickers, primary := p.tagTickers(deps.tagger, article)
	newsType, _ := p.classifier.Classify(classifier.Input{
		Title: article.Title, Body: article.Body, PrimarySymbols: primary, SourceTier: src.Tier,
	})

	summary, provider, err := p.summarize(ctx, src, article, primary)
	if err != nil {
		return outcomeSkipped, err
	}

	if check := p.validator.Validate(summary.SummaryMD, article.Body); !check.Grounded {
		p.log.Warn("revision failed numeric grounding, published summary kept",
			"article_id", known.ID, "violations", check.Violations)
		return outcomeHeld, p.store.HoldPendingRevision(ctx, known.ID, summary.SummaryMD,
			"numbers_not_grounded: "+check.Error(), seenAt)
	}

	if err := p.store.ApplyRevision(ctx, store.Revision{
		ArticleID: known.ID, Title: article.Title, SummaryMD: summary.SummaryMD,
		Body: article.Body, ContentHash: contentHash, NewsType: newsType,
		Provider: provider, Tickers: tickers, UpdatedAt: seenAt,
	}); err != nil {
		return outcomeSkipped, err
	}
	p.log.Info("article updated at source, summary regenerated",
		"article_id", known.ID, "revision", known.Revision+1)
	return outcomeUpdated, nil
}

// publish ghi một bài mới cùng mã CK và short link trong một transaction (I-03).
func (p *Pipeline) publish(
	ctx context.Context,
	src domain.Source,
	item fetcher.FeedItem,
	article extractor.Article,
	canonicalURL, urlHash, contentHash string,
	deps ingestDeps,
) (outcome, error) {
	fetchedAt := p.now().UTC()

	parsed, hasDate := article.PublishedAt, article.HasDate
	if !hasDate && item.HasDate {
		parsed, hasDate = item.PublishedHint, true
	}
	publishedAt, estimated := extractor.SettlePublishedAt(parsed, hasDate, fetchedAt)
	if extractor.TooOld(publishedAt, fetchedAt) {
		return outcomeRejected, fmt.Errorf("published %s is older than the freshness window", publishedAt.Format(time.RFC3339))
	}

	titleHash := dedupe.TitleSimhash(article.Title)
	clusterID, isCanonical, err := p.resolveCluster(ctx, titleHash)
	if err != nil {
		return outcomeSkipped, err
	}

	tickers, primary := p.tagTickers(deps.tagger, article)

	if verdict := classifier.IsOnTopic(classifier.TopicalInput{
		Title: article.Title, Body: article.Body, TickerCount: len(tickers),
	}); !verdict.OnTopic {
		p.log.Info("article rejected as off topic",
			"source", src.Code, "title", article.Title, "finance_signals", verdict.Signals)
		return outcomeOffTopic, nil
	}

	newsType, _ := p.classifier.Classify(classifier.Input{
		Title: article.Title, Body: article.Body, PrimarySymbols: primary, SourceTier: src.Tier,
	})

	summary, provider, summaryErr := p.summarize(ctx, src, article, primary)
	if summaryErr != nil && summary.SummaryMD == "" {
		return outcomeSkipped, summaryErr
	}

	status, rejectReason := domain.StatusPublished, ""
	if summaryErr != nil {
		// Tóm tắt vi phạm rule hình thức (thường là ngôn ngữ khuyến nghị đầu tư
		// copy từ báo cáo môi giới). Giữ bài lại ở pending_review thay vì vứt:
		// mất bài là mất tin, còn publish lời khuyên đầu tư là rủi ro pháp lý.
		status, rejectReason = domain.StatusPendingReview, "summary_rejected"
		p.log.Warn("summary rejected, article held for review",
			"url", canonicalURL, "err", summaryErr)
	}

	duplicate, err := p.store.SummaryAlreadyUsed(ctx, src.ID, summary.SummaryMD, 0)
	if err != nil {
		return outcomeSkipped, err
	}
	if duplicate {
		status, rejectReason = domain.StatusPendingReview, "duplicate_summary"
		p.log.Error("summary identical to another article from the same source, extraction is likely broken",
			"source", src.Code, "url", canonicalURL, "title", article.Title)
	}

	if check := p.validator.Validate(summary.SummaryMD, article.Body); !check.Grounded {
		status, rejectReason = domain.StatusPendingReview, "numbers_not_grounded"
		p.log.Warn("summary failed numeric grounding", "url", canonicalURL, "violations", check.Violations)
	}
	// Rule sao chép nguyên văn chỉ áp cho provider sinh ngữ (LLM). `extractive`
	// hoạt động bằng cách trích nguyên câu — đó là cơ chế của nó, và cũng là thứ
	// giữ cho mọi con số luôn truy vết được về bài gốc.
	if provider != p.extractive.Name() {
		if run := summarizer.LongestVerbatimRun(summary.SummaryMD, article.Body); run >= summarizer.VerbatimRunLimit {
			status, rejectReason = domain.StatusPendingReview, "verbatim_copy"
			p.log.Warn("summary copies source verbatim", "url", canonicalURL, "words", run)
		}
	}

	record := store.NewArticle{
		Article: domain.Article{
			SourceID: src.ID, CanonicalURL: canonicalURL, URLHash: urlHash,
			Title: article.Title, SummaryMD: summary.SummaryMD, Excerpt: article.Body,
			NewsType: newsType, TitleSimhash: dedupe.EncodeSimhash(titleHash),
			ClusterID: clusterID, IsCanonicalInCluster: isCanonical,
			Status: status, SummaryProvider: provider,
			PublishedAt: publishedAt, FetchedAt: fetchedAt,
			ContentHash: contentHash, PublishedAtEstimated: estimated,
			IsLive: article.IsLive, RejectReason: rejectReason,
		},
		Tickers:   tickers,
		TargetURL: canonicalURL,
	}

	for attempt := 1; attempt <= shortCodeAttempts; attempt++ {
		code, err := shortlink.GenerateCode()
		if err != nil {
			return outcomeSkipped, err
		}
		record.ShortCode = code

		if _, err := p.store.CreateArticle(ctx, record); err != nil {
			if errors.Is(err, store.ErrDuplicateArticle) {
				return outcomeUnchanged, nil
			}
			if errors.Is(err, store.ErrCodeConflict) && attempt < shortCodeAttempts {
				continue
			}
			return outcomeSkipped, err
		}
		if status != domain.StatusPublished {
			return outcomeHeld, nil
		}
		return outcomeCreated, nil
	}
	return outcomeSkipped, fmt.Errorf("no unique short code after %d attempts", shortCodeAttempts)
}

func (p *Pipeline) summarize(ctx context.Context, src domain.Source, article extractor.Article, primary []string) (summarizer.Output, string, error) {
	out, provider, err := summarizer.SummarizeWithFallback(ctx, p.summarizer, p.extractive, summarizer.Input{
		Title: article.Title, Body: article.Body, SourceName: src.Name, Tickers: primary,
	}, p.log)
	if err != nil {
		// Trả kèm nội dung: caller cần nó để lưu bản chờ duyệt thay vì mất bài.
		return out, provider, fmt.Errorf("summarize: %w", err)
	}
	return out, provider, nil
}

func (p *Pipeline) tagTickers(tg *tagger.Tagger, article extractor.Article) ([]domain.TickerRef, []string) {
	matches := tg.Tag(tagger.Input{Title: article.Title, Body: article.Body})
	refs := make([]domain.TickerRef, 0, len(matches))
	primary := make([]string, 0, len(matches))
	for _, m := range matches {
		refs = append(refs, domain.TickerRef{Symbol: m.Symbol, Relevance: m.Relevance, Score: m.Score})
		if m.Relevance == domain.RelevancePrimary {
			primary = append(primary, m.Symbol)
		}
	}
	return refs, primary
}

// resolveCluster tìm cụm trùng trong cửa sổ 72h; không thấy thì mở cụm mới.
func (p *Pipeline) resolveCluster(ctx context.Context, titleHash uint64) (*int64, bool, error) {
	since := p.now().Add(-ClusterWindow)
	rows, err := p.store.RecentSimhashes(ctx, since)
	if err != nil {
		return nil, false, err
	}
	for _, r := range rows {
		other, err := dedupe.DecodeSimhash(r.Simhash)
		if err != nil {
			continue
		}
		if dedupe.IsNearDuplicate(titleHash, other) && r.ClusterID != nil {
			id := *r.ClusterID
			return &id, false, nil
		}
	}
	id, err := p.store.CreateCluster(ctx, dedupe.EncodeSimhash(titleHash), since)
	if err != nil {
		return nil, false, err
	}
	return &id, true, nil
}

// sameArticle so tiêu đề để chắc chắn hai URL cùng trỏ về một bài.
//
// Chỉ dựa vào content_hash là không đủ: khi selector bóc tách trượt, nhiều bài
// của cùng một nguồn cho ra đúng một đoạn boilerplate và sẽ bị gộp nhầm thành
// "bài đã đổi URL", ghi đè canonical_url của nhau.
func sameArticle(storedTitle, fetchedTitle string) bool {
	return textutil.Normalize(storedTitle) == textutil.Normalize(fetchedTitle)
}

func isGoneStatus(code int) bool {
	return code == http.StatusNotFound || code == http.StatusGone
}

func toTaggerTickers(list []domain.Ticker) []tagger.Ticker {
	out := make([]tagger.Ticker, 0, len(list))
	for _, t := range list {
		entry := tagger.Ticker{
			Symbol:          t.Symbol,
			Sector:          t.Sector,
			InVN30:          t.InVN30,
			NegativeAliases: t.NegativeAliases,
		}
		for _, a := range t.Aliases {
			entry.Aliases = append(entry.Aliases, a.Alias)
		}
		out = append(out, entry)
	}
	return out
}
