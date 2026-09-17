package dedupe_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tenpoint/tenpoint-api/internal/ingest/dedupe"
)

// E-02: mọi biến thể tham số theo dõi của cùng một bài phải cho ra một url_hash.
// Không có rule này, một bài chia sẻ qua Zalo, Facebook và newsletter sẽ thành
// ba dòng khác nhau trên bảng digest.
func TestTrackingParametersCollapseToOneHash(t *testing.T) {
	bare := "https://vnexpress.net/12-ngan-hang-cam-ket-408000-ty-dong-4790001.html"
	want := dedupe.URLHash(mustCanonical(t, bare))

	variants := []string{
		bare + "?utm_source=zalo",
		bare + "?utm_source=facebook&utm_medium=social&utm_campaign=tet",
		bare + "?fbclid=IwAR0abcdef",
		bare + "?gclid=EAIaIQobChMI",
		bare + "?zarsrc=30",
		bare + "?src=newsletter",
		bare + "?ref=home",
		bare + "#box-comment",
		bare + "?utm_source=zalo#comment",
	}
	for _, raw := range variants {
		assert.Equal(t, want, dedupe.URLHash(mustCanonical(t, raw)), "biến thể %q", raw)
	}
}

// E-02: tham số thật sự thuộc về bài thì KHÔNG được xoá.
func TestMeaningfulQueryParametersSurvive(t *testing.T) {
	got := mustCanonical(t, "https://vietstock.vn/bai-viet.htm?page=2&utm_source=zalo")
	assert.Contains(t, got, "page=2")
	assert.NotContains(t, got, "utm_source")
}

// E-04: http/https, www, dấu "/" cuối, hoa thường phần host.
func TestSchemeHostAndTrailingSlashNormalisation(t *testing.T) {
	want := mustCanonical(t, "https://cafef.vn/bai-viet.chn")

	for _, raw := range []string{
		"http://cafef.vn/bai-viet.chn",
		"https://www.cafef.vn/bai-viet.chn",
		"http://WWW.CafeF.VN/bai-viet.chn",
		"https://CAFEF.vn/bai-viet.chn",
	} {
		assert.Equal(t, want, mustCanonical(t, raw), "biến thể %q", raw)
	}
}

// E-04: phần path phân biệt hoa thường — nhiều CMS coi đó là hai bài khác nhau.
func TestPathCaseIsPreserved(t *testing.T) {
	got := mustCanonical(t, "https://cafef.vn/Bai-Viet-Quan-Trong.chn")
	assert.Contains(t, got, "Bai-Viet-Quan-Trong.chn")
}

func TestAmpSuffixRemoved(t *testing.T) {
	want := mustCanonical(t, "https://thanhnien.vn/bai-viet")
	for _, raw := range []string{
		"https://thanhnien.vn/bai-viet/amp",
		"https://thanhnien.vn/bai-viet/amp/",
		"https://thanhnien.vn/bai-viet?amp=1",
	} {
		assert.Equal(t, want, mustCanonical(t, raw), "biến thể %q", raw)
	}
}

func TestCanonicalRejectsUnsafeSchemes(t *testing.T) {
	for _, raw := range []string{"javascript:alert(1)", "data:text/html,x", "file:///etc/passwd", ""} {
		_, err := dedupe.CanonicalURL(raw)
		assert.Error(t, err, raw)
	}
}

func TestDisplayDomain(t *testing.T) {
	assert.Equal(t, "vnexpress.net", dedupe.DisplayDomain("https://www.vnexpress.net/bai-viet.html"))
	assert.Equal(t, "cafef.vn", dedupe.DisplayDomain("https://cafef.vn/x.chn"))
	assert.Empty(t, dedupe.DisplayDomain("không-phải-url"))
}

func mustCanonical(t *testing.T, raw string) string {
	t.Helper()
	got, err := dedupe.CanonicalURL(raw)
	require.NoError(t, err, raw)
	return got
}
