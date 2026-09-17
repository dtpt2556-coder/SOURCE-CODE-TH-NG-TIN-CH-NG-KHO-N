package domain

// Source là một nguồn tin đã được duyệt (xem §1 của BA2).
// Tier 0 = nguồn sơ cấp (HOSE/HNX/UBCKNN), 1 = báo tài chính uy tín, 2 = báo tổng hợp.
type Source struct {
	ID          int
	Code        string
	Name        string
	Domain      string
	RSSURL      string
	ListURL     string
	Tier        int
	Enabled     bool
	RateLimitMS int
}
