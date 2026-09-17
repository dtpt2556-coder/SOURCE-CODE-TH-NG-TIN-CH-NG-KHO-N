package domain

import "time"

// Trạng thái vòng đời của một bài viết trong pipeline.
const (
	StatusPublished     = "published"
	StatusPendingReview = "pending_review"
	StatusRejected      = "rejected"
	// StatusSourceGone: bài gốc đã bị gỡ khỏi nguồn. Tóm tắt vẫn giữ vì đó là
	// nội dung của TenPoint, nhưng link không còn bấm được (U-04).
	StatusSourceGone = "source_gone"
)

// Mức độ liên quan của một mã CK với bài viết (mục 8 architecture.md).
const (
	RelevancePrimary   = "primary"
	RelevanceMentioned = "mentioned"
)

// Tên provider sinh tóm tắt.
const (
	SummaryProviderAnthropic  = "anthropic"
	SummaryProviderExtractive = "extractive"
)

// TickerRef là một mã CK đã gắn vào bài viết.
type TickerRef struct {
	Symbol    string
	Relevance string
	Score     float64
}

// Article là entity thuần, không phụ thuộc hạ tầng.
// `Excerpt` chỉ dùng nội bộ (validator) — KHÔNG BAO GIỜ trả ra API public (ADR-006).
type Article struct {
	ID                   int64
	SourceID             int
	CanonicalURL         string
	URLHash              string
	Title                string
	SummaryMD            string
	Excerpt              string
	NewsType             NewsType
	TitleSimhash         string
	ClusterID            *int64
	IsCanonicalInCluster bool
	Status               string
	SummaryProvider      string
	PublishedAt          time.Time
	FetchedAt            time.Time

	ContentHash          string
	Revision             int
	UpdatedAt            *time.Time
	LastSeenAt           *time.Time
	PublishedAtEstimated bool
	IsLive               bool
	IsDemo               bool
	RejectReason         string

	// Dữ liệu join sẵn phục vụ tầng trình bày.
	Source         Source
	Tickers        []TickerRef
	ShortCode      string
	ShortLinkAlive bool
}

// PublishedDateVN định dạng ngày theo cột "Ngày" của bảng digest: DD/MM/YYYY
// theo múi giờ truyền vào (Asia/Ho_Chi_Minh).
func (a Article) PublishedDateVN(loc *time.Location) string {
	if loc == nil {
		loc = time.UTC
	}
	return a.PublishedAt.In(loc).Format("02/01/2006")
}

// PrimaryTickers lọc các mã ở mức primary, giữ nguyên thứ tự.
func (a Article) PrimaryTickers() []TickerRef {
	out := make([]TickerRef, 0, len(a.Tickers))
	for _, t := range a.Tickers {
		if t.Relevance == RelevancePrimary {
			out = append(out, t)
		}
	}
	return out
}

// ArticleFilter là bộ lọc cho GET /api/v1/news (mục 6 architecture.md).
type ArticleFilter struct {
	Symbols   []string
	Relevance string // "primary" | "all"
	NewsTypes []NewsType
	From      *time.Time
	To        *time.Time
	Query     string
	Limit     int
	Offset    int
	// IncludeDemo chỉ bật khi CSDL chưa có bài thật nào, để trang digest không
	// trắng trơn ở lần chạy đầu (T3 mục 0 của ingest-edge-cases).
	IncludeDemo bool
}

// CrawlRun ghi nhận một lần chạy pipeline.
type CrawlRun struct {
	ID               int64
	Trigger          string
	StartedAt        time.Time
	FinishedAt       *time.Time
	ArticlesFound    int
	ArticlesNew      int
	ArticlesRejected int
	Errors           int
}

// Trigger hợp lệ của CrawlRun.
const (
	TriggerCron   = "cron"
	TriggerManual = "manual"
)

// CrawlRunSource là kết quả xử lý của một nguồn trong một run.
type CrawlRunSource struct {
	ID           int64
	RunID        int64
	SourceID     int
	Found        int
	NewCount     int
	ErrorMessage string
	// Truncated: feed vượt MAX_ARTICLES_PER_SOURCE nên có khoảng trống dữ liệu.
	Truncated bool
	// OffTopic: số bài bị loại vì không thuộc chủ đề tài chính.
	OffTopic int
}
