package stringutil

import "strings"

func capitalize(s string, b *strings.Builder) string {

	if len(s) == 0 {
		return s
	}

	if s[0] > 'z' || s[0] < 'a' {
		return s
	}

	// Builder is reused across multiple calls.
	defer b.Reset()
	b.WriteByte(s[0] - 32)
	for i := 1; i < len(s); i++ {
		b.WriteByte(s[i])
	}
	return b.String()
}

// CapitalizeArray takes an array of string values representing database column
// names consisting of ASCII characters, and returns the same array with each
// word having the starting letter capitalized.
func CapitalizeArray(values []string) []string {
	out := make([]string, len(values))
	b := strings.Builder{}
	for i := 0; i < len(values); i++ {
		out[i] = capitalize(values[i], &b)
	}
	return out
}
