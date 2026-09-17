package store

import (
	"context"
	"fmt"

	"github.com/tenpoint/tenpoint-api/internal/domain"
)

// ListTickers trả về danh sách mã kèm số tin đã gắn (cho ô lọc/autocomplete).
func (s *Store) ListTickers(ctx context.Context) ([]domain.Ticker, error) {
	rows, err := s.pool.Query(ctx, `
        SELECT t.symbol, t.company_name, t.short_name, t.exchange, t.sector,
               t.in_vn30, count(at.article_id)
        FROM tickers t
        LEFT JOIN article_tickers at ON at.symbol = t.symbol
        LEFT JOIN articles a ON a.id = at.article_id AND a.status = 'published'
        WHERE t.status = 'active'
        GROUP BY t.symbol
        ORDER BY t.symbol`)
	if err != nil {
		return nil, fmt.Errorf("query tickers: %w", err)
	}
	defer rows.Close()

	var out []domain.Ticker
	for rows.Next() {
		var t domain.Ticker
		if err := rows.Scan(&t.Symbol, &t.CompanyName, &t.ShortName, &t.Exchange,
			&t.Sector, &t.InVN30, &t.ArticleCount); err != nil {
			return nil, fmt.Errorf("scan ticker: %w", err)
		}
		out = append(out, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate tickers: %w", err)
	}
	return out, nil
}

// GetTicker đọc thông tin một mã.
func (s *Store) GetTicker(ctx context.Context, symbol string) (domain.Ticker, error) {
	var t domain.Ticker
	err := s.pool.QueryRow(ctx, `
        SELECT symbol, company_name, short_name, exchange, sector, in_vn30
        FROM tickers WHERE symbol = $1 AND status = 'active'`, symbol).
		Scan(&t.Symbol, &t.CompanyName, &t.ShortName, &t.Exchange, &t.Sector, &t.InVN30)
	if err != nil {
		return domain.Ticker{}, wrapNoRows(err, fmt.Sprintf("ticker %s", symbol))
	}
	return t, nil
}

// ListTickersForTagger nạp toàn bộ mã + alias để dựng từ điển tagger.
func (s *Store) ListTickersForTagger(ctx context.Context) ([]domain.Ticker, error) {
	rows, err := s.pool.Query(ctx, `
        SELECT t.symbol, t.company_name, t.short_name, t.sector, t.in_vn30,
               COALESCE(array_agg(ta.alias) FILTER (WHERE ta.alias IS NOT NULL), '{}'),
               COALESCE(t.negative_aliases, '{}')
        FROM tickers t
        LEFT JOIN ticker_aliases ta ON ta.symbol = t.symbol
        WHERE t.status = 'active'
        GROUP BY t.symbol
        ORDER BY t.symbol`)
	if err != nil {
		return nil, fmt.Errorf("load ticker dictionary: %w", err)
	}
	defer rows.Close()

	var out []domain.Ticker
	for rows.Next() {
		var t domain.Ticker
		var aliases []string
		if err := rows.Scan(&t.Symbol, &t.CompanyName, &t.ShortName, &t.Sector, &t.InVN30,
			&aliases, &t.NegativeAliases); err != nil {
			return nil, fmt.Errorf("scan ticker dictionary row: %w", err)
		}
		for _, a := range aliases {
			t.Aliases = append(t.Aliases, domain.TickerAlias{Symbol: t.Symbol, Alias: a})
		}
		out = append(out, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate ticker dictionary: %w", err)
	}
	return out, nil
}
