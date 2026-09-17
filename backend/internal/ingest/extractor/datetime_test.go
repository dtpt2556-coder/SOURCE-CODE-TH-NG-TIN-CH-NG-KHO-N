package extractor_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tenpoint/tenpoint-api/internal/ingest/extractor"
)

// now là mốc cố định để "2 giờ trước" cho ra kết quả tất định.
func fixedNow(t *testing.T) time.Time {
	t.Helper()
	return time.Date(2026, 8, 26, 10, 0, 0, 0, vnLocation(t))
}

func TestParsePublishedAtVietnameseFormats(t *testing.T) {
	loc := vnLocation(t)
	now := fixedNow(t)

	cases := []struct {
		id   string
		raw  string
		want string
	}{
		{"T-01", "Thứ Tư, 26/8/2026, 08:15 (GMT+7)", "2026-08-26T01:15:00Z"},
		{"T-01", "Thứ Ba, 26/08/2026 08:15", "2026-08-26T01:15:00Z"},
		{"T-02", "26/08/2026 08:15", "2026-08-26T01:15:00Z"},
		{"T-02", "26/08/2026", "2026-08-25T17:00:00Z"},
		{"T-03", "2 giờ trước", "2026-08-26T01:00:00Z"},
		{"T-03", "30 phút trước", "2026-08-26T02:30:00Z"},
		{"T-03", "3 ngày trước", "2026-08-23T03:00:00Z"},
		{"T-03", "hôm qua, 14:30", "2026-08-25T07:30:00Z"},
		{"T-04", "Wed, 26 Aug 2026 08:15:00 +0700", "2026-08-26T01:15:00Z"},
		{"T-05", "2026-08-26T08:15:00+07:00", "2026-08-26T01:15:00Z"},
		{"T-05", "2026-08-26T01:15:00Z", "2026-08-26T01:15:00Z"},
		{"", "2026-08-26 08:15:00", "2026-08-26T01:15:00Z"},
	}

	for _, tc := range cases {
		t.Run(tc.id+" "+tc.raw, func(t *testing.T) {
			got, ok := extractor.ParsePublishedAt(tc.raw, now, loc)
			require.True(t, ok, "phải parse được %q", tc.raw)
			assert.Equal(t, tc.want, got.UTC().Format(time.RFC3339))
		})
	}
}

func TestParsePublishedAtRejectsGarbage(t *testing.T) {
	loc := vnLocation(t)
	now := fixedNow(t)

	for _, raw := range []string{"", "   ", "cập nhật mới nhất", "32/13/2026", "31/02/2026"} {
		_, ok := extractor.ParsePublishedAt(raw, now, loc)
		assert.False(t, ok, "không được parse %q", raw)
	}
}

// T-06: không có thời gian nào thì dùng fetched_at và bật cờ ước lượng.
func TestSettlePublishedAtFallsBackToFetchedAt(t *testing.T) {
	fetchedAt := time.Date(2026, 8, 26, 3, 0, 0, 0, time.UTC)

	got, estimated := extractor.SettlePublishedAt(time.Time{}, false, fetchedAt)
	assert.True(t, estimated)
	assert.Equal(t, fetchedAt, got)
}

// T-07: CMS trả ngày ở tương lai thì kẹp về fetched_at.
func TestSettlePublishedAtClampsFutureTimestamps(t *testing.T) {
	fetchedAt := time.Date(2026, 8, 26, 3, 0, 0, 0, time.UTC)

	future := fetchedAt.Add(5 * time.Hour)
	got, estimated := extractor.SettlePublishedAt(future, true, fetchedAt)
	assert.True(t, estimated, "thời gian tương lai phải bị kẹp")
	assert.Equal(t, fetchedAt, got)

	// Lệch trong ngưỡng dung sai thì vẫn tin nguồn.
	slight := fetchedAt.Add(time.Hour)
	got, estimated = extractor.SettlePublishedAt(slight, true, fetchedAt)
	assert.False(t, estimated)
	assert.Equal(t, slight, got)
}

// T-08: bài quá cũ lọt vào feed "mới nhất".
func TestTooOld(t *testing.T) {
	now := time.Date(2026, 8, 26, 3, 0, 0, 0, time.UTC)
	assert.True(t, extractor.TooOld(now.AddDate(0, 0, -100), now))
	assert.False(t, extractor.TooOld(now.AddDate(0, 0, -10), now))
}

// T-09: bài đăng 23:50 ngày 26/08 giờ VN phải hiển thị 26/08, không phải 27/08.
func TestLateNightArticleKeepsVietnamCalendarDay(t *testing.T) {
	loc := vnLocation(t)
	now := fixedNow(t)

	got, ok := extractor.ParsePublishedAt("26/08/2026 23:50", now, loc)
	require.True(t, ok)

	assert.Equal(t, "2026-08-26T16:50:00Z", got.UTC().Format(time.RFC3339),
		"lưu trữ theo UTC")
	assert.Equal(t, "26/08/2026", got.In(loc).Format("02/01/2006"),
		"hiển thị phải đổi sang giờ VN TRƯỚC khi cắt ngày")
	assert.Equal(t, "26/08/2026", got.In(loc).Format("02/01/2006"))
	assert.NotEqual(t, "27/08/2026", got.UTC().Format("02/01/2006"))
}

func TestExtractorReadsVietnameseVisibleDate(t *testing.T) {
	ex := extractor.New(vnLocation(t))
	ex.Now = func() time.Time { return fixedNow(t) }

	html := `<html><head><meta property="og:title" content="Tiêu đề bài viết thử nghiệm"></head>
	<body><article><span class="date">Thứ Tư, 26/8/2026, 08:15 (GMT+7)</span>
	<div class="detail-content">` + longParagraphs() + `</div></article></body></html>`

	got, err := ex.Extract(html, "https://cafef.vn/bai-viet.chn")
	require.NoError(t, err)
	require.True(t, got.HasDate, "phải đọc được ngày từ text hiển thị")
	assert.Equal(t, "2026-08-26T01:15:00Z", got.PublishedAt.UTC().Format(time.RFC3339))
}

func longParagraphs() string {
	var out strings.Builder
	for i := 0; i < 8; i++ {
		fmt.Fprintf(&out,
			"<p>Phiên thứ %d: VN-Index đóng cửa 1.821 điểm, tăng 29,91 điểm tương đương 1,67%%, "+
				"thanh khoản toàn thị trường đạt khoảng 20.000 tỷ đồng với 182 mã tăng giá.</p>", i+1)
	}
	return out.String()
}
