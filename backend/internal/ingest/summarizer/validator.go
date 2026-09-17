package summarizer

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/tenpoint/tenpoint-api/internal/textutil"
)

// boldRe bắt các cụm **...** trong markdown tóm tắt.
var boldRe = regexp.MustCompile(`\*\*([^*]+)\*\*`)

// MaxSummaryWords — tóm tắt dài hơn ngưỡng này bị loại, sinh lại (S-09).
const MaxSummaryWords = 180

// VerbatimRunLimit — sao chép nguyên văn từ này trở lên là vi phạm ranh giới
// bản quyền của PRD mục 4 (S-10).
const VerbatimRunLimit = 15

// ValidationResult là kết quả đối chiếu số liệu.
type ValidationResult struct {
	Grounded   bool
	Violations []string
}

// Error mô tả vi phạm để ghi log cho QA rà soát.
func (r ValidationResult) Error() string {
	if r.Grounded {
		return ""
	}
	return fmt.Sprintf("summary contains %d ungrounded numbers: %s",
		len(r.Violations), strings.Join(r.Violations, ", "))
}

// Validator đối chiếu mọi số trong tóm tắt với body bài gốc.
// BA3 xếp "ảo giác số liệu" vào nhóm rủi ro mức tồn vong — đây là hàng rào an
// toàn quan trọng nhất của sản phẩm.
type Validator struct{}

// NewValidator tạo validator (stateless, an toàn khi dùng đồng thời).
func NewValidator() *Validator { return &Validator{} }

// Validate trả Grounded=false nếu có số nào trong summary không tồn tại trong
// body. Khi đó pipeline phải đặt status='pending_review' (S-01, S-02, S-03).
func (v *Validator) Validate(summaryMD, body string) ValidationResult {
	grounded := textutil.NumberSet(body)

	seen := map[string]bool{}
	var violations []string
	for _, raw := range textutil.RawNumbers(summaryMD) {
		norm, ok := textutil.NormalizeVNNumber(raw)
		if !ok || grounded[norm] || seen[raw] {
			continue
		}
		seen[raw] = true
		violations = append(violations, raw)
	}
	return ValidationResult{Grounded: len(violations) == 0, Violations: violations}
}

// BoldSpans trích các cụm được in đậm trong markdown.
func BoldSpans(summaryMD string) []string {
	matches := boldRe.FindAllStringSubmatch(summaryMD, -1)
	out := make([]string, 0, len(matches))
	for _, m := range matches {
		out = append(out, m[1])
	}
	return out
}

var (
	mdHeadingRe = regexp.MustCompile(`(?m)^\s{0,3}#{1,6}\s*`)
	mdBulletRe  = regexp.MustCompile(`(?m)^\s{0,3}[-*+]\s+`)
	mdLinkRe    = regexp.MustCompile(`\[([^\]]*)\]\([^)]*\)`)
	mdCodeRe    = regexp.MustCompile("`+")
)

// SanitizeMarkdown giữ lại đúng cú pháp in đậm, loại mọi markdown khác mà LLM
// có thể trả về (S-07). FE chỉ sanitize cho phép <strong>/<em>, nên heading và
// list lọt xuống sẽ hiển thị thành ký tự thô trên bảng digest.
func SanitizeMarkdown(s string) string {
	s = mdHeadingRe.ReplaceAllString(s, "")
	s = mdBulletRe.ReplaceAllString(s, "")
	s = mdLinkRe.ReplaceAllString(s, "$1")
	s = mdCodeRe.ReplaceAllString(s, "")
	s = strings.ReplaceAll(s, "\n", " ")
	return strings.Join(strings.Fields(s), " ")
}

// LongestVerbatimRun trả về độ dài chuỗi từ liên tiếp dài nhất mà tóm tắt sao
// chép nguyên văn từ bài gốc (S-10).
func LongestVerbatimRun(summaryMD, body string) int {
	summaryWords := textutil.Tokenize(textutil.Normalize(strings.ReplaceAll(summaryMD, "*", "")))
	if len(summaryWords) == 0 {
		return 0
	}
	bodyText := " " + strings.Join(textutil.Tokenize(textutil.Normalize(body)), " ") + " "

	longest := 0
	for start := range summaryWords {
		for end := start + longest + 1; end <= len(summaryWords); end++ {
			phrase := " " + strings.Join(summaryWords[start:end], " ") + " "
			if !strings.Contains(bodyText, phrase) {
				break
			}
			longest = end - start
		}
	}
	return longest
}

// QualityIssue mô tả lý do một bản tóm tắt bị loại.
type QualityIssue struct {
	Code   string
	Detail string
}

// CheckQuality chạy các rule hình thức chạy được mà không cần bài gốc đầy đủ.
// Trả nil khi đạt.
func CheckQuality(summaryMD string) *QualityIssue {
	if strings.TrimSpace(summaryMD) == "" {
		return &QualityIssue{Code: "empty", Detail: "summary is empty"}
	}
	if words := len(strings.Fields(summaryMD)); words > MaxSummaryWords {
		return &QualityIssue{Code: "too_long", Detail: fmt.Sprintf("%d words > %d", words, MaxSummaryWords)}
	}
	if phrase, banned := FindBannedPhrase(summaryMD); banned {
		return &QualityIssue{Code: "banned_phrase", Detail: phrase}
	}
	return nil
}
