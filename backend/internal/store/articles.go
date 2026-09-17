package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/tenpoint/tenpoint-api/internal/domain"
)

// ErrCodeConflict báo mã short link đã tồn tại — caller sinh mã mới rồi thử lại.
var ErrCodeConflict = errors.New("short link code already exists")

// ErrDuplicateArticle báo bài đã có trong DB (đụng UNIQUE url_hash) — bỏ qua,
// KHÔNG thử lại.
var ErrDuplicateArticle = errors.New("article already exists")

// articleFilterSQL dùng chung cho cả truy vấn đếm và truy vấn phân trang.
// $1 tu, $2 den, $3 news_types, $4 symbols, $5 relevance_all, $6 q, $7 include_demo
const articleFilterSQL = `
    a.status IN ('published', 'source_gone')
    AND ($7::boolean OR a.is_demo = false)
    AND ($1::timestamptz IS NULL OR a.published_at >= $1)
    AND ($2::timestamptz IS NULL OR a.published_at < $2)
    AND ($3::text[] IS NULL OR a.news_type = ANY($3))
    AND ($4::text[] IS NULL OR EXISTS (
            SELECT 1 FROM article_tickers t
            WHERE t.article_id = a.id
              AND t.symbol = ANY($4)
              AND ($5::boolean OR t.relevance = 'primary')))
    AND ($6::text = '' OR a.search_vec @@ plainto_tsquery('simple', tenpoint_unaccent($6)))`

const articleSelectColumns = `
    a.id, a.published_at, a.title, a.summary_md, a.news_type,
    a.status, a.summary_provider, a.canonical_url, a.cluster_id,
    a.is_demo, a.revision, a.updated_at, a.published_at_estimated, a.is_live,
    s.id, s.name, s.domain, s.tier,
    COALESCE(sl.code, ''), COALESCE(sl.target_alive, false)`

// ListArticles trả về trang dữ liệu cho GET /api/v1/news kèm tổng số bản ghi.
func (s *Store) ListArticles(ctx context.Context, f domain.ArticleFilter) ([]domain.Article, int, error) {
	var (
		from, to  any
		newsTypes any
		symbols   any
	)
	if f.From != nil {
		from = *f.From
	}
	if f.To != nil {
		to = *f.To
	}
	if len(f.NewsTypes) > 0 {
		vals := make([]string, 0, len(f.NewsTypes))
		for _, nt := range f.NewsTypes {
			vals = append(vals, string(nt))
		}
		newsTypes = vals
	}
	if len(f.Symbols) > 0 {
		symbols = f.Symbols
	}
	relevanceAll := f.Relevance == "all"

	args := []any{from, to, newsTypes, symbols, relevanceAll, f.Query, f.IncludeDemo}

	var total int
	countSQL := "SELECT count(*) FROM articles a WHERE" + articleFilterSQL
	if err := s.pool.QueryRow(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count articles: %w", err)
	}

	limit := f.Limit
	if limit <= 0 {
		limit = 20
	}
	listSQL := `SELECT` + articleSelectColumns + `
        FROM articles a
        JOIN sources s ON s.id = a.source_id
        LEFT JOIN short_links sl ON sl.article_id = a.id
        WHERE` + articleFilterSQL + `
        ORDER BY a.published_at DESC, a.id DESC
        LIMIT $8 OFFSET $9`

	rows, err := s.pool.Query(ctx, listSQL, append(args, limit, f.Offset)...)
	if err != nil {
		return nil, 0, fmt.Errorf("list articles: %w", err)
	}
	defer rows.Close()

	articles, err := scanArticles(rows)
	if err != nil {
		return nil, 0, err
	}
	if err := s.attachTickers(ctx, articles); err != nil {
		return nil, 0, err
	}
	return articles, total, nil
}

// GetArticle đọc 1 tin theo id (chỉ tin đã publish).
func (s *Store) GetArticle(ctx context.Context, id int64) (domain.Article, error) {
	q := `SELECT` + articleSelectColumns + `
        FROM articles a
        JOIN sources s ON s.id = a.source_id
        LEFT JOIN short_links sl ON sl.article_id = a.id
        WHERE a.id = $1 AND a.status IN ('published', 'source_gone')`

	rows, err := s.pool.Query(ctx, q, id)
	if err != nil {
		return domain.Article{}, fmt.Errorf("query article %d: %w", id, err)
	}
	defer rows.Close()

	list, err := scanArticles(rows)
	if err != nil {
		return domain.Article{}, err
	}
	if len(list) == 0 {
		return domain.Article{}, fmt.Errorf("article id=%d: %w", id, ErrNotFound)
	}
	if err := s.attachTickers(ctx, list); err != nil {
		return domain.Article{}, err
	}
	return list[0], nil
}

// ListClusterSiblings trả về các tin cùng cụm trùng (trừ chính nó).
func (s *Store) ListClusterSiblings(ctx context.Context, articleID int64, clusterID *int64) ([]domain.Article, error) {
	if clusterID == nil {
		return nil, nil
	}
	q := `SELECT` + articleSelectColumns + `
        FROM articles a
        JOIN sources s ON s.id = a.source_id
        LEFT JOIN short_links sl ON sl.article_id = a.id
        WHERE a.cluster_id = $1 AND a.id <> $2 AND a.status = 'published' AND a.is_demo = false
        ORDER BY a.published_at DESC
        LIMIT 10`

	rows, err := s.pool.Query(ctx, q, *clusterID, articleID)
	if err != nil {
		return nil, fmt.Errorf("query cluster siblings of article %d: %w", articleID, err)
	}
	defer rows.Close()
	return scanArticles(rows)
}

// ListArticlesBySymbol trả về tin gần đây của một mã.
func (s *Store) ListArticlesBySymbol(ctx context.Context, symbol string, limit int) ([]domain.Article, error) {
	if limit <= 0 {
		limit = 20
	}
	q := `SELECT` + articleSelectColumns + `
        FROM articles a
        JOIN sources s ON s.id = a.source_id
        LEFT JOIN short_links sl ON sl.article_id = a.id
        JOIN article_tickers t ON t.article_id = a.id
        WHERE t.symbol = $1 AND a.status IN ('published', 'source_gone')
        ORDER BY a.published_at DESC, a.id DESC
        LIMIT $2`

	rows, err := s.pool.Query(ctx, q, symbol, limit)
	if err != nil {
		return nil, fmt.Errorf("query articles for ticker %s: %w", symbol, err)
	}
	defer rows.Close()

	list, err := scanArticles(rows)
	if err != nil {
		return nil, err
	}
	if err := s.attachTickers(ctx, list); err != nil {
		return nil, err
	}
	return list, nil
}

func scanArticles(rows pgx.Rows) ([]domain.Article, error) {
	var out []domain.Article
	for rows.Next() {
		var a domain.Article
		var newsType string
		if err := rows.Scan(
			&a.ID, &a.PublishedAt, &a.Title, &a.SummaryMD, &newsType,
			&a.Status, &a.SummaryProvider, &a.CanonicalURL, &a.ClusterID,
			&a.IsDemo, &a.Revision, &a.UpdatedAt, &a.PublishedAtEstimated, &a.IsLive,
			&a.Source.ID, &a.Source.Name, &a.Source.Domain, &a.Source.Tier,
			&a.ShortCode, &a.ShortLinkAlive,
		); err != nil {
			return nil, fmt.Errorf("scan article row: %w", err)
		}
		a.NewsType = domain.NewsType(newsType)
		a.SourceID = a.Source.ID
		out = append(out, a)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate article rows: %w", err)
	}
	return out, nil
}

// attachTickers nạp mã CK cho một tập bài trong 1 truy vấn (tránh N+1).
func (s *Store) attachTickers(ctx context.Context, articles []domain.Article) error {
	if len(articles) == 0 {
		return nil
	}
	ids := make([]int64, 0, len(articles))
	for _, a := range articles {
		ids = append(ids, a.ID)
	}
	rows, err := s.pool.Query(ctx, `
        SELECT article_id, symbol, relevance, score
        FROM article_tickers
        WHERE article_id = ANY($1)
        ORDER BY article_id, score DESC, symbol`, ids)
	if err != nil {
		return fmt.Errorf("query article tickers: %w", err)
	}
	defer rows.Close()

	byArticle := map[int64][]domain.TickerRef{}
	for rows.Next() {
		var id int64
		var ref domain.TickerRef
		var score float32
		if err := rows.Scan(&id, &ref.Symbol, &ref.Relevance, &score); err != nil {
			return fmt.Errorf("scan article ticker: %w", err)
		}
		ref.Score = float64(score)
		byArticle[id] = append(byArticle[id], ref)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate article tickers: %w", err)
	}
	for i := range articles {
		articles[i].Tickers = byArticle[articles[i].ID]
	}
	return nil
}

// ---------------------------------------------------------------- ghi dữ liệu

// NewArticle gom mọi thứ cần ghi cho một bài mới trong một transaction.
type NewArticle struct {
	Article   domain.Article
	Tickers   []domain.TickerRef
	ShortCode string
	TargetURL string
}

// ExistsByURLHash kiểm tra bài đã có chưa (L0 dedupe — rẻ nhất).
func (s *Store) ExistsByURLHash(ctx context.Context, urlHash string) (bool, error) {
	var exists bool
	err := s.pool.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM articles WHERE url_hash = $1)", urlHash).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check url_hash: %w", err)
	}
	return exists, nil
}

// SimhashRow là bản ghi tối thiểu dùng để so trùng tiêu đề trong cửa sổ 72h.
type SimhashRow struct {
	ID        int64
	Simhash   string
	ClusterID *int64
}

// RecentSimhashes trả về simhash tiêu đề của các bài trong cửa sổ thời gian.
func (s *Store) RecentSimhashes(ctx context.Context, since time.Time) ([]SimhashRow, error) {
	rows, err := s.pool.Query(ctx, `
        SELECT id, title_simhash, cluster_id
        FROM articles
        WHERE published_at >= $1 AND title_simhash <> ''`, since)
	if err != nil {
		return nil, fmt.Errorf("query recent simhashes: %w", err)
	}
	defer rows.Close()

	var out []SimhashRow
	for rows.Next() {
		var r SimhashRow
		if err := rows.Scan(&r.ID, &r.Simhash, &r.ClusterID); err != nil {
			return nil, fmt.Errorf("scan simhash: %w", err)
		}
		out = append(out, r)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate simhashes: %w", err)
	}
	return out, nil
}

// CreateCluster tạo cụm trùng mới và trả về id.
func (s *Store) CreateCluster(ctx context.Context, simhash string, windowStart time.Time) (int64, error) {
	var id int64
	err := s.pool.QueryRow(ctx, `
        INSERT INTO article_clusters (representative_simhash, window_start)
        VALUES ($1, $2) RETURNING id`, simhash, windowStart).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("create cluster: %w", err)
	}
	return id, nil
}

// CreateArticle ghi bài + mã CK + short link trong MỘT transaction.
// Trả ErrCodeConflict khi mã short link bị trùng để caller sinh mã khác.
func (s *Store) CreateArticle(ctx context.Context, in NewArticle) (int64, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("begin article insert tx: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	a := in.Article
	var id int64
	err = tx.QueryRow(ctx, `
        INSERT INTO articles (
            source_id, canonical_url, url_hash, title, summary_md, excerpt,
            news_type, title_simhash, cluster_id, is_canonical_in_cluster,
            status, summary_provider, published_at, fetched_at,
            content_hash, last_seen_at, published_at_estimated, is_live, reject_reason)
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19)
        RETURNING id`,
		a.SourceID, a.CanonicalURL, a.URLHash, a.Title, a.SummaryMD, a.Excerpt,
		string(a.NewsType), a.TitleSimhash, a.ClusterID, a.IsCanonicalInCluster,
		a.Status, a.SummaryProvider, a.PublishedAt, a.FetchedAt,
		a.ContentHash, a.FetchedAt, a.PublishedAtEstimated, a.IsLive, a.RejectReason).Scan(&id)
	if err != nil {
		if isUniqueViolation(err) {
			return 0, fmt.Errorf("article %s: %w", a.CanonicalURL, ErrDuplicateArticle)
		}
		return 0, fmt.Errorf("insert article %s: %w", a.CanonicalURL, err)
	}

	for _, t := range in.Tickers {
		if _, err := tx.Exec(ctx, `
            INSERT INTO article_tickers (article_id, symbol, relevance, score)
            VALUES ($1,$2,$3,$4)
            ON CONFLICT (article_id, symbol) DO NOTHING`,
			id, t.Symbol, t.Relevance, float32(t.Score)); err != nil {
			return 0, fmt.Errorf("attach ticker %s to article %d: %w", t.Symbol, id, err)
		}
	}

	if in.ShortCode != "" && in.TargetURL != "" {
		_, err := tx.Exec(ctx, `
            INSERT INTO short_links (code, article_id, target_url)
            VALUES ($1,$2,$3)`, in.ShortCode, id, in.TargetURL)
		if err != nil {
			if isUniqueViolation(err) {
				return 0, fmt.Errorf("code %s: %w", in.ShortCode, ErrCodeConflict)
			}
			return 0, fmt.Errorf("create short link for article %d: %w", id, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("commit article insert: %w", err)
	}
	return id, nil
}

// MetaStats phục vụ GET /api/v1/meta.
type MetaStats struct {
	TotalArticles int
	TickersCount  int
	NewsTypes     map[domain.NewsType]int
	LastCrawlAt   *time.Time
}

// Meta tổng hợp số liệu trạng thái hệ thống.
func (s *Store) Meta(ctx context.Context) (MetaStats, error) {
	out := MetaStats{NewsTypes: map[domain.NewsType]int{}}

	if err := s.pool.QueryRow(ctx,
		"SELECT count(*) FROM articles WHERE status = 'published'").Scan(&out.TotalArticles); err != nil {
		return out, fmt.Errorf("count published articles: %w", err)
	}
	if err := s.pool.QueryRow(ctx,
		"SELECT count(*) FROM tickers WHERE status = 'active'").Scan(&out.TickersCount); err != nil {
		return out, fmt.Errorf("count active tickers: %w", err)
	}

	rows, err := s.pool.Query(ctx, `
        SELECT news_type, count(*)
        FROM articles WHERE status = 'published'
        GROUP BY news_type`)
	if err != nil {
		return out, fmt.Errorf("count articles by news type: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var nt string
		var n int
		if err := rows.Scan(&nt, &n); err != nil {
			return out, fmt.Errorf("scan news type stats: %w", err)
		}
		out.NewsTypes[domain.NewsType(nt)] = n
	}
	if err := rows.Err(); err != nil {
		return out, fmt.Errorf("iterate news type stats: %w", err)
	}

	last, err := s.LastCrawlAt(ctx)
	if err != nil {
		return out, err
	}
	out.LastCrawlAt = last
	return out, nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
