package shortlink_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tenpoint/tenpoint-api/internal/shortlink"
)

func testAllowlist() *shortlink.Allowlist {
	return shortlink.NewAllowlist([]string{
		"vnexpress.net", "cafef.vn", "thanhnien.vn", "vietstock.vn",
		"tinnhanhchungkhoan.vn", "hsx.vn",
	})
}

func TestAllowlistAcceptsKnownSourceDomains(t *testing.T) {
	a := testAllowlist()
	ok := []string{
		"https://vnexpress.net/12-ngan-hang-cam-ket-408000-ty-dong-tin-dung-cho-dnnvv-4790001.html",
		"https://cafef.vn/vn-index-dong-cua-1821-diem.chn",
		"http://thanhnien.vn/bai-viet.htm",
		"https://www.vnexpress.net/bai-viet.html",
		// Subdomain của domain đã duyệt vẫn hợp lệ.
		"https://finance.vietstock.vn/bao-cao.htm",
	}
	for _, u := range ok {
		assert.NoError(t, a.Check(u), u)
		assert.True(t, a.Allowed(u), u)
	}
}

func TestAllowlistRejectsUnknownDomain(t *testing.T) {
	a := testAllowlist()
	bad := []string{
		"https://evil.com/phishing",
		"https://vnexpress.net.evil.com/phishing",
		"https://notcafef.vn/x",
		"https://cafef.vn.attacker.io/x",
	}
	for _, u := range bad {
		err := a.Check(u)
		require.Error(t, err, u)
		assert.True(t, errors.Is(err, shortlink.ErrNotAllowed), "phải là ErrNotAllowed: %v", err)
		assert.False(t, a.Allowed(u), u)
	}
}

func TestAllowlistRejectsDangerousSchemes(t *testing.T) {
	a := testAllowlist()
	bad := []string{
		"javascript:alert(document.cookie)",
		"data:text/html;base64,PHNjcmlwdD4=",
		"file:///etc/passwd",
		"//evil.com/phishing", // protocol-relative
		"ftp://cafef.vn/x",
		"",
		"   ",
	}
	for _, u := range bad {
		assert.Error(t, a.Check(u), "phải từ chối %q", u)
	}
}

func TestAllowlistReplaceUpdatesRules(t *testing.T) {
	a := shortlink.NewAllowlist([]string{"cafef.vn"})
	require.NoError(t, a.Check("https://cafef.vn/x"))
	assert.Error(t, a.Check("https://vnexpress.net/x"))
	assert.Equal(t, 1, a.Size())

	a.Replace([]string{"vnexpress.net"})
	assert.Error(t, a.Check("https://cafef.vn/x"))
	assert.NoError(t, a.Check("https://vnexpress.net/x"))
	assert.Equal(t, 1, a.Size())
}

func TestEmptyAllowlistRejectsEverything(t *testing.T) {
	a := shortlink.NewAllowlist(nil)
	assert.Error(t, a.Check("https://vnexpress.net/x"),
		"allowlist rỗng phải từ chối tất cả — mặc định an toàn")
}
