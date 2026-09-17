// Package shortlink sinh và phân giải mã rút gọn /r/{code}
// (mục 7 architecture.md, §7 BA2).
package shortlink

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// CodeLength — base62 6 ký tự = 56,8 tỷ tổ hợp (mục 7 architecture.md).
const CodeLength = 6

const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// GenerateCode sinh mã base62 ngẫu nhiên bằng crypto/rand.
// Mã phải KHÔNG đoán được và KHÔNG tuần tự để chống enumeration (§7.4 BA2).
func GenerateCode() (string, error) {
	return generate(CodeLength)
}

func generate(n int) (string, error) {
	max := big.NewInt(int64(len(alphabet)))
	buf := make([]byte, n)
	for i := 0; i < n; i++ {
		idx, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", fmt.Errorf("read from crypto/rand: %w", err)
		}
		buf[i] = alphabet[idx.Int64()]
	}
	return string(buf), nil
}

// IsValidCode kiểm tra định dạng mã trước khi truy vấn DB — chặn sớm các
// request rác của scanner.
func IsValidCode(code string) bool {
	if len(code) != CodeLength {
		return false
	}
	for i := 0; i < len(code); i++ {
		c := code[i]
		switch {
		case c >= '0' && c <= '9':
		case c >= 'A' && c <= 'Z':
		case c >= 'a' && c <= 'z':
		default:
			return false
		}
	}
	return true
}
