package store

import (
	"context"
	"fmt"
	"time"

	"github.com/tenpoint/tenpoint-api/internal/domain"
)

// StartRun mở một crawl_run mới và trả về id.
func (s *Store) StartRun(ctx context.Context, trigger string) (int64, error) {
	if trigger != domain.TriggerCron && trigger != domain.TriggerManual {
		return 0, fmt.Errorf("invalid trigger %q", trigger)
	}
	var id int64
	err := s.pool.QueryRow(ctx,
		"INSERT INTO crawl_runs (trigger) VALUES ($1) RETURNING id", trigger).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("open crawl_run: %w", err)
	}
	return id, nil
}

// FinishRun cập nhật tổng hợp khi run kết thúc.
func (s *Store) FinishRun(ctx context.Context, run domain.CrawlRun) error {
	_, err := s.pool.Exec(ctx, `
        UPDATE crawl_runs
        SET finished_at = now(), articles_found = $2, articles_new = $3,
            articles_rejected = $4, errors = $5
        WHERE id = $1`,
		run.ID, run.ArticlesFound, run.ArticlesNew, run.ArticlesRejected, run.Errors)
	if err != nil {
		return fmt.Errorf("close crawl_run %d: %w", run.ID, err)
	}
	return nil
}

// RecordRunSource ghi kết quả xử lý của một nguồn trong run.
// Một nguồn lỗi KHÔNG làm hỏng cả run — lỗi được ghi vào error_message.
func (s *Store) RecordRunSource(ctx context.Context, rs domain.CrawlRunSource) error {
	_, err := s.pool.Exec(ctx, `
        INSERT INTO crawl_run_sources (run_id, source_id, found, new_count, error_message, truncated, off_topic)
        VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		rs.RunID, rs.SourceID, rs.Found, rs.NewCount, rs.ErrorMessage, rs.Truncated, rs.OffTopic)
	if err != nil {
		return fmt.Errorf("insert crawl_run_sources (run=%d, source=%d): %w", rs.RunID, rs.SourceID, err)
	}
	return nil
}

// LastCrawlAt trả về thời điểm run gần nhất kết thúc (nil nếu chưa có run nào).
func (s *Store) LastCrawlAt(ctx context.Context) (*time.Time, error) {
	var t *time.Time
	err := s.pool.QueryRow(ctx,
		"SELECT max(finished_at) FROM crawl_runs WHERE finished_at IS NOT NULL").Scan(&t)
	if err != nil {
		return nil, fmt.Errorf("read last crawl timestamp: %w", err)
	}
	return t, nil
}
