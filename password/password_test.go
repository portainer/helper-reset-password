package password

import (
	"strings"
	"testing"
	"unicode"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	symbols = "~!@#$%^&*()_+`-={}|[]\\:\"<>?,./"
)

func TestGeneratePlainTextPassword(t *testing.T) {
	t.Parallel()
	t.Run("returns no error", func(t *testing.T) {
		_, err := GeneratePlainTextPassword()
		require.NoError(t, err)
	})

	t.Run("returns a non-empty password", func(t *testing.T) {
		pass, err := GeneratePlainTextPassword()
		require.NoError(t, err)
		assert.NotEmpty(t, pass)
	})

	t.Run("returns a password of length 32", func(t *testing.T) {
		pass, err := GeneratePlainTextPassword()
		require.NoError(t, err)
		assert.Len(t, pass, 32)
	})

	t.Run("contains at least 8 digits", func(t *testing.T) {
		pass, err := GeneratePlainTextPassword()
		require.NoError(t, err)

		digitCount := 0
		for _, ch := range pass {
			if unicode.IsDigit(ch) {
				digitCount++
			}
		}
		assert.GreaterOrEqual(t, digitCount, 8, "expected at least 8 digits, got %d in %q", digitCount, pass)
	})

	t.Run("contains at least 8 symbols", func(t *testing.T) {
		pass, err := GeneratePlainTextPassword()
		require.NoError(t, err)

		symbolCount := 0
		for _, ch := range pass {
			if strings.ContainsRune(symbols, ch) {
				symbolCount++
			}
		}
		assert.GreaterOrEqual(t, symbolCount, 8, "expected at least 8 symbols, got %d in %q", symbolCount, pass)
	})

	t.Run("generates unique passwords on repeated calls", func(t *testing.T) {
		const iterations = 5
		generated := make(map[string]struct{}, iterations)
		for range iterations {
			pass, err := GeneratePlainTextPassword()
			require.NoError(t, err)
			generated[pass] = struct{}{}
		}
		assert.Greater(t, len(generated), 1, "expected unique passwords across %d calls", iterations)
	})
}

func TestGenerateRandomString(t *testing.T) {
	t.Parallel()
	t.Run("returns no error", func(t *testing.T) {
		_, err := GenerateRandomString()
		require.NoError(t, err)
	})

	t.Run("returns a non-empty string", func(t *testing.T) {
		s, err := GenerateRandomString()
		require.NoError(t, err)
		assert.NotEmpty(t, s)
	})

	t.Run("returns a string of length 16", func(t *testing.T) {
		s, err := GenerateRandomString()
		require.NoError(t, err)
		assert.Len(t, s, 16)
	})

	t.Run("contains at least 8 digits", func(t *testing.T) {
		s, err := GenerateRandomString()
		require.NoError(t, err)

		digitCount := 0
		for _, ch := range s {
			if unicode.IsDigit(ch) {
				digitCount++
			}
		}
		assert.GreaterOrEqual(t, digitCount, 8, "expected at least 8 digits, got %d in %q", digitCount, s)
	})

	t.Run("contains no symbols", func(t *testing.T) {
		s, err := GenerateRandomString()
		require.NoError(t, err)

		for _, ch := range s {
			assert.False(t, strings.ContainsRune(symbols, ch), "unexpected symbol %q found in %q", ch, s)
		}
	})

	t.Run("contains no uppercase letters", func(t *testing.T) {
		s, err := GenerateRandomString()
		require.NoError(t, err)

		for _, ch := range s {
			assert.False(t, unicode.IsUpper(ch), "unexpected uppercase letter %q found in %q", ch, s)
		}
	})

	t.Run("generates unique strings on repeated calls", func(t *testing.T) {
		const iterations = 5
		generated := make(map[string]struct{}, iterations)
		for range iterations {
			s, err := GenerateRandomString()
			require.NoError(t, err)
			generated[s] = struct{}{}
		}
		assert.Greater(t, len(generated), 1, "expected unique strings across %d calls", iterations)
	})
}
