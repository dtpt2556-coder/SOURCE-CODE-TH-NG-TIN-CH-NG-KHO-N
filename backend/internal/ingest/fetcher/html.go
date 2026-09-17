package fetcher

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// MaxBodyBytes chặn bài quá lớn làm nổ bộ nhớ.
const MaxBodyBytes = 4 << 20 // 4 MiB

// MaxRedirects — R2.4 của BA2: tối đa 3 hop.
const MaxRedirects = 3

// Page là kết quả tải một trang.
type Page struct {
	FinalURL   string
	StatusCode int
	Body       string
}

// HTTPFetcher tải HTML/RSS với User-Agent định danh, timeout và rate limit.
type HTTPFetcher struct {
	client    *http.Client
	userAgent string
	limiter   *RateLimiter
}

// NewHTTPFetcher dựng fetcher.
func NewHTTPFetcher(userAgent string, timeout time.Duration, limiter *RateLimiter) *HTTPFetcher {
	if timeout <= 0 {
		timeout = 20 * time.Second
	}
	if limiter == nil {
		limiter = NewRateLimiter(time.Second)
	}
	client := &http.Client{
		Timeout: timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= MaxRedirects {
				return fmt.Errorf("exceeded %d redirects", MaxRedirects)
			}
			return nil
		},
	}
	return &HTTPFetcher{client: client, userAgent: userAgent, limiter: limiter}
}

// UserAgent trả về UA đang dùng.
func (f *HTTPFetcher) UserAgent() string { return f.userAgent }

// Get tải một URL, có rate limit theo host.
func (f *HTTPFetcher) Get(ctx context.Context, rawURL string, interval time.Duration) (Page, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return Page{}, fmt.Errorf("parse URL %q: %w", rawURL, err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return Page{}, fmt.Errorf("scheme %q is not accepted", u.Scheme)
	}
	if err := f.limiter.Wait(ctx, strings.ToLower(u.Host), interval); err != nil {
		return Page{}, fmt.Errorf("wait on rate limit for %s: %w", u.Host, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return Page{}, fmt.Errorf("build request %s: %w", rawURL, err)
	}
	req.Header.Set("User-Agent", f.userAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "vi,en;q=0.8")

	resp, err := f.client.Do(req)
	if err != nil {
		return Page{}, fmt.Errorf("fetch %s: %w", rawURL, err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(io.LimitReader(resp.Body, MaxBodyBytes))
	if err != nil {
		return Page{}, fmt.Errorf("read body of %s: %w", rawURL, err)
	}
	page := Page{
		FinalURL:   resp.Request.URL.String(),
		StatusCode: resp.StatusCode,
		Body:       string(body),
	}
	if resp.StatusCode != http.StatusOK {
		return page, fmt.Errorf("%s returned HTTP %d", rawURL, resp.StatusCode)
	}
	return page, nil
}

// GetString là helper cho RobotsCache (chỉ cần phần thân).
func (f *HTTPFetcher) GetString(ctx context.Context, rawURL string) (string, error) {
	p, err := f.Get(ctx, rawURL, 0)
	if err != nil {
		return "", err
	}
	return p.Body, nil
}
