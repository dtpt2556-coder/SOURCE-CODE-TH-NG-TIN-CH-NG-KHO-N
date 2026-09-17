package store

import (
	"context"
	"fmt"
	"time"

	"github.com/tenpoint/tenpoint-api/internal/domain"
)

// StoredArticle là trạng thái hiện có của một bài, đủ để quyết định xử lý lần
// crawl tiếp theo (mục 6 của ingest-edge-cases).
type StoredArticle struct {
	ID           int64
	CanonicalURL string
	Title        string
	Body         string
	SummaryMD    string
	ContentHash  string
	Status       string
	Revision     int
	ShortCode    string
}

const storedArticleColumns = `
    a.id, a.canonical_url, a.title, a.excerpt, a.summary_md,
    COALESCE(a.content_hash, ''), a.status, a.revision, COALESCE(sl.code, '')`

// FindArticleByURLHash tìm bài theo url_hash đã chuẩn hoá.
func (s *Store) FindArticleByURLHash(ctx context.Context, urlHash string) (StoredArticle, error) {
	return s.findArticle(ctx, `
        SELECT`+storedArticleColumns+`
        FROM articles a
        LEFT JOIN short_links sl ON sl.article_id = a.id
        WHERE a.url_hash = $1`, urlHash)
}

// FindArticleByContentHash tìm bài cùng nguồn có nội dung y hệt.
//
// Đây là cách bắt được bài đổi URL: nội dung không đổi nhưng `url_hash` mới nên
// tra theo URL sẽ trượt và pipeline tạo bản ghi trùng (U-05, C-05).
func (s *Store) FindArticleByContentHash(ctx context.Context, sourceID int, contentHash string) (StoredArticle, error) {
	if contentHash == "" {
		return StoredArticle{}, ErrNotFound
	}
	return s.findArticle(ctx, `
        SELECT`+storedArticleColumns+`
        FROM articles a
        LEFT JOIN short_links sl ON sl.article_id = a.id
        WHERE a.source_id = $1 AND a.content_hash = $2
        ORDER BY a.id
        LIMIT 1`, sourceID, contentHash)
}

func (s *Store) findArticle(ctx context.Context, query string, args ...any) (StoredArticle, error) {
	var a StoredArticle
	err := s.pool.QueryRow(ctx, query, args...).Scan(
		&a.ID, &a.CanonicalURL, &a.Title, &a.Body, &a.SummaryMD,
		&a.ContentHash, &a.Status, &a.Revision, &a.ShortCode)
	if err != nil {
		return StoredArticle{}, wrapNoRows(err, "find article")
	}
	return a, nil
}

// TouchArticle ghi nhận bài vẫn còn sống ở nguồn và nội dung chưa đổi đáng kể.
// Không đụng tới summary_md nên không tốn một lượt gọi LLM nào (U-01, U-02).
func (s *Store) TouchArticle(ctx context.Context, id int64, contentHash string, seenAt time.Time) error {
	_, err := s.pool.Exec(ctx, `
        UPDATE articles SET last_seen_at = $2, content_hash = $3
        WHERE id = $1`, id, seenAt, contentHash)
	if err != nil {
		return fmt.Errorf("touch article %d: %w", id, err)
	}
	return nil
}

// Revision là bản cập nhật đã qua kiểm chứng của một bài.
type Revision struct {
	ArticleID   int64
	Title       string
	SummaryMD   string
	Body        string
	ContentHash string
	NewsType    domain.NewsType
	Provider    string
	Tickers     []domain.TickerRef
	UpdatedAt   time.Time
}

// ApplyRevision thay nội dung đang publish bằng bản mới đã qua validator, tăng
// revision và ghi updated_at để FE hiện nhãn "Đã cập nhật lúc HH:mm" (U-03).
func (s *Store) ApplyRevision(ctx context.Context, rev Revision) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin revision tx: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	_, err = tx.Exec(ctx, `
        UPDATE articles
        SET title = $2, summary_md = $3, excerpt = $4, content_hash = $5,
            news_type = $6, summary_provider = $7, revision = revision + 1,
            updated_at = $8, last_seen_at = $8, status = 'published',
            pending_summary_md = NULL, reject_reason = NULL
        WHERE id = $1`,
		rev.ArticleID, rev.Title, rev.SummaryMD, rev.Body, rev.ContentHash,
		string(rev.NewsType), rev.Provider, rev.UpdatedAt)
	if err != nil {
		return fmt.Errorf("update article %d: %w", rev.ArticleID, err)
	}

	if _, err := tx.Exec(ctx, `DELETE FROM article_tickers WHERE article_id = $1`, rev.ArticleID); err != nil {
		return fmt.Errorf("clear tickers of article %d: %w", rev.ArticleID, err)
	}
	for _, t := range rev.Tickers {
		if _, err := tx.Exec(ctx, `
            INSERT INTO article_tickers (article_id, symbol, relevance, score)
            VALUES ($1,$2,$3,$4) ON CONFLICT DO NOTHING`,
			rev.ArticleID, t.Symbol, t.Relevance, float32(t.Score)); err != nil {
			return fmt.Errorf("attach ticker %s to article %d: %w", t.Symbol, rev.ArticleID, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit revision of article %d: %w", rev.ArticleID, err)
	}
	return nil
}

// HoldPendingRevision giữ nguyên bản tóm tắt đang publish và cất bản mới sang
// pending_summary_md.
//
// Bản mới trượt validator nghĩa là nó chứa số không truy vết được về bài gốc.
// Thay nội dung đã kiểm chứng bằng nội dung chưa kiểm chứng là đúng thứ mà hàng
// rào an toàn này sinh ra để ngăn (U-07).
func (s *Store) HoldPendingRevision(ctx context.Context, id int64, summaryMD, reason string, seenAt time.Time) error {
	_, err := s.pool.Exec(ctx, `
        UPDATE articles
        SET pending_summary_md = $2, reject_reason = $3, last_seen_at = $4
        WHERE id = $1`, id, summaryMD, reason, seenAt)
	if err != nil {
		return fmt.Errorf("hold pending revision of article %d: %w", id, err)
	}
	return nil
}

// RelocateArticle cập nhật URL mới cho bài đã đổi địa chỉ, giữ nguyên article.id
// và short_links.code để link đã chia sẻ không chết (U-05, L-05).
func (s *Store) RelocateArticle(ctx context.Context, id int64, canonicalURL, urlHash string) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin relocate tx: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if _, err := tx.Exec(ctx, `
        UPDATE articles SET canonical_url = $2, url_hash = $3 WHERE id = $1`,
		id, canonicalURL, urlHash); err != nil {
		return fmt.Errorf("relocate article %d: %w", id, err)
	}
	if _, err := tx.Exec(ctx, `
        UPDATE short_links SET target_url = $2, target_alive = true, last_checked_at = now()
        WHERE article_id = $1`, id, canonicalURL); err != nil {
		return fmt.Errorf("retarget short link of article %d: %w", id, err)
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit relocate of article %d: %w", id, err)
	}
	return nil
}

// MarkSourceGone xử lý bài đã bị gỡ khỏi nguồn: giữ lại tóm tắt vì đó là nội
// dung của TenPoint, nhưng tắt link để không dẫn người dùng tới trang 404 (U-04).
func (s *Store) MarkSourceGone(ctx context.Context, id int64) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin source-gone tx: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if _, err := tx.Exec(ctx,
		`UPDATE articles SET status = 'source_gone' WHERE id = $1`, id); err != nil {
		return fmt.Errorf("mark article %d gone: %w", id, err)
	}
	if _, err := tx.Exec(ctx, `
        UPDATE short_links SET target_alive = false, last_checked_at = now()
        WHERE article_id = $1`, id); err != nil {
		return fmt.Errorf("disable short link of article %d: %w", id, err)
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit source-gone of article %d: %w", id, err)
	}
	return nil
}

// HasRealArticles cho biết đã có bài thật (không phải seed demo) hay chưa.
// Quyết định cả việc bootstrap lúc khởi động lẫn việc có ẩn demo hay không.
func (s *Store) HasRealArticles(ctx context.Context) (bool, error) {
	var exists bool
	err := s.pool.QueryRow(ctx, `
        SELECT EXISTS (
            SELECT 1 FROM articles
            WHERE is_demo = false AND status IN ('published', 'source_gone'))`).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check for real articles: %w", err)
	}
	return exists, nil
}

// TickerSymbolsForShortCode trả về mã CK của bài đứng sau một short link, dùng
// để trang 410 dẫn người đọc sang trang mã liên quan thay vì bỏ họ ở ngõ cụt.
func (s *Store) TickerSymbolsForShortCode(ctx context.Context, code string) ([]string, error) {
	rows, err := s.pool.Query(ctx, `
        SELECT t.symbol
        FROM short_links sl
        JOIN article_tickers t ON t.article_id = sl.article_id
        WHERE sl.code = $1
        ORDER BY t.score DESC, t.symbol
        LIMIT 5`, code)
	if err != nil {
		return nil, fmt.Errorf("query tickers for short code %s: %w", code, err)
	}
	defer rows.Close()

	var out []string
	for rows.Next() {
		var symbol string
		if err := rows.Scan(&symbol); err != nil {
			return nil, fmt.Errorf("scan ticker symbol: %w", err)
		}
		out = append(out, symbol)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate ticker symbols: %w", err)
	}
	return out, nil
}

// SummaryAlreadyUsed cho biết một nguồn đã dùng đúng bản tóm tắt này cho bài
// khác chưa.
//
// Hai bài khác nhau cùng nguồn mà ra đúng một tóm tắt luôn là dấu hiệu bóc tách
// hỏng, không phải trường hợp hợp lệ: đó chính là cách 36 bài VnEconomy cùng
// hiển thị nội dung về hội thảo tầng ô-dôn mà không cổng nào kêu.
func (s *Store) SummaryAlreadyUsed(ctx context.Context, sourceID int, summaryMD string, exceptID int64) (bool, error) {
	var exists bool
	err := s.pool.QueryRow(ctx, `
        SELECT EXISTS (
            SELECT 1 FROM articles
            WHERE source_id = $1 AND summary_md = $2 AND id <> $3 AND summary_md <> '')`,
		sourceID, summaryMD, exceptID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check duplicate summary for source %d: %w", sourceID, err)
	}
	return exists, nil
}

// DuplicateSummaryGroups liệt kê các nguồn đang có tóm tắt dùng chung, phục vụ
// cảnh báo vận hành và kiểm thử nghiệm thu.
func (s *Store) DuplicateSummaryGroups(ctx context.Context) (map[int]int, error) {
	rows, err := s.pool.Query(ctx, `
        SELECT source_id, count(*)
        FROM (
            SELECT source_id, summary_md
            FROM articles
            WHERE status = 'published' AND is_demo = false AND summary_md <> ''
            GROUP BY source_id, summary_md
            HAVING count(*) > 1
        ) AS dup
        GROUP BY source_id`)
	if err != nil {
		return nil, fmt.Errorf("query duplicate summary groups: %w", err)
	}
	defer rows.Close()

	out := map[int]int{}
	for rows.Next() {
		var sourceID, groups int
		if err := rows.Scan(&sourceID, &groups); err != nil {
			return nil, fmt.Errorf("scan duplicate summary group: %w", err)
		}
		out[sourceID] = groups
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate duplicate summary groups: %w", err)
	}
	return out, nil
}

// SourceHealth ghi lại kết quả feed rỗng liên tiếp để cảnh báo nguồn đổi cấu
// trúc (D-04).
func (s *Store) SourceHealth(ctx context.Context, sourceID int, found int) (int, error) {
	var streak int
	err := s.pool.QueryRow(ctx, `
        UPDATE sources
        SET consecutive_empty_runs = CASE WHEN $2 = 0 THEN consecutive_empty_runs + 1 ELSE 0 END
        WHERE id = $1
        RETURNING consecutive_empty_runs`, sourceID, found).Scan(&streak)
	if err != nil {
		return 0, fmt.Errorf("update source health %d: %w", sourceID, err)
	}
	return streak, nil
}
