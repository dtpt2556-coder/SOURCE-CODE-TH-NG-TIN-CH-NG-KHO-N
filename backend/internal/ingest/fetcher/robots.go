package fetcher

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"
)

// RobotsTTL — cache robots.txt 24 giờ theo mục 5 architecture.md.
const RobotsTTL = 24 * time.Hour

type robotsRules struct {
	disallow  []string
	allow     []string
	crawlWait time.Duration
	fetchedAt time.Time
}

// RobotsCache tải và cache robots.txt theo host.
type RobotsCache struct {
	mu    sync.Mutex
	rules map[string]*robotsRules
	get   func(ctx context.Context, rawURL string) (string, error)
	ttl   time.Duration
}

// NewRobotsCache nhận hàm tải trang (inject để test offline được).
func NewRobotsCache(get func(ctx context.Context, rawURL string) (string, error)) *RobotsCache {
	return &RobotsCache{rules: make(map[string]*robotsRules), get: get, ttl: RobotsTTL}
}

// Allowed cho biết UA của ta có được phép crawl rawURL hay không.
// Không tải được robots.txt => coi như cho phép (fail-open như mọi crawler lịch
// sự), nhưng vẫn giữ rate limit.
func (c *RobotsCache) Allowed(ctx context.Context, rawURL, userAgent string) (bool, time.Duration, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return false, 0, fmt.Errorf("parse URL %q: %w", rawURL, err)
	}
	r, err := c.rulesFor(ctx, u)
	if err != nil {
		return true, 0, nil
	}
	path := u.EscapedPath()
	if path == "" {
		path = "/"
	}
	// Allow thắng Disallow khi cùng độ dài (theo chuẩn của Google).
	longestDisallow := matchLen(r.disallow, path)
	longestAllow := matchLen(r.allow, path)
	if longestDisallow > longestAllow {
		return false, r.crawlWait, nil
	}
	return true, r.crawlWait, nil
}

func (c *RobotsCache) rulesFor(ctx context.Context, u *url.URL) (*robotsRules, error) {
	host := strings.ToLower(u.Host)

	c.mu.Lock()
	if r, ok := c.rules[host]; ok && time.Since(r.fetchedAt) < c.ttl {
		c.mu.Unlock()
		return r, nil
	}
	c.mu.Unlock()

	robotsURL := u.Scheme + "://" + u.Host + "/robots.txt"
	body, err := c.get(ctx, robotsURL)
	if err != nil {
		return nil, fmt.Errorf("fetch %s: %w", robotsURL, err)
	}
	r := parseRobots(body)
	r.fetchedAt = time.Now()

	c.mu.Lock()
	c.rules[host] = r
	c.mu.Unlock()
	return r, nil
}

// parseRobots đọc các nhóm User-agent: * (nhóm chung). TenPointBot không yêu cầu
// luật riêng nên nhóm "*" là đủ và an toàn nhất.
func parseRobots(body string) *robotsRules {
	r := &robotsRules{}
	applies := false
	for _, line := range strings.Split(body, "\n") {
		line = strings.TrimSpace(line)
		if i := strings.Index(line, "#"); i >= 0 {
			line = strings.TrimSpace(line[:i])
		}
		if line == "" {
			continue
		}
		key, value, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		key = strings.ToLower(strings.TrimSpace(key))
		value = strings.TrimSpace(value)

		switch key {
		case "user-agent":
			applies = value == "*"
		case "disallow":
			if applies && value != "" {
				r.disallow = append(r.disallow, value)
			}
		case "allow":
			if applies && value != "" {
				r.allow = append(r.allow, value)
			}
		case "crawl-delay":
			if applies {
				if d, err := time.ParseDuration(value + "s"); err == nil {
					r.crawlWait = d
				}
			}
		}
	}
	return r
}

// matchLen trả về độ dài prefix khớp dài nhất (0 nếu không khớp).
func matchLen(patterns []string, path string) int {
	best := 0
	for _, p := range patterns {
		if p == "/" {
			if best < 1 {
				best = 1
			}
			continue
		}
		clean := strings.TrimSuffix(p, "*")
		if strings.HasPrefix(path, clean) && len(clean) > best {
			best = len(clean)
		}
	}
	return best
}
