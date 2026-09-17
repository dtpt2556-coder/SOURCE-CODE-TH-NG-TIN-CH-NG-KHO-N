package http

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tenpoint/tenpoint-api/internal/domain"
)

func vnLoc(t *testing.T) *time.Location {
	t.Helper()
	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	require.NoError(t, err)
	return loc
}

// Kiểm tra response khớp NGUYÊN VĂN JSON mẫu ở mục 6 architecture.md.
func TestToNewsItemMatchesAPIContract(t *testing.T) {
	loc := vnLoc(t)
	updatedAt := time.Date(2026, 8, 26, 9, 0, 0, 0, time.UTC)
	article := domain.Article{
		ID:             1042,
		Revision:       2,
		UpdatedAt:      &updatedAt,
		ShortLinkAlive: true,
		PublishedAt:    time.Date(2026, 8, 26, 1, 30, 0, 0, time.UTC),
		Title:          "12 ngân hàng cam kết 408.000 tỷ đồng tín dụng cho DNNVV",
		SummaryMD:      "12 ngân hàng cam kết **408.000 tỷ đồng** tín dụng cho DNNVV.",
		NewsType:       domain.NewsTypeNganh,
		ShortCode:      "a7Kx2p",
		Source:         domain.Source{Name: "VnExpress", Domain: "vnexpress.net", Tier: 1},
		Tickers: []domain.TickerRef{
			{Symbol: "BID", Relevance: domain.RelevancePrimary},
			{Symbol: "VCB", Relevance: domain.RelevancePrimary},
		},
	}

	item := toNewsItem(article, loc)
	assert.Equal(t, int64(1042), item.ID)
	assert.Equal(t, "2026-08-26T01:30:00Z", item.PublishedAt)
	assert.Equal(t, "26/08/2026", item.PublishedDateVN, "published_date_vn phải là DD/MM/YYYY giờ VN")
	assert.Equal(t, "nganh", item.NewsType)
	assert.Equal(t, "Ngành", item.NewsTypeLabel)
	assert.Equal(t, "/r/a7Kx2p", item.ShortLink)
	assert.Equal(t, "vnexpress.net", item.Source.Domain)
	assert.Len(t, item.Tickers, 2)

	// FE ghi article_title vào title/aria-label của link để người dùng biết bấm
	// vào sẽ ra bài nào — câu trả lời cho "link chưa đưa đến thông tin trang cụ thể".
	assert.Equal(t, article.Title, item.Source.ArticleTitle)
	assert.True(t, item.ShortLinkAlive)
	assert.False(t, item.IsDemo)
	assert.False(t, item.PublishedAtEstimated)
	assert.Equal(t, 2, item.Revision)
	require.NotNil(t, item.UpdatedAt)
	assert.Equal(t, "2026-08-26T09:00:00Z", *item.UpdatedAt)

	raw, err := json.Marshal(item)
	require.NoError(t, err)
	for _, field := range []string{
		`"id"`, `"published_at"`, `"published_date_vn"`, `"title"`, `"summary_md"`,
		`"news_type"`, `"news_type_label"`, `"tickers"`, `"symbol"`, `"relevance"`,
		`"source"`, `"name"`, `"domain"`, `"tier"`, `"short_link"`,
		`"article_title"`, `"short_link_alive"`, `"is_demo"`, `"updated_at"`,
		`"published_at_estimated"`,
	} {
		assert.Contains(t, string(raw), field)
	}
}

// Ngày biên: 23:30 UTC ngày 26/8 = 06:30 ngày 27/8 giờ VN.
func TestPublishedDateVNCrossesMidnightCorrectly(t *testing.T) {
	loc := vnLoc(t)
	a := domain.Article{PublishedAt: time.Date(2026, 8, 26, 23, 30, 0, 0, time.UTC)}
	assert.Equal(t, "27/08/2026", a.PublishedDateVN(loc))
}

func TestWriteErrorFormat(t *testing.T) {
	rec := httptest.NewRecorder()
	writeError(rec, 404, CodeNotFound, "Không tìm thấy tin này.")

	assert.Equal(t, 404, rec.Code)
	assert.Equal(t, "application/json; charset=utf-8", rec.Header().Get("Content-Type"))

	var body struct {
		Error ErrorBody `json:"error"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	assert.Equal(t, "not_found", body.Error.Code)
	assert.Equal(t, "Không tìm thấy tin này.", body.Error.Message)
	assert.Contains(t, rec.Body.String(), "Không tìm thấy tin này.",
		"tiếng Việt không được bị escape trong JSON")
}

func TestParseDateRange(t *testing.T) {
	loc := vnLoc(t)

	from, to, err := parseDateRange("2026-08-20", "2026-08-27", loc)
	require.NoError(t, err)
	require.NotNil(t, from)
	require.NotNil(t, to)
	// 00:00 ngày 20/8 giờ VN == 17:00 ngày 19/8 UTC.
	assert.Equal(t, "2026-08-19T17:00:00Z", from.Format(time.RFC3339))
	// `den` bao gồm cả ngày đó -> chặn trên là 00:00 ngày 28/8 giờ VN.
	assert.Equal(t, "2026-08-27T17:00:00Z", to.Format(time.RFC3339))

	// Mặc định: 7 ngày gần nhất.
	from, to, err = parseDateRange("", "", loc)
	require.NoError(t, err)
	require.NotNil(t, from)
	assert.Nil(t, to)
	assert.WithinDuration(t, time.Now().AddDate(0, 0, -defaultDays), *from, 25*time.Hour)

	_, _, err = parseDateRange("26-08-2026", "", loc)
	assert.Error(t, err)

	_, _, err = parseDateRange("2026-08-27", "2026-08-20", loc)
	assert.Error(t, err, "tu phải trước den")
}

func TestToNewsItemWithoutShortLink(t *testing.T) {
	item := toNewsItem(domain.Article{PublishedAt: time.Now()}, vnLoc(t))
	assert.Empty(t, item.ShortLink)
	assert.False(t, item.ShortLinkAlive)
	assert.Nil(t, item.UpdatedAt, "bài chưa từng cập nhật phải trả null")
	assert.NotNil(t, item.Tickers, "tickers phải là mảng rỗng, không phải null")
}

// Bài seed demo dùng URL không tồn tại nên short link phải tắt: thà hiện "Dữ
// liệu mẫu" trung thực còn hơn một link đẹp mà bấm vào ra 404.
func TestDemoArticleIsFlaggedAndLinkDisabled(t *testing.T) {
	item := toNewsItem(domain.Article{
		PublishedAt:    time.Now(),
		IsDemo:         true,
		ShortCode:      "a7Kx2p",
		ShortLinkAlive: false,
	}, vnLoc(t))

	assert.True(t, item.IsDemo)
	assert.False(t, item.ShortLinkAlive)
}

func TestEstimatedPublishDateIsFlagged(t *testing.T) {
	item := toNewsItem(domain.Article{
		PublishedAt:          time.Now(),
		PublishedAtEstimated: true,
	}, vnLoc(t))
	assert.True(t, item.PublishedAtEstimated, "FE hiển thị dấu ~ trước ngày")
}

func TestInvalidNewsTypeListsValidValues(t *testing.T) {
	values := domain.NewsTypeValues()
	for _, want := range []string{"vi_mo", "nganh", "phan_tich", "trai_phieu_tin_dung"} {
		assert.Contains(t, values, want)
	}
}
