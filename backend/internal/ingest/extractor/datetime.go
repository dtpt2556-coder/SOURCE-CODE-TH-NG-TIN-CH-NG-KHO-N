package extractor

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/tenpoint/tenpoint-api/internal/textutil"
)

// FutureTolerance — CMS nguồn thỉnh thoảng đẩy ngày về tương lai. Lệch quá
// ngưỡng này thì kẹp về thời điểm fetch thay vì tin vào nguồn (T-07).
const FutureTolerance = 2 * time.Hour

// MaxArticleAge — bài cũ hơn ngưỡng này trong feed "mới nhất" gần như chắc chắn
// là bài cũ được CMS đẩy lại (T-08).
const MaxArticleAge = 90 * 24 * time.Hour

var (
	relativeRe = regexp.MustCompile(`(\d+)\s*(giay|phut|gio|ngay|tuan|thang)\s+truoc`)
	dmyRe      = regexp.MustCompile(`\b(\d{1,2})[/-](\d{1,2})[/-](\d{4})\b`)
	ymdRe      = regexp.MustCompile(`\b(\d{4})-(\d{1,2})-(\d{1,2})\b`)
	clockRe    = regexp.MustCompile(`\b(\d{1,2}):(\d{2})(?::(\d{2}))?\b`)
	tzSuffixRe = regexp.MustCompile(`\((?:gmt|utc)[+-]?\d{0,2}(?::\d{2})?\)`)
)

// absoluteLayouts là các định dạng có offset rõ ràng, thử trước mọi heuristic.
var absoluteLayouts = []string{
	time.RFC3339,
	"2006-01-02T15:04:05Z0700",
	"2006-01-02T15:04:05-07:00",
	time.RFC1123Z,
	time.RFC1123,
	"Mon, 02 Jan 2006 15:04:05 -0700",
	"Mon, 02 Jan 06 15:04:05 -0700",
}

// ParsePublishedAt đọc thời gian đăng từ chuỗi bất kỳ trên báo Việt Nam và trả
// về thời điểm UTC.
//
// Thứ tự: định dạng có offset tuyệt đối -> thời gian tương đối ("2 giờ trước")
// -> ngày tuyệt đối kiểu Việt Nam. Chuỗi không ghi múi giờ được hiểu là ICT,
// vì đó là mặc định của mọi CMS báo trong nước (T-02).
func ParsePublishedAt(raw string, now time.Time, loc *time.Location) (time.Time, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return time.Time{}, false
	}
	if loc == nil {
		loc = time.UTC
	}

	for _, layout := range absoluteLayouts {
		if t, err := time.Parse(layout, raw); err == nil {
			return t.UTC(), true
		}
	}

	flat := textutil.Normalize(raw)
	flat = tzSuffixRe.ReplaceAllString(flat, " ")
	flat = strings.Join(strings.Fields(flat), " ")

	if t, ok := parseRelative(flat, now, loc); ok {
		return t, true
	}
	return parseAbsoluteVN(flat, loc)
}

func parseRelative(flat string, now time.Time, loc *time.Location) (time.Time, bool) {
	if m := relativeRe.FindStringSubmatch(flat); m != nil {
		n, err := strconv.Atoi(m[1])
		if err != nil {
			return time.Time{}, false
		}
		var delta time.Duration
		switch m[2] {
		case "giay":
			delta = time.Duration(n) * time.Second
		case "phut":
			delta = time.Duration(n) * time.Minute
		case "gio":
			delta = time.Duration(n) * time.Hour
		case "ngay":
			delta = time.Duration(n) * 24 * time.Hour
		case "tuan":
			delta = time.Duration(n) * 7 * 24 * time.Hour
		case "thang":
			delta = time.Duration(n) * 30 * 24 * time.Hour
		default:
			return time.Time{}, false
		}
		return now.Add(-delta).UTC(), true
	}

	dayOffset := 0
	switch {
	case strings.Contains(flat, "hom nay"), strings.Contains(flat, "vua xong"):
		dayOffset = 0
	case strings.Contains(flat, "hom qua"):
		dayOffset = -1
	case strings.Contains(flat, "hom kia"):
		dayOffset = -2
	default:
		return time.Time{}, false
	}

	base := now.In(loc).AddDate(0, 0, dayOffset)
	hour, minute, sec := 0, 0, 0
	if c := clockRe.FindStringSubmatch(flat); c != nil {
		hour, _ = strconv.Atoi(c[1])
		minute, _ = strconv.Atoi(c[2])
		if c[3] != "" {
			sec, _ = strconv.Atoi(c[3])
		}
	} else if dayOffset == 0 {
		t := now.In(loc)
		hour, minute, sec = t.Hour(), t.Minute(), t.Second()
	}
	return time.Date(base.Year(), base.Month(), base.Day(), hour, minute, sec, 0, loc).UTC(), true
}

func parseAbsoluteVN(flat string, loc *time.Location) (time.Time, bool) {
	var year, month, day int
	if m := dmyRe.FindStringSubmatch(flat); m != nil {
		day, _ = strconv.Atoi(m[1])
		month, _ = strconv.Atoi(m[2])
		year, _ = strconv.Atoi(m[3])
	} else if m := ymdRe.FindStringSubmatch(flat); m != nil {
		year, _ = strconv.Atoi(m[1])
		month, _ = strconv.Atoi(m[2])
		day, _ = strconv.Atoi(m[3])
	} else {
		return time.Time{}, false
	}
	if month < 1 || month > 12 || day < 1 || day > 31 {
		return time.Time{}, false
	}

	hour, minute, sec := 0, 0, 0
	if c := clockRe.FindStringSubmatch(flat); c != nil {
		hour, _ = strconv.Atoi(c[1])
		minute, _ = strconv.Atoi(c[2])
		if c[3] != "" {
			sec, _ = strconv.Atoi(c[3])
		}
	}
	if hour > 23 || minute > 59 || sec > 59 {
		return time.Time{}, false
	}

	t := time.Date(year, time.Month(month), day, hour, minute, sec, 0, loc)
	// time.Date cuộn ngày không hợp lệ (31/02 -> 03/03); từ chối thay vì lặng lẽ
	// nhận một ngày khác với bài gốc.
	if t.Day() != day || int(t.Month()) != month {
		return time.Time{}, false
	}
	return t.UTC(), true
}

// SettlePublishedAt chốt thời gian đăng cuối cùng và cho biết nó có phải ước
// lượng hay không.
//
// Không parse được, hoặc ở tương lai quá FutureTolerance, thì lấy fetchedAt và
// bật cờ estimated để FE hiển thị dấu "~" trước ngày (T-06, T-07).
func SettlePublishedAt(parsed time.Time, hasValue bool, fetchedAt time.Time) (time.Time, bool) {
	if !hasValue || parsed.IsZero() {
		return fetchedAt.UTC(), true
	}
	if parsed.After(fetchedAt.Add(FutureTolerance)) {
		return fetchedAt.UTC(), true
	}
	return parsed.UTC(), false
}

// TooOld nhận diện bài quá cũ lọt vào feed "mới nhất" (T-08).
func TooOld(publishedAt, now time.Time) bool {
	return publishedAt.Before(now.Add(-MaxArticleAge))
}
