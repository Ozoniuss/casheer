package stringutil

import (
	"strings"
	"testing"
)

/*
Tests generated using Chat-GPT, which actually caught an error. Prompt below.
*/

func TestCapitalize(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "Hello"},
		{"world", "World"},
		{"gO", "GO"},
		{"", ""},
		{"123", "123"},
		{"@#$", "@#$"},
	}

	b := &strings.Builder{}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := capitalize(tt.input, b)
			if result != tt.expected {
				t.Errorf("Expected '%s', but got '%s'", tt.expected, result)
			}
		})
	}
}

func TestCapitalizeArray(t *testing.T) {
	tests := []struct {
		input    []string
		expected []string
	}{
		{
			[]string{"hello", "world"},
			[]string{"Hello", "World"},
		},
		{
			[]string{"go", "gopher", "golang"},
			[]string{"Go", "Gopher", "Golang"},
		},
		{
			[]string{"123", "@#$", "test"},
			[]string{"123", "@#$", "Test"},
		},
		{
			[]string{},
			[]string{},
		},
	}

	for _, tt := range tests {
		t.Run(strings.Join(tt.input, ","), func(t *testing.T) {
			CapitalizeArray(tt.input)
			if len(tt.input) != len(tt.expected) {
				t.Errorf("Expected length %d, but got length %d", len(tt.expected), len(tt.input))
			}

			for i := 0; i < len(tt.input); i++ {
				if tt.input[i] != tt.expected[i] {
					t.Errorf("Expected '%s', but got '%s'", tt.expected[i], tt.input[i])
				}
			}
		})
	}
}

/* Chat GPT prompt
---
Here is some go code:

import "strings"

func capitalize(s string, b *strings.Builder) string {

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
func CapitalizeArray(values []string) {
	b := strings.Builder{}
	for i := 0; i < len(values); i++ {
		values[i] = capitalize(values[i], &b)
	}
}

Please write unit tests for these two functions in a go idiomatic way. Do not use the testify or testify/assert libraries. Use go standard's "testing" library. Write the tests with a table-driven approach.
---
*/
