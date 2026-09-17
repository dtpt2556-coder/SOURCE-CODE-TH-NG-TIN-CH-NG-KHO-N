package store

import (
	"context"
	"fmt"

	"github.com/tenpoint/tenpoint-api/internal/domain"
)

// ListNotes trả về danh sách research note, mới nhất trước.
func (s *Store) ListNotes(ctx context.Context, limit int) ([]domain.ResearchNote, error) {
	if limit <= 0 {
		limit = 50
	}
	rows, err := s.pool.Query(ctx, `
        SELECT id, slug, COALESCE(symbol, ''), title, section_heading, disclaimer, published_at
        FROM research_notes
        ORDER BY published_at DESC, id DESC
        LIMIT $1`, limit)
	if err != nil {
		return nil, fmt.Errorf("query research notes: %w", err)
	}
	defer rows.Close()

	var out []domain.ResearchNote
	for rows.Next() {
		var n domain.ResearchNote
		if err := rows.Scan(&n.ID, &n.Slug, &n.Symbol, &n.Title, &n.SectionHeading,
			&n.Disclaimer, &n.PublishedAt); err != nil {
			return nil, fmt.Errorf("scan research note: %w", err)
		}
		out = append(out, n)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate research notes: %w", err)
	}
	return out, nil
}

// ListNotesBySymbol trả về research note của một mã.
func (s *Store) ListNotesBySymbol(ctx context.Context, symbol string) ([]domain.ResearchNote, error) {
	rows, err := s.pool.Query(ctx, `
        SELECT id, slug, COALESCE(symbol, ''), title, section_heading, disclaimer, published_at
        FROM research_notes
        WHERE symbol = $1
        ORDER BY published_at DESC`, symbol)
	if err != nil {
		return nil, fmt.Errorf("query research notes for ticker %s: %w", symbol, err)
	}
	defer rows.Close()

	var out []domain.ResearchNote
	for rows.Next() {
		var n domain.ResearchNote
		if err := rows.Scan(&n.ID, &n.Slug, &n.Symbol, &n.Title, &n.SectionHeading,
			&n.Disclaimer, &n.PublishedAt); err != nil {
			return nil, fmt.Errorf("scan research note: %w", err)
		}
		out = append(out, n)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate research notes: %w", err)
	}
	return out, nil
}

// GetNote đọc 1 research note kèm đầy đủ luận điểm.
func (s *Store) GetNote(ctx context.Context, slug string) (domain.ResearchNote, error) {
	var n domain.ResearchNote
	err := s.pool.QueryRow(ctx, `
        SELECT id, slug, COALESCE(symbol, ''), title, section_heading, disclaimer, published_at
        FROM research_notes WHERE slug = $1`, slug).
		Scan(&n.ID, &n.Slug, &n.Symbol, &n.Title, &n.SectionHeading, &n.Disclaimer, &n.PublishedAt)
	if err != nil {
		return domain.ResearchNote{}, wrapNoRows(err, fmt.Sprintf("research note %s", slug))
	}

	rows, err := s.pool.Query(ctx, `
        SELECT id, note_id, ordinal, lead, body
        FROM research_note_points
        WHERE note_id = $1
        ORDER BY ordinal`, n.ID)
	if err != nil {
		return domain.ResearchNote{}, fmt.Errorf("query note points for %s: %w", slug, err)
	}
	defer rows.Close()

	for rows.Next() {
		var p domain.ResearchNotePoint
		if err := rows.Scan(&p.ID, &p.NoteID, &p.Ordinal, &p.Lead, &p.Body); err != nil {
			return domain.ResearchNote{}, fmt.Errorf("scan note point: %w", err)
		}
		n.Points = append(n.Points, p)
	}
	if err := rows.Err(); err != nil {
		return domain.ResearchNote{}, fmt.Errorf("iterate note points: %w", err)
	}
	return n, nil
}
