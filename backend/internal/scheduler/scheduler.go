// Package scheduler chạy pipeline theo lịch cron (ADR-004: 06:00 / 14:00 /
// 22:00 giờ Việt Nam — trước phiên / trước ATC / sau phiên).
package scheduler

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/robfig/cron/v3"
)

// Job là việc cần chạy mỗi chu kỳ.
type Job func(ctx context.Context) error

// Scheduler bọc robfig/cron với múi giờ và logging structured.
type Scheduler struct {
	cron    *cron.Cron
	spec    string
	entryID cron.EntryID
	log     *slog.Logger
}

// New tạo scheduler chạy theo `spec` trong múi giờ `loc`.
// Dùng cron.WithLocation để 06:00 luôn là 06:00 giờ VN, không phụ thuộc TZ của
// tiến trình hay của host.
func New(spec string, loc *time.Location, job Job, log *slog.Logger) (*Scheduler, error) {
	if log == nil {
		log = slog.Default()
	}
	if loc == nil {
		loc = time.UTC
	}
	if job == nil {
		return nil, fmt.Errorf("scheduler requires a non-nil job")
	}

	c := cron.New(
		cron.WithLocation(loc),
		cron.WithLogger(cron.DiscardLogger),
		cron.WithChain(cron.SkipIfStillRunning(cron.DiscardLogger)),
	)

	s := &Scheduler{cron: c, spec: spec, log: log}
	id, err := c.AddFunc(spec, func() {
		start := time.Now()
		s.log.Info("cron tick", "spec", spec)
		// Mỗi lần chạy có context riêng, timeout rộng hơn NFR 10 phút.
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Minute)
		defer cancel()
		if err := job(ctx); err != nil {
			s.log.Error("cron job failed", "err", err,
				"duration", time.Since(start).String())
			return
		}
		s.log.Info("cron job finished", "duration", time.Since(start).String())
	})
	if err != nil {
		return nil, fmt.Errorf("register cron spec %q: %w", spec, err)
	}
	s.entryID = id
	return s, nil
}

// Start khởi động scheduler (không chặn).
func (s *Scheduler) Start() {
	s.cron.Start()
	s.log.Info("scheduler started", "spec", s.spec, "next_run", s.NextRun())
}

// Stop dừng scheduler và chờ job đang chạy kết thúc (tối đa theo ctx).
func (s *Scheduler) Stop(ctx context.Context) {
	stopCtx := s.cron.Stop()
	select {
	case <-stopCtx.Done():
		s.log.Info("scheduler stopped")
	case <-ctx.Done():
		s.log.Warn("scheduler stopped while a job was still running")
	}
}

// NextRun trả về thời điểm chạy kế tiếp (chuỗi rỗng nếu chưa start).
func (s *Scheduler) NextRun() string {
	entry := s.cron.Entry(s.entryID)
	if entry.Next.IsZero() {
		return ""
	}
	return entry.Next.Format(time.RFC3339)
}
