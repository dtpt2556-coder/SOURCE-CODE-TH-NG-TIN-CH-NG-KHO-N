package extractor_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tenpoint/tenpoint-api/internal/ingest/extractor"
	"github.com/tenpoint/tenpoint-api/internal/ingest/tagger"
)

// testdataDir trỏ tới backend/testdata (dùng chung cho extractor và tagger).
const testdataDir = "../../../testdata"

func loadHTML(t *testing.T, name string) string {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join(testdataDir, name))
	require.NoError(t, err, "đọc %s", name)
	return string(raw)
}

func vnLocation(t *testing.T) *time.Location {
	t.Helper()
	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	require.NoError(t, err)
	return loc
}

func TestExtractVnExpressArticle(t *testing.T) {
	ex := extractor.New(vnLocation(t))
	art, err := ex.Extract(loadHTML(t, "vnexpress_article.html"),
		"https://vnexpress.net/12-ngan-hang-cam-ket-408000-ty-dong-tin-dung-cho-dnnvv-4790001.html")
	require.NoError(t, err)

	assert.Equal(t, "12 ngân hàng cam kết 408.000 tỷ đồng tín dụng cho DNNVV", art.Title)
	assert.Contains(t, art.Body, "408.000 tỷ đồng")
	assert.Contains(t, art.Body, "220.000 tỷ đồng")
	assert.Contains(t, art.Body, "8,8%")

	require.True(t, art.HasDate, "phải đọc được article:published_time")
	// 2026-08-26T08:30:00+07:00 == 2026-08-26T01:30:00Z
	assert.Equal(t, "2026-08-26T01:30:00Z", art.PublishedAt.UTC().Format(time.RFC3339))
}

func TestExtractRemovesNoiseBlocks(t *testing.T) {
	ex := extractor.New(vnLocation(t))
	art, err := ex.Extract(loadHTML(t, "vnexpress_article.html"), "https://vnexpress.net/x.html")
	require.NoError(t, err)

	assert.NotContains(t, art.Body, "Tin liên quan",
		"khối 'Tin liên quan' phải bị loại, nếu không tóm tắt sẽ lẫn nội dung bài khác (F3.c)")
	assert.NotContains(t, art.Body, "Bản quyền thuộc VnExpress")
	assert.NotContains(t, art.Body, "Menu điều hướng")
	assert.NotContains(t, art.Body, "var tracking")
}

func TestExtractCafeFArticleWithLocalTime(t *testing.T) {
	ex := extractor.New(vnLocation(t))
	art, err := ex.Extract(loadHTML(t, "cafef_article.html"), "https://cafef.vn/vn-index.chn")
	require.NoError(t, err)

	assert.Equal(t, "VN-Index đóng cửa 1.821 điểm, lần đầu vượt mốc 1.800", art.Title)
	assert.Contains(t, art.Body, "1.821 điểm")
	assert.Contains(t, art.Body, "39,17 triệu cp")
	assert.NotContains(t, art.Body, "Quảng cáo")

	require.True(t, art.HasDate)
	// Nguồn VN không ghi offset -> hiểu theo ICT: 15:45 ICT == 08:45 UTC.
	assert.Equal(t, "2026-08-26T08:45:00Z", art.PublishedAt.UTC().Format(time.RFC3339))
}

func TestExtractRejectsPaywalledArticle(t *testing.T) {
	ex := extractor.New(vnLocation(t))
	_, err := ex.Extract(loadHTML(t, "paywalled_article.html"), "https://vneconomy.vn/bao-cao.htm")
	require.Error(t, err, "bài quá ngắn phải bị loại (R2.5/R3.7)")
	assert.Contains(t, err.Error(), "paywalled")
}

func TestExtractRejectsHTMLWithoutTitle(t *testing.T) {
	ex := extractor.New(vnLocation(t))
	_, err := ex.Extract("<html><body><p>"+strings.Repeat("nội dung dài ", 60)+"</p></body></html>",
		"https://cafef.vn/x.chn")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no title")
}

// Extractor + tagger chạy nối tiếp trên cùng file testdata: kết quả phải khớp
// đúng các mã mà ảnh tham chiếu hiển thị.
func TestExtractThenTagUsesTestdata(t *testing.T) {
	ex := extractor.New(vnLocation(t))
	art, err := ex.Extract(loadHTML(t, "cafef_article.html"), "https://cafef.vn/vn-index.chn")
	require.NoError(t, err)

	tg := tagger.New([]tagger.Ticker{
		{Symbol: "TCB", Sector: "Ngân hàng", Aliases: []string{"Techcombank"}},
		{Symbol: "VIC", Sector: "Bất động sản", Aliases: []string{"Vingroup"}},
		{Symbol: "FPT", Sector: "Công nghệ", Aliases: []string{"Tập đoàn FPT"}},
		{Symbol: "BCM", Sector: "BĐS khu công nghiệp", Aliases: []string{"Becamex"}},
	})
	got := tg.Tag(tagger.Input{Title: art.Title, Body: art.Body})

	found := map[string]bool{}
	for _, m := range got {
		found[m.Symbol] = true
	}
	for _, sym := range []string{"TCB", "BCM", "VIC", "FPT"} {
		assert.True(t, found[sym], "phải gắn %s từ nội dung đã extract, nhận được %v", sym, found)
	}
}
