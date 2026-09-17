package shortlink

import (
	"fmt"
	"net/url"
	"strings"
	"sync"
)

// ErrNotAllowed được bọc trong lỗi trả về khi domain đích ngoài allowlist.
var ErrNotAllowed = fmt.Errorf("target domain is not in the allowlist")

// Allowlist chặn open redirect: chỉ cho phép tạo short link trỏ tới domain đã
// duyệt trong bảng `sources` (hoặc subdomain của nó).
//
// Không có biện pháp này, /r/{code} trở thành công cụ phishing — đây là yêu cầu
// bảo mật bắt buộc (BA1 US-3.3, §7.4 BA2).
type Allowlist struct {
	mu      sync.RWMutex
	domains map[string]bool
}

// NewAllowlist dựng allowlist từ danh sách domain của các nguồn tin.
func NewAllowlist(domains []string) *Allowlist {
	a := &Allowlist{domains: make(map[string]bool, len(domains))}
	a.Replace(domains)
	return a
}

// Replace thay toàn bộ allowlist (gọi sau mỗi lần nạp lại `sources`).
func (a *Allowlist) Replace(domains []string) {
	next := make(map[string]bool, len(domains))
	for _, d := range domains {
		d = normalizeHost(d)
		if d != "" {
			next[d] = true
		}
	}
	a.mu.Lock()
	a.domains = next
	a.mu.Unlock()
}

// Size trả về số domain đang cho phép.
func (a *Allowlist) Size() int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return len(a.domains)
}

// Check xác thực URL đích. Trả về nil nếu hợp lệ.
//
// Ràng buộc:
//   - scheme chỉ được http/https (chặn javascript:, data:, file:, //evil.com)
//   - host phải khớp chính xác một domain trong allowlist, hoặc là subdomain
//     của nó ("finance.vietstock.vn" khớp "vietstock.vn")
func (a *Allowlist) Check(rawURL string) error {
	raw := strings.TrimSpace(rawURL)
	if raw == "" {
		return fmt.Errorf("empty target URL")
	}
	// URL protocol-relative ("//evil.com/x") phải bị từ chối thẳng.
	if strings.HasPrefix(raw, "//") {
		return fmt.Errorf("protocol-relative URL %q: %w", raw, ErrNotAllowed)
	}
	u, err := url.Parse(raw)
	if err != nil {
		return fmt.Errorf("parse target URL %q: %w", raw, err)
	}
	scheme := strings.ToLower(u.Scheme)
	if scheme != "http" && scheme != "https" {
		return fmt.Errorf("scheme %q is not accepted: %w", u.Scheme, ErrNotAllowed)
	}
	host := normalizeHost(u.Hostname())
	if host == "" {
		return fmt.Errorf("URL %q has no host: %w", raw, ErrNotAllowed)
	}

	a.mu.RLock()
	defer a.mu.RUnlock()
	if a.domains[host] {
		return nil
	}
	for d := range a.domains {
		if strings.HasSuffix(host, "."+d) {
			return nil
		}
	}
	return fmt.Errorf("host %q: %w", host, ErrNotAllowed)
}

// Allowed là dạng boolean tiện dùng trong điều kiện.
func (a *Allowlist) Allowed(rawURL string) bool {
	return a.Check(rawURL) == nil
}

func normalizeHost(h string) string {
	h = strings.ToLower(strings.TrimSpace(h))
	h = strings.TrimPrefix(h, "www.")
	h = strings.TrimSuffix(h, ".")
	return h
}
