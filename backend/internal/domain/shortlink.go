package domain

import "time"

// ShortLink là bản ghi redirect nội bộ /r/{code} (mục 7 architecture.md).
// TargetURL bắt buộc thuộc allowlist domain của `sources` — chống open redirect.
type ShortLink struct {
	ID            int64
	Code          string
	ArticleID     int64
	TargetURL     string
	ClickCount    int64
	LastCheckedAt *time.Time
	TargetAlive   bool
}

// Path trả về đường dẫn tương đối dùng trong API response: "/r/a7Kx2p".
func (s ShortLink) Path() string {
	if s.Code == "" {
		return ""
	}
	return "/r/" + s.Code
}
