package store

import (
	"context"
	"fmt"

	"github.com/tenpoint/tenpoint-api/internal/domain"
)

const sourceColumns = `id, code, name, domain, rss_url, list_url, tier, enabled, rate_limit_ms`

// ListEnabledSources trả về các nguồn đang bật, ưu tiên tier thấp (sơ cấp) trước.
func (s *Store) ListEnabledSources(ctx context.Context) ([]domain.Source, error) {
	return s.querySources(ctx, `SELECT `+sourceColumns+`
        FROM sources WHERE enabled = true ORDER BY tier, code`)
}

// ListSources trả về toàn bộ nguồn (kể cả nguồn tắt) — dùng để dựng allowlist.
func (s *Store) ListSources(ctx context.Context) ([]domain.Source, error) {
	return s.querySources(ctx, `SELECT `+sourceColumns+` FROM sources ORDER BY tier, code`)
}

// ListSourceDomains trả về danh sách domain cho allowlist chống open redirect.
func (s *Store) ListSourceDomains(ctx context.Context) ([]string, error) {
	rows, err := s.pool.Query(ctx, "SELECT domain FROM sources")
	if err != nil {
		return nil, fmt.Errorf("query source domains: %w", err)
	}
	defer rows.Close()

	var out []string
	for rows.Next() {
		var d string
		if err := rows.Scan(&d); err != nil {
			return nil, fmt.Errorf("scan source domain: %w", err)
		}
		out = append(out, d)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate source domains: %w", err)
	}
	return out, nil
}

func (s *Store) querySources(ctx context.Context, q string) ([]domain.Source, error) {
	rows, err := s.pool.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("query sources: %w", err)
	}
	defer rows.Close()

	var out []domain.Source
	for rows.Next() {
		var src domain.Source
		if err := rows.Scan(&src.ID, &src.Code, &src.Name, &src.Domain, &src.RSSURL,
			&src.ListURL, &src.Tier, &src.Enabled, &src.RateLimitMS); err != nil {
			return nil, fmt.Errorf("scan source: %w", err)
		}
		out = append(out, src)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate sources: %w", err)
	}
	return out, nil
}
