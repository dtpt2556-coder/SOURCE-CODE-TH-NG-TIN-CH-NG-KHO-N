package fetcher

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
)

// DiscoverWindow — R1.2 của BA2: chỉ nhận bài đăng trong 48 giờ gần nhất.
const DiscoverWindow = 48 * time.Hour

// blacklistPathParts — R1.5: các URL không phải bài text tóm tắt được.
var blacklistPathParts = []string{
	"/video/", "/podcast/", "/infographic/", "/emagazine/", "/rao-vat/",
	"/quang-cao/", "/tag/", "/chuyen-muc/", "/photo/", "/anh/",
}

// FeedItem là một ứng viên bài viết phát hiện từ RSS.
type FeedItem struct {
	URL           string
	TitleHint     string
	PublishedHint time.Time
	HasDate       bool
}

// ParseFeed đọc nội dung RSS/Atom và trả về các ứng viên hợp lệ, cùng cờ cho
// biết danh sách có bị cắt vì vượt `limit` hay không.
// Bài cũ hơn `window` hoặc nằm trong blacklist path bị loại ngay.
func ParseFeed(ctx context.Context, body string, now time.Time, window time.Duration, limit int) ([]FeedItem, bool, error) {
	if strings.TrimSpace(body) == "" {
		return nil, false, fmt.Errorf("feed body is empty")
	}
	if err := ctx.Err(); err != nil {
		return nil, false, fmt.Errorf("context cancelled before parsing feed: %w", err)
	}
	parser := gofeed.NewParser()
	feed, err := parser.ParseString(body)
	if err != nil {
		return nil, false, fmt.Errorf("parse feed: %w", err)
	}
	if window <= 0 {
		window = DiscoverWindow
	}
	cutoff := now.Add(-window)

	out := make([]FeedItem, 0, len(feed.Items))
	for _, it := range feed.Items {
		if it == nil || strings.TrimSpace(it.Link) == "" {
			continue
		}
		if isBlacklisted(it.Link) {
			continue
		}
		item := FeedItem{URL: strings.TrimSpace(it.Link), TitleHint: strings.TrimSpace(it.Title)}
		if it.PublishedParsed != nil {
			item.PublishedHint = it.PublishedParsed.UTC()
			item.HasDate = true
		} else if it.UpdatedParsed != nil {
			item.PublishedHint = it.UpdatedParsed.UTC()
			item.HasDate = true
		}
		if item.HasDate && item.PublishedHint.Before(cutoff) {
			continue
		}
		out = append(out, item)
	}

	// Sắp theo thời gian giảm dần TRƯỚC khi cắt. Feed backfill lịch sử hay trả
	// bài cũ lên đầu; cắt theo thứ tự gặp sẽ lấy nhầm 40 bài cũ nhất (D-05).
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].HasDate != out[j].HasDate {
			return out[i].HasDate
		}
		return out[i].PublishedHint.After(out[j].PublishedHint)
	})

	truncated := false
	if limit > 0 && len(out) > limit {
		out = out[:limit]
		truncated = true
	}
	return out, truncated, nil
}

func isBlacklisted(rawURL string) bool {
	lower := strings.ToLower(rawURL)
	for _, p := range blacklistPathParts {
		if strings.Contains(lower, p) {
			return true
		}
	}
	return false
}
