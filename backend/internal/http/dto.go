// Package http là tầng trình bày REST v1 (mục 6 architecture.md).
// Mọi thông điệp lỗi bằng tiếng Việt vì hiển thị trực tiếp cho người dùng.
package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/tenpoint/tenpoint-api/internal/domain"
)

// Mã lỗi thống nhất.
const (
	CodeNotFound       = "not_found"
	CodeInvalidRequest = "invalid_request"
	CodeUnauthorized   = "unauthorized"
	CodeConflict       = "conflict"
	CodeGone           = "gone"
	CodeInternal       = "internal_error"
)

// ErrorBody là khối lỗi thống nhất: {"error":{"code":..., "message":...}}.
type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type errorEnvelope struct {
	Error ErrorBody `json:"error"`
}

// ListMeta là khối phân trang.
type ListMeta struct {
	Total   int  `json:"total"`
	Limit   int  `json:"limit"`
	Offset  int  `json:"offset"`
	HasMore bool `json:"has_more"`
}

type listEnvelope struct {
	Data any       `json:"data"`
	Meta *ListMeta `json:"meta,omitempty"`
}

type itemEnvelope struct {
	Data any `json:"data"`
}

// SourceDTO là khối nguồn trong response.
//
// ArticleTitle là tiêu đề bài tại nguồn; FE ghi nó vào title và aria-label của
// link để người dùng biết bấm vào sẽ ra bài nào. Đây là câu trả lời trực tiếp
// cho phản hồi "link chưa đưa đến thông tin trang cụ thể".
type SourceDTO struct {
	Name         string `json:"name"`
	Domain       string `json:"domain"`
	Tier         int    `json:"tier"`
	ArticleTitle string `json:"article_title"`
}

// TickerRefDTO là mã CK gắn với bài.
type TickerRefDTO struct {
	Symbol    string `json:"symbol"`
	Relevance string `json:"relevance"`
}

// NewsItem là một dòng của bảng digest.
type NewsItem struct {
	ID                   int64          `json:"id"`
	PublishedAt          string         `json:"published_at"`
	PublishedDateVN      string         `json:"published_date_vn"`
	Title                string         `json:"title"`
	SummaryMD            string         `json:"summary_md"`
	NewsType             string         `json:"news_type"`
	NewsTypeLabel        string         `json:"news_type_label"`
	Tickers              []TickerRefDTO `json:"tickers"`
	Source               SourceDTO      `json:"source"`
	ShortLink            string         `json:"short_link"`
	ShortLinkAlive       bool           `json:"short_link_alive"`
	IsDemo               bool           `json:"is_demo"`
	UpdatedAt            *string        `json:"updated_at"`
	PublishedAtEstimated bool           `json:"published_at_estimated"`
	Revision             int            `json:"revision"`
}

// NewsDetail bổ sung các tin cùng cụm trùng.
type NewsDetail struct {
	NewsItem
	AlsoReportedBy []NewsItem `json:"also_reported_by"`
}

// TickerItem là một mã trong danh sách lọc/autocomplete.
type TickerItem struct {
	Symbol       string `json:"symbol"`
	CompanyName  string `json:"company_name"`
	ShortName    string `json:"short_name"`
	Exchange     string `json:"exchange"`
	Sector       string `json:"sector"`
	InVN30       bool   `json:"in_vn30"`
	ArticleCount int    `json:"article_count"`
}

// TickerDetail gồm thông tin mã + tin gần đây + research notes.
type TickerDetail struct {
	TickerItem
	RecentNews    []NewsItem `json:"recent_news"`
	ResearchNotes []NoteItem `json:"research_notes"`
}

// NoteItem là một research note trong danh sách.
type NoteItem struct {
	Slug            string `json:"slug"`
	Symbol          string `json:"symbol"`
	Title           string `json:"title"`
	SectionHeading  string `json:"section_heading"`
	PublishedAt     string `json:"published_at"`
	PublishedDateVN string `json:"published_date_vn"`
}

// NotePoint là một luận điểm đầu tư.
type NotePoint struct {
	Ordinal int    `json:"ordinal"`
	Lead    string `json:"lead"`
	Body    string `json:"body"`
}

// NoteDetail là research note đầy đủ.
type NoteDetail struct {
	NoteItem
	Disclaimer string      `json:"disclaimer"`
	Points     []NotePoint `json:"points"`
}

// NewsTypeCount là loại tin kèm số lượng.
type NewsTypeCount struct {
	Value string `json:"value"`
	Label string `json:"label"`
	Count int    `json:"count"`
}

// MetaResponse phục vụ GET /api/v1/meta.
type MetaResponse struct {
	LastCrawlAt     *string         `json:"last_crawl_at"`
	IsStale         bool            `json:"is_stale"`
	TotalArticles   int             `json:"total_articles"`
	TickersCount    int             `json:"tickers_count"`
	HasRealArticles bool            `json:"has_real_articles"`
	NewsTypes       []NewsTypeCount `json:"news_types"`
}

// ---------------------------------------------------------------- serializers

func toNewsItem(a domain.Article, loc *time.Location) NewsItem {
	tickers := make([]TickerRefDTO, 0, len(a.Tickers))
	for _, t := range a.Tickers {
		tickers = append(tickers, TickerRefDTO{Symbol: t.Symbol, Relevance: t.Relevance})
	}
	shortLink := ""
	if a.ShortCode != "" {
		shortLink = "/r/" + a.ShortCode
	}
	var updatedAt *string
	if a.UpdatedAt != nil {
		formatted := a.UpdatedAt.UTC().Format(time.RFC3339)
		updatedAt = &formatted
	}
	return NewsItem{
		ID:              a.ID,
		PublishedAt:     a.PublishedAt.UTC().Format(time.RFC3339),
		PublishedDateVN: a.PublishedDateVN(loc),
		Title:           a.Title,
		SummaryMD:       a.SummaryMD,
		NewsType:        string(a.NewsType),
		NewsTypeLabel:   a.NewsType.Label(),
		Tickers:         tickers,
		Source: SourceDTO{
			Name:         a.Source.Name,
			Domain:       a.Source.Domain,
			Tier:         a.Source.Tier,
			ArticleTitle: a.Title,
		},
		ShortLink:            shortLink,
		ShortLinkAlive:       a.ShortLinkAlive,
		IsDemo:               a.IsDemo,
		UpdatedAt:            updatedAt,
		PublishedAtEstimated: a.PublishedAtEstimated,
		Revision:             a.Revision,
	}
}

func toNewsItems(list []domain.Article, loc *time.Location) []NewsItem {
	out := make([]NewsItem, 0, len(list))
	for _, a := range list {
		out = append(out, toNewsItem(a, loc))
	}
	return out
}

func toNoteItem(n domain.ResearchNote, loc *time.Location) NoteItem {
	return NoteItem{
		Slug:            n.Slug,
		Symbol:          n.Symbol,
		Title:           n.Title,
		SectionHeading:  n.SectionHeading,
		PublishedAt:     n.PublishedAt.UTC().Format(time.RFC3339),
		PublishedDateVN: n.PublishedDateVN(loc),
	}
}

// ---------------------------------------------------------------- writers

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(payload); err != nil {
		slog.Default().Error("write JSON response", "err", err)
	}
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, errorEnvelope{Error: ErrorBody{Code: code, Message: message}})
}
