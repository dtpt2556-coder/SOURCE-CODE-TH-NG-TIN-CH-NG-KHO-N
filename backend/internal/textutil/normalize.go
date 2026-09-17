// Package textutil gom các hàm chuẩn hoá văn bản tiếng Việt dùng chung cho
// tagger, dedupe và extractor. Không phụ thuộc hạ tầng, test được offline.
package textutil

import (
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

var dReplacer = strings.NewReplacer("đ", "d", "Đ", "D")

// RemoveDiacritics bỏ dấu tiếng Việt: "Hoà Phát" -> "Hoa Phat".
//
// Tách NFD rồi lọc rune thủ công thay vì dùng transform.Chain: chuỗi chain
// NFD -> runes.Remove -> NFC panic "slice bounds out of range" trên một số body
// crawl về từ báo VN, và một panic ở đây giết cả goroutine của nguồn đó.
// Đ/đ phải thay riêng vì hai ký tự này không phân rã được bằng NFD.
func RemoveDiacritics(s string) string {
	if s == "" {
		return ""
	}
	s = dReplacer.Replace(strings.ToValidUTF8(s, ""))

	decomposed := norm.NFD.String(s)
	var b strings.Builder
	b.Grow(len(decomposed))
	for _, r := range decomposed {
		if unicode.Is(unicode.Mn, r) {
			continue
		}
		b.WriteRune(r)
	}
	return norm.NFC.String(b.String())
}

// Normalize = bỏ dấu + lowercase + gom khoảng trắng.
// Đây là dạng khoá dùng cho `ticker_aliases.alias_normalized`.
func Normalize(s string) string {
	return CollapseSpaces(strings.ToLower(RemoveDiacritics(s)))
}

// CollapseSpaces gom mọi chuỗi khoảng trắng thành đúng một dấu cách và trim.
func CollapseSpaces(s string) string {
	return strings.Join(strings.FieldsFunc(s, unicode.IsSpace), " ")
}

// Fold chuẩn hoá văn bản để SO KHỚP CỤM TỪ mà vẫn GIỮ NGUYÊN DẤU: hợp nhất
// Unicode về NFC, hạ chữ thường, thay mọi ký tự không phải chữ/số (trừ "/")
// bằng dấu cách rồi gom khoảng trắng.
//
// Đây là hàm chuẩn cho mọi phép so khớp mang ngữ nghĩa — cụm bị cấm, alias tên
// doanh nghiệp, từ khoá phân loại.
//
// Bỏ dấu là phép biến đổi LÀM MẤT THÔNG TIN nên không được dùng ở đây: trong
// tiếng Việt dấu phân biệt nghĩa. "nền tảng" và "nên tăng" cùng ra "nen tang",
// "căn bản" và "cần bán" cùng ra "can ban", "Khi Việt Nam" và "Khí Việt Nam"
// cùng ra "khi viet nam". Ba va chạm đó đã thật sự cách ly oan 6 bài và gắn sai
// mã GAS trên dữ liệu production.
//
// Chỉ dùng RemoveDiacritics/Normalize cho tìm kiếm do NGƯỜI DÙNG gõ (US-2.4),
// nơi người ta chủ động gõ không dấu và chấp nhận kết quả rộng hơn.
func Fold(s string) string {
	if s == "" {
		return ""
	}
	folded := strings.ToLower(norm.NFC.String(strings.ToValidUTF8(s, "")))

	var b strings.Builder
	b.Grow(len(folded))
	for _, r := range folded {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '/' {
			b.WriteRune(r)
		} else {
			b.WriteRune(' ')
		}
	}
	return CollapseSpaces(b.String())
}

// FoldPad bọc Fold bằng dấu cách hai đầu để kiểm tra biên từ bằng Contains.
func FoldPad(s string) string {
	return " " + Fold(s) + " "
}

// ContainsPhrase kiểm tra `phrase` xuất hiện trọn vẹn theo biên từ trong văn bản
// đã Fold. Cả hai đầu vào đều giữ dấu.
func ContainsPhrase(foldedPadded, phrase string) bool {
	return CountPhrase(foldedPadded, Fold(phrase)) > 0
}

// CountPhrase đếm số lần `core` xuất hiện như một cụm từ hoàn chỉnh trong chuỗi
// đã được bọc dấu cách hai đầu.
func CountPhrase(padded, core string) int {
	if core == "" {
		return 0
	}
	n, start := 0, 0
	for start < len(padded) {
		i := strings.Index(padded[start:], core)
		if i < 0 {
			break
		}
		abs := start + i
		end := abs + len(core)
		leftOK := abs == 0 || padded[abs-1] == ' '
		rightOK := end == len(padded) || padded[end] == ' '
		if leftOK && rightOK {
			n++
		}
		start = abs + 1
	}
	return n
}

// NormalizeWords = Normalize nhưng thay mọi ký tự không phải chữ/số bằng dấu
// cách. CHỈ dùng cho tìm kiếm không dấu của người dùng, không dùng để so khớp
// cụm từ mang ngữ nghĩa — xem Fold.
func NormalizeWords(s string) string {
	n := Normalize(s)
	var b strings.Builder
	b.Grow(len(n))
	for _, r := range n {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
		} else {
			b.WriteRune(' ')
		}
	}
	return CollapseSpaces(b.String())
}

// Tokenize tách chuỗi đã normalize thành các token chữ/số.
func Tokenize(s string) []string {
	return strings.FieldsFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
}

// FirstParagraph trả về đoạn đầu tiên của body (ngăn cách bằng dòng trống hoặc
// xuống dòng). Dùng cho rule phân định primary/mentioned.
func FirstParagraph(body string) string {
	body = strings.TrimSpace(body)
	if body == "" {
		return ""
	}
	if idx := strings.Index(body, "\n\n"); idx >= 0 {
		return strings.TrimSpace(body[:idx])
	}
	if idx := strings.Index(body, "\n"); idx >= 0 {
		return strings.TrimSpace(body[:idx])
	}
	return body
}
