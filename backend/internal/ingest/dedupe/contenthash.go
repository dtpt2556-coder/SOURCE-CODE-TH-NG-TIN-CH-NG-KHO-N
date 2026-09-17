package dedupe

import (
	"crypto/sha256"
	"encoding/hex"
	"regexp"
	"strings"

	"github.com/tenpoint/tenpoint-api/internal/textutil"
)

// RewriteThreshold — dưới mức thay đổi này coi là sửa chính tả, chỉ cập nhật
// content_hash chứ không tóm tắt lại (U-02). Từ mức này trở lên mới gọi
// summarizer, vì mỗi lần gọi LLM là một lần tốn tiền.
const RewriteThreshold = 0.05

// volatilePatterns là những đoạn thay đổi ở MỌI lần crawl dù bài không sửa gì.
// Không loại chúng thì bài nào cũng bị coi là "đã cập nhật" và pipeline sẽ tóm
// tắt lại toàn bộ kho tin sau mỗi run (U-06).
var volatilePatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)cập nhật\s*(lúc)?\s*[:\-]?\s*\d{1,2}[:h]\d{2}`),
	regexp.MustCompile(`(?i)\d{1,2}[:h]\d{2}(:\d{2})?\s*(gmt[+-]?\d*)?`),
	regexp.MustCompile(`\d{1,2}[/-]\d{1,2}[/-]\d{2,4}`),
	regexp.MustCompile(`\d{4}-\d{2}-\d{2}`),
	regexp.MustCompile(`(?i)[\d.,]+\s*(lượt xem|lượt đọc|lượt thích|lượt chia sẻ|người theo dõi)`),
	regexp.MustCompile(`(?i)[\d.,]+\s*(bình luận|phản hồi|comment)`),
	regexp.MustCompile(`(?i)(thứ\s+(hai|ba|tư|năm|sáu|bảy)|chủ\s+nhật)\s*,?`),
}

// ContentHash sinh SHA-256 của thân bài đã chuẩn hoá.
func ContentHash(body string) string {
	sum := sha256.Sum256([]byte(normalizeBody(body)))
	return hex.EncodeToString(sum[:])
}

// normalizeBody bỏ phần biến động rồi gom khoảng trắng, hạ chữ thường.
func normalizeBody(body string) string {
	out := body
	for _, re := range volatilePatterns {
		out = re.ReplaceAllString(out, " ")
	}
	return strings.ToLower(textutil.CollapseSpaces(out))
}

// ChangeRatio đo mức thay đổi giữa hai phiên bản thân bài, 0 là giống hệt và 1
// là không còn từ nào chung. Dùng multiset từ thay vì diff ký tự vì báo hay đảo
// vị trí đoạn mà không đổi nội dung.
func ChangeRatio(before, after string) float64 {
	a := wordCounts(before)
	b := wordCounts(after)
	if len(a) == 0 && len(b) == 0 {
		return 0
	}

	total, shared := 0, 0
	for word, countA := range a {
		total += countA
		if countB, ok := b[word]; ok {
			shared += min(countA, countB)
		}
	}
	for _, countB := range b {
		total += countB
	}
	if total == 0 {
		return 0
	}
	return 1 - float64(2*shared)/float64(total)
}

// NeedsResummarize quyết định có gọi lại summarizer hay không (U-02, U-03).
func NeedsResummarize(oldTitle, oldBody, newTitle, newBody string) bool {
	if strings.TrimSpace(oldTitle) != strings.TrimSpace(newTitle) {
		return true
	}
	if textutil.NumbersDiffer(oldBody, newBody) {
		return true
	}
	return ChangeRatio(oldBody, newBody) >= RewriteThreshold
}

func wordCounts(text string) map[string]int {
	out := map[string]int{}
	for _, w := range textutil.Tokenize(textutil.Normalize(text)) {
		out[w]++
	}
	return out
}
