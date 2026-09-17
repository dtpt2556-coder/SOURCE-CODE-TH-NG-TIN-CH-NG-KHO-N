package textutil

import (
	"regexp"
	"strconv"
	"strings"
)

// numberRe bắt mọi token số theo spec mục 9 architecture.md: [0-9][0-9.,]*
var numberRe = regexp.MustCompile(`[0-9][0-9.,]*`)

// NormalizeVNNumber chuẩn hoá số kiểu Việt Nam về dạng so sánh được:
//
//	"408.000" -> "408000"   (dấu "." phân nhóm hàng nghìn)
//	"1,67"    -> "1.67"     (dấu "," phân cách thập phân)
//
// Đây là bẫy lớn nhất khi dùng thư viện parse số kiểu Anh-Mỹ: `strconv` sẽ đọc
// "408.000" thành 408 (R3.6 của BA2, S-04).
func NormalizeVNNumber(raw string) (string, bool) {
	s := strings.Trim(raw, ".,")
	if s == "" {
		return "", false
	}
	s = strings.ReplaceAll(s, ".", "")
	s = strings.ReplaceAll(s, ",", ".")
	if strings.Count(s, ".") > 1 {
		return "", false
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return "", false
	}
	return strconv.FormatFloat(f, 'f', -1, 64), true
}

// RawNumbers trả về các token số nguyên văn trong văn bản.
func RawNumbers(text string) []string {
	return numberRe.FindAllString(text, -1)
}

// NumberSet trả về tập số đã chuẩn hoá của văn bản.
func NumberSet(text string) map[string]bool {
	out := map[string]bool{}
	for _, raw := range RawNumbers(text) {
		if norm, ok := NormalizeVNNumber(raw); ok {
			out[norm] = true
		}
	}
	return out
}

// NumbersDiffer cho biết hai văn bản có tập số liệu khác nhau hay không.
// Số liệu đổi là tín hiệu mạnh nhất để quyết định tóm tắt lại (U-03).
func NumbersDiffer(before, after string) bool {
	a, b := NumberSet(before), NumberSet(after)
	if len(a) != len(b) {
		return true
	}
	for k := range a {
		if !b[k] {
			return true
		}
	}
	return false
}
