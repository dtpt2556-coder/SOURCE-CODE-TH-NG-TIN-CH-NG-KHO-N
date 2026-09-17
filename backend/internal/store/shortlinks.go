package store

import (
	"context"
	"fmt"

	"github.com/tenpoint/tenpoint-api/internal/domain"
)

// GetShortLink đọc đích redirect theo mã.
func (s *Store) GetShortLink(ctx context.Context, code string) (domain.ShortLink, error) {
	var l domain.ShortLink
	err := s.pool.QueryRow(ctx, `
        SELECT id, code, article_id, target_url, click_count, last_checked_at, target_alive
        FROM short_links WHERE code = $1`, code).
		Scan(&l.ID, &l.Code, &l.ArticleID, &l.TargetURL, &l.ClickCount, &l.LastCheckedAt, &l.TargetAlive)
	if err != nil {
		return domain.ShortLink{}, wrapNoRows(err, fmt.Sprintf("short link %s", code))
	}
	return l, nil
}

// FlushClicks cộng dồn click theo batch — được ClickCounter gọi mỗi 5 giây.
// KHÔNG chạy trên đường redirect nên không ảnh hưởng độ trễ người dùng.
func (s *Store) FlushClicks(ctx context.Context, counts map[string]int64) error {
	if len(counts) == 0 {
		return nil
	}
	codes := make([]string, 0, len(counts))
	deltas := make([]int64, 0, len(counts))
	for code, n := range counts {
		codes = append(codes, code)
		deltas = append(deltas, n)
	}
	_, err := s.pool.Exec(ctx, `
        UPDATE short_links sl
        SET click_count = sl.click_count + v.delta
        FROM unnest($1::text[], $2::bigint[]) AS v(code, delta)
        WHERE sl.code = v.code`, codes, deltas)
	if err != nil {
		return fmt.Errorf("increment click_count for %d codes: %w", len(codes), err)
	}
	return nil
}

// MarkShortLinkDead vô hiệu hoá short link khi bài gốc không còn truy cập được
// (F2.e của BA2) — /r/{code} sẽ trả 410 thay vì 302 tới link vỡ.
func (s *Store) MarkShortLinkDead(ctx context.Context, code string) error {
	_, err := s.pool.Exec(ctx, `
        UPDATE short_links SET target_alive = false, last_checked_at = now()
        WHERE code = $1`, code)
	if err != nil {
		return fmt.Errorf("mark short link %s dead: %w", code, err)
	}
	return nil
}
