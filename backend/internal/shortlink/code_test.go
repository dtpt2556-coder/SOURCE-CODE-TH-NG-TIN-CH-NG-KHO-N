package shortlink_test

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tenpoint/tenpoint-api/internal/shortlink"
)

var base62Re = regexp.MustCompile(`^[0-9A-Za-z]{6}$`)

func TestGenerateCodeFormat(t *testing.T) {
	for i := 0; i < 100; i++ {
		code, err := shortlink.GenerateCode()
		require.NoError(t, err)
		assert.Len(t, code, shortlink.CodeLength)
		assert.Regexp(t, base62Re, code)
		assert.True(t, shortlink.IsValidCode(code))
	}
}

func TestGenerateCodeIsNotSequential(t *testing.T) {
	// Mã phải ngẫu nhiên (crypto/rand) để chống enumeration (§7.4 BA2).
	seen := make(map[string]bool, 2000)
	for i := 0; i < 2000; i++ {
		code, err := shortlink.GenerateCode()
		require.NoError(t, err)
		require.False(t, seen[code], "trùng mã sau %d lần sinh: %s", i, code)
		seen[code] = true
	}
	assert.Len(t, seen, 2000)
}

func TestIsValidCode(t *testing.T) {
	valid := []string{"a7Kx2p", "000000", "ZzZzZz", "b3Qw9m"}
	for _, c := range valid {
		assert.True(t, shortlink.IsValidCode(c), c)
	}

	invalid := []string{"", "abc", "abcdefg", "a7Kx2-", "a7Kx2/", "a7Kx 2", "../../x"}
	for _, c := range invalid {
		assert.False(t, shortlink.IsValidCode(c), c)
	}
}
