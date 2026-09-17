package shortlink

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

// FlushInterval — gộp UPDATE click_count mỗi 5 giây (mục 7 architecture.md).
const FlushInterval = 5 * time.Second

// FlushFunc ghi batch click xuống DB.
type FlushFunc func(ctx context.Context, counts map[string]int64) error

// ClickCounter đếm click trong RAM rồi flush theo batch.
// Add() KHÔNG BAO GIỜ chặn redirect — lỗi ghi log không được ảnh hưởng người dùng.
type ClickCounter struct {
	mu     sync.Mutex
	counts map[string]int64

	flush    FlushFunc
	interval time.Duration
	log      *slog.Logger
}

// NewClickCounter tạo bộ đếm.
func NewClickCounter(flush FlushFunc, interval time.Duration, log *slog.Logger) *ClickCounter {
	if interval <= 0 {
		interval = FlushInterval
	}
	if log == nil {
		log = slog.Default()
	}
	return &ClickCounter{
		counts:   make(map[string]int64),
		flush:    flush,
		interval: interval,
		log:      log,
	}
}

// Add ghi nhận 1 click. Thao tác chỉ khoá mutex trong RAM.
func (c *ClickCounter) Add(code string) {
	if code == "" {
		return
	}
	c.mu.Lock()
	c.counts[code]++
	c.mu.Unlock()
}

// Run chạy vòng flush cho tới khi ctx bị huỷ, sau đó flush lần cuối.
func (c *ClickCounter) Run(ctx context.Context) {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			// Flush lần cuối với context độc lập để không mất số liệu khi shutdown.
			final, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
			c.Flush(final)
			cancel()
			return
		case <-ticker.C:
			c.Flush(ctx)
		}
	}
}

// Flush đẩy toàn bộ bộ đếm hiện tại xuống DB.
func (c *ClickCounter) Flush(ctx context.Context) {
	c.mu.Lock()
	if len(c.counts) == 0 {
		c.mu.Unlock()
		return
	}
	batch := c.counts
	c.counts = make(map[string]int64, len(batch))
	c.mu.Unlock()

	if c.flush == nil {
		return
	}
	if err := c.flush(ctx, batch); err != nil {
		// Không nuốt lỗi im lặng, nhưng cũng không làm hỏng luồng redirect.
		c.log.Error("flush click_count failed", "err", err, "codes", len(batch))
	}
}

// Pending trả về số mã đang chờ flush (dùng cho test và metric).
func (c *ClickCounter) Pending() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.counts)
}
