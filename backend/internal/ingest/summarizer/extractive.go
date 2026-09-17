package summarizer

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"unicode"

	"github.com/tenpoint/tenpoint-api/internal/textutil"
)

// vnNumberRe dùng để định vị offset của số trong câu khi bôi đậm.
var vnNumberRe = regexp.MustCompile(`[0-9][0-9.,]*`)

// Extractive là provider mặc định — KHÔNG cần API key (ADR-007).
// Heuristic: chọn các câu giàu số liệu nhất (mật độ số + đơn vị tiền tệ + vị
// trí câu), giữ nguyên văn để mọi con số luôn truy vết được về bài gốc, rồi
// bold cụm số lớn nhất ở câu đầu.
type Extractive struct {
	MinWords int
	MaxWords int
}

// NewExtractive tạo provider với ngưỡng độ dài theo rubric C4 của BA2.
func NewExtractive() *Extractive {
	return &Extractive{MinWords: 80, MaxWords: TargetMaxWords}
}

// Name thoả Provider.
func (e *Extractive) Name() string { return "extractive" }

// Summarize chọn câu giàu số liệu nhất và in đậm cụm số quan trọng nhất.
func (e *Extractive) Summarize(_ context.Context, in Input) (Output, error) {
	body := strings.TrimSpace(in.Body)
	if body == "" {
		return Output{}, fmt.Errorf("body is empty, nothing to summarize")
	}

	sentences := SplitSentences(body)
	if len(sentences) == 0 {
		return Output{}, fmt.Errorf("no sentence could be extracted from body")
	}

	type scored struct {
		idx   int
		text  string
		score float64
		words int
	}
	cands := make([]scored, 0, len(sentences))
	for i, s := range sentences {
		w := len(strings.Fields(s))
		if w < 5 || w > 70 {
			continue
		}
		cands = append(cands, scored{idx: i, text: s, score: sentenceScore(s, i), words: w})
	}
	if len(cands) == 0 {
		// Body toàn câu quá ngắn/quá dài: lấy câu đầu làm tối thiểu.
		cands = append(cands, scored{idx: 0, text: sentences[0], words: len(strings.Fields(sentences[0]))})
	}

	byScore := make([]scored, len(cands))
	copy(byScore, cands)
	sort.SliceStable(byScore, func(i, j int) bool {
		if byScore[i].score != byScore[j].score {
			return byScore[i].score > byScore[j].score
		}
		return byScore[i].idx < byScore[j].idx
	})

	picked := make([]scored, 0, 6)
	total := 0
	for _, c := range byScore {
		if len(picked) > 0 && total+c.words > e.MaxWords {
			continue
		}
		picked = append(picked, c)
		total += c.words
		if total >= e.MinWords || len(picked) >= 6 {
			break
		}
	}
	sort.Slice(picked, func(i, j int) bool { return picked[i].idx < picked[j].idx })

	parts := make([]string, 0, len(picked))
	for _, p := range picked {
		parts = append(parts, p.text)
	}
	summary := strings.Join(parts, " ")

	bolded, span := boldHeadlineNumber(summary)
	if span == "" {
		bolded, span = boldKeyPhrase(summary)
	}
	out := Output{SummaryMD: bolded}
	if span != "" {
		out.BoldSpans = []string{span}
	}
	return out, nil
}

// moneyUnits là các đơn vị làm tăng giá trị thông tin của câu.
var moneyUnits = []string{
	"tỷ đồng", "nghìn tỷ", "tỷ usd", "triệu usd", "usd/tấn", "đồng/kg", "đ/kg",
	"triệu cp", "điểm %", "tỷ", "triệu", "nghìn", "usd", "đồng", "điểm",
	"tấn", "%", "ha", "mw",
}

// sentenceScore chấm điểm một câu theo mật độ số liệu + đơn vị + vị trí.
func sentenceScore(s string, idx int) float64 {
	nums := textutil.RawNumbers(s)
	score := float64(len(nums)) * 2.0

	lower := strings.ToLower(s)
	for _, u := range moneyUnits {
		if strings.Contains(lower, u) {
			score += 1.2
		}
	}
	// Câu đầu bài thường chứa số liệu headline.
	score += 3.0 / float64(idx+1)

	words := len(strings.Fields(s))
	if words < 8 {
		score -= 2.0
	}
	if words > 45 {
		score -= 1.5
	}
	return score
}

// SplitSentences tách câu tiếng Việt an toàn với định dạng số kiểu VN:
// "408.000" KHÔNG bị coi là hết câu vì sau dấu chấm là chữ số, không phải
// khoảng trắng.
func SplitSentences(body string) []string {
	var out []string
	runes := []rune(body)
	start := 0
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		if r != '.' && r != '!' && r != '?' && r != '\n' {
			continue
		}
		if r == '\n' {
			if s := strings.TrimSpace(string(runes[start : i+1])); s != "" {
				out = append(out, s)
			}
			start = i + 1
			continue
		}
		// Phải có khoảng trắng (hoặc hết chuỗi) ngay sau dấu kết câu.
		j := i + 1
		if j < len(runes) && !unicode.IsSpace(runes[j]) {
			continue
		}
		for j < len(runes) && unicode.IsSpace(runes[j]) {
			j++
		}
		// Câu mới phải bắt đầu bằng chữ hoa hoặc chữ số.
		if j < len(runes) && !unicode.IsUpper(runes[j]) && !unicode.IsDigit(runes[j]) {
			continue
		}
		if s := strings.TrimSpace(string(runes[start : i+1])); s != "" {
			out = append(out, s)
		}
		start = j
		// Nhảy con trỏ tới đầu câu mới: nếu không, các ký tự khoảng trắng đã bỏ
		// qua (đặc biệt là "\n") sẽ bị xử lý lại với start > i gây panic.
		i = j - 1
	}
	if s := strings.TrimSpace(string(runes[start:])); s != "" {
		out = append(out, s)
	}
	return out
}

// boldHeadlineNumber in đậm cụm số lớn nhất trong CÂU ĐẦU của tóm tắt
// (rule C5 của BA2: đúng 1 cụm bold, nằm ở câu đầu).
func boldHeadlineNumber(summary string) (string, string) {
	sentences := SplitSentences(summary)
	if len(sentences) == 0 {
		return summary, ""
	}
	first := sentences[0]

	locs := vnNumberRe.FindAllStringIndex(first, -1)
	if len(locs) == 0 {
		return summary, ""
	}

	bestStart, bestEnd := -1, -1
	bestVal := -1.0
	for _, loc := range locs {
		raw := first[loc[0]:loc[1]]
		norm, ok := textutil.NormalizeVNNumber(raw)
		if !ok {
			continue
		}
		val := parseFloatSafe(norm)
		end := loc[1] + unitLength(first[loc[1]:])
		if val > bestVal {
			bestVal, bestStart, bestEnd = val, loc[0], end
		}
	}
	if bestStart < 0 {
		return summary, ""
	}

	span := strings.TrimRight(first[bestStart:bestEnd], " ,.;")
	bestEnd = bestStart + len(span)
	newFirst := first[:bestStart] + "**" + span + "**" + first[bestEnd:]

	// Thay đúng 1 lần, ở câu đầu.
	return strings.Replace(summary, first, newFirst, 1), span
}

// boldKeyPhrase là phương án dự phòng khi câu đầu không có con số nào đáng nhấn.
//
// Rubric NUM-3 yêu cầu 1–3 cụm bold; QA tìm thấy 3 tin có 0 cụm. Bài không có
// số liệu thì nhấn cụm danh từ mang thông tin nhất ở câu đầu thay vì im lặng
// trả về 0 cụm.
func boldKeyPhrase(summary string) (string, string) {
	sentences := SplitSentences(summary)
	if len(sentences) == 0 {
		return summary, ""
	}
	first := sentences[0]

	words := strings.Fields(first)
	if len(words) < 4 {
		return summary, ""
	}

	// Cụm 3 từ đầu tiên bắt đầu bằng chữ hoa (thường là tên riêng, tổ chức, sản
	// phẩm) là thứ đáng nhấn nhất khi không có số.
	for i := 0; i+2 < len(words) && i < 12; i++ {
		if !startsUpper(words[i]) {
			continue
		}
		phrase := strings.Trim(strings.Join(words[i:i+3], " "), " ,.;:")
		if phrase == "" || strings.Contains(phrase, "*") {
			continue
		}
		idx := strings.Index(first, phrase)
		if idx < 0 {
			continue
		}
		newFirst := first[:idx] + "**" + phrase + "**" + first[idx+len(phrase):]
		return strings.Replace(summary, first, newFirst, 1), phrase
	}
	return summary, ""
}

func startsUpper(w string) bool {
	for _, r := range w {
		return unicode.IsUpper(r)
	}
	return false
}

// unitLength đo độ dài phần đơn vị đứng ngay sau con số ("  tỷ đồng").
func unitLength(rest string) int {
	trimmed := strings.TrimLeft(rest, " ")
	lead := len(rest) - len(trimmed)
	lower := strings.ToLower(trimmed)
	best := 0
	for _, u := range moneyUnits {
		if strings.HasPrefix(lower, u) && len(u) > best {
			best = len(u)
		}
	}
	if best == 0 {
		return 0
	}
	// Gom thêm danh từ liền sau đơn vị nếu có ("tấn thép", "tỷ đồng").
	return lead + best
}

func parseFloatSafe(s string) float64 {
	var f float64
	_, err := fmt.Sscanf(s, "%g", &f)
	if err != nil {
		return 0
	}
	return f
}
