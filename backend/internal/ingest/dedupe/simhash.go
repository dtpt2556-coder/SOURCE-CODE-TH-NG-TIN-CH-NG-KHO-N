package dedupe

import (
	"encoding/hex"
	"hash/fnv"
	"math/bits"
	"strconv"
	"strings"

	"github.com/tenpoint/tenpoint-api/internal/textutil"
)

// MaxTitleHammingDistance — ngưỡng coi 2 tiêu đề là trùng (§6.2 L2 của BA2).
const MaxTitleHammingDistance = 3

// titleStopwords bị loại trước khi hash để tiêu đề xào lại vẫn hội tụ.
var titleStopwords = map[string]bool{
	"cua": true, "va": true, "trong": true, "cho": true, "voi": true,
	"ve": true, "la": true, "cac": true, "nhung": true, "tai": true,
	"tu": true, "den": true, "mot": true, "co": true, "da": true,
}

// ngramSize — kích thước n-gram ký tự dùng làm đặc trưng phụ.
const ngramSize = 3

// wordWeight — trọng số của đặc trưng từ so với n-gram ký tự.
const wordWeight = 3

// TitleSimhash tính SimHash 64-bit của tiêu đề.
//
// Đặc trưng = từ (trọng số 3) + n-gram ký tự (trọng số 1). Tiêu đề chỉ có
// 8–15 từ nên SimHash thuần theo từ quá nhiễu; bổ sung n-gram ký tự làm chữ ký
// ổn định hơn nhiều với các thay đổi nhỏ (thêm/bớt vài ký tự).
//
// Bắt buộc bỏ dấu trước khi hash: "Hoà Phát" và "Hòa Phát" là hai chuỗi Unicode
// khác nhau hoàn toàn nhưng cùng một tiêu đề (§6.3 BA2).
func TitleSimhash(title string) uint64 {
	weights := titleFeatures(title)
	if len(weights) == 0 {
		return 0
	}

	var vec [64]int
	for tok, w := range weights {
		h := hashToken(tok)
		for i := 0; i < 64; i++ {
			if h&(1<<uint(i)) != 0 {
				vec[i] += w
			} else {
				vec[i] -= w
			}
		}
	}

	var out uint64
	for i := 0; i < 64; i++ {
		if vec[i] > 0 {
			out |= 1 << uint(i)
		}
	}
	return out
}

// HammingDistance đếm số bit khác nhau giữa hai simhash.
func HammingDistance(a, b uint64) int {
	return bits.OnesCount64(a ^ b)
}

// IsNearDuplicate cho biết hai tiêu đề có thuộc cùng một cụm trùng hay không.
func IsNearDuplicate(a, b uint64) bool {
	return HammingDistance(a, b) <= MaxTitleHammingDistance
}

// EncodeSimhash chuyển simhash sang chuỗi hex 16 ký tự để lưu DB.
func EncodeSimhash(h uint64) string {
	var buf [8]byte
	for i := 0; i < 8; i++ {
		buf[7-i] = byte(h >> (8 * uint(i)))
	}
	return hex.EncodeToString(buf[:])
}

// DecodeSimhash đọc lại simhash từ chuỗi hex.
func DecodeSimhash(s string) (uint64, error) {
	return strconv.ParseUint(strings.TrimSpace(s), 16, 64)
}

// titleFeatures dựng tập đặc trưng có trọng số của tiêu đề.
func titleFeatures(title string) map[string]int {
	tokens := tokenizeTitle(title)
	if len(tokens) == 0 {
		return nil
	}
	weights := make(map[string]int, len(tokens)*4)
	for _, tok := range tokens {
		weights["w:"+tok] += wordWeight
	}
	runes := []rune(" " + strings.Join(tokens, " ") + " ")
	for i := 0; i+ngramSize <= len(runes); i++ {
		weights["g:"+string(runes[i:i+ngramSize])]++
	}
	return weights
}

func tokenizeTitle(title string) []string {
	raw := textutil.Tokenize(textutil.Normalize(title))
	out := make([]string, 0, len(raw))
	for _, tok := range raw {
		if titleStopwords[tok] {
			continue
		}
		if len([]rune(tok)) < 2 {
			continue
		}
		out = append(out, tok)
	}
	return out
}

func hashToken(tok string) uint64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(tok))
	return h.Sum64()
}
