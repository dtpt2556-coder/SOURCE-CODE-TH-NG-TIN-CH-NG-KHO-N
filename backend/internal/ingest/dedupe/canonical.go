// Package dedupe khử trùng lặp bài viết: chuẩn hoá URL (L0) và SimHash tiêu đề
// (L2, cửa sổ 72h) theo §6 BA2 và ADR-009.
package dedupe

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"
)

// trackingParams là các query param quảng cáo/đo lường cần loại bỏ (E-02).
// Mọi param có tiền tố "utm_" đều bị loại, danh sách này là phần còn lại.
var trackingParams = map[string]bool{
	"fbclid": true, "gclid": true, "zarsrc": true, "source": true, "src": true,
	"ref": true, "campaign": true, "cmp": true, "spm": true,
	"amp": true, "__twitter_impression": true, "utm_id": true,
}

// CanonicalURL chuẩn hoá URL theo §6.1 BA2:
// lowercase scheme+host, bỏ "www.", bỏ tracking param, bỏ fragment, bỏ hậu tố
// AMP, bỏ trailing slash. Trả lỗi nếu URL không hợp lệ hoặc scheme không phải
// http/https.
func CanonicalURL(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", fmt.Errorf("empty URL")
	}
	u, err := url.Parse(raw)
	if err != nil {
		return "", fmt.Errorf("parse URL %q: %w", raw, err)
	}
	u.Scheme = strings.ToLower(u.Scheme)
	if u.Scheme != "http" && u.Scheme != "https" {
		return "", fmt.Errorf("scheme %q is not accepted (http/https only)", u.Scheme)
	}
	if u.Host == "" {
		return "", fmt.Errorf("URL %q has no host", raw)
	}
	// Ép https để http:// và https:// của cùng một bài cho ra cùng url_hash.
	// Mọi nguồn trong allowlist đều phục vụ https (E-04).
	u.Scheme = "https"
	u.Host = strings.TrimPrefix(strings.ToLower(u.Host), "www.")
	u.Fragment = ""
	u.RawFragment = ""
	u.User = nil

	q := u.Query()
	for k := range q {
		lower := strings.ToLower(k)
		if trackingParams[lower] || strings.HasPrefix(lower, "utm_") {
			q.Del(k)
		}
	}
	u.RawQuery = q.Encode()

	p := u.Path
	// Bỏ hậu tố AMP.
	for _, suffix := range []string{"/amp", "/amp/", ".amp"} {
		if strings.HasSuffix(strings.ToLower(p), suffix) {
			p = p[:len(p)-len(suffix)]
			break
		}
	}
	if len(p) > 1 {
		p = strings.TrimRight(p, "/")
	}
	u.Path = p

	return u.String(), nil
}

// URLHash = SHA256 hex của canonical URL — khoá UNIQUE chặn trùng khi ingest.
func URLHash(canonical string) string {
	sum := sha256.Sum256([]byte(canonical))
	return hex.EncodeToString(sum[:])
}

// DisplayDomain trả về host hiển thị ở cột "Source" (đã bỏ "www.").
func DisplayDomain(raw string) string {
	u, err := url.Parse(raw)
	if err != nil || u.Host == "" {
		return ""
	}
	return strings.TrimPrefix(strings.ToLower(u.Hostname()), "www.")
}
