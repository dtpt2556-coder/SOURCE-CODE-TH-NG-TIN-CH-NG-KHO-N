// Package fetcher tải nội dung từ nguồn tin: RSS, HTML, robots.txt.
// Luôn tôn trọng robots.txt, rate limit theo domain và dùng User-Agent định danh
// rõ ràng (R2.1/R2.2 của BA2) — không bao giờ giả mạo UA trình duyệt.
package fetcher

import (
	"context"
	"sync"
	"time"
)

// RateLimiter giãn cách request theo từng host.
type RateLimiter struct {
	mu       sync.Mutex
	next     map[string]time.Time
	fallback time.Duration
}

// NewRateLimiter tạo limiter với khoảng cách mặc định khi nguồn không khai báo.
func NewRateLimiter(fallback time.Duration) *RateLimiter {
	if fallback <= 0 {
		fallback = time.Second
	}
	return &RateLimiter{next: make(map[string]time.Time), fallback: fallback}
}

// Wait chặn cho tới khi được phép gọi host tiếp theo, hoặc ctx bị huỷ.
func (r *RateLimiter) Wait(ctx context.Context, host string, interval time.Duration) error {
	if interval <= 0 {
		interval = r.fallback
	}

	r.mu.Lock()
	now := time.Now()
	allowed, ok := r.next[host]
	if !ok || allowed.Before(now) {
		allowed = now
	}
	r.next[host] = allowed.Add(interval)
	r.mu.Unlock()

	delay := time.Until(allowed)
	if delay <= 0 {
		return ctx.Err()
	}
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
