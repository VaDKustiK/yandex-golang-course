package common

import (
	"strings"
	"unicode"
)

// Tokenize разбивает выражение на токены.
func Tokenize(expression string) []string {
	var tokens []string
	var buf strings.Builder

	for _, r := range expression {
		if unicode.IsDigit(r) || r == '.' {
			buf.WriteRune(r)
		} else {
			if buf.Len() > 0 {
				tokens = append(tokens, buf.String())
				buf.Reset()
			}
			if !unicode.IsSpace(r) {
				tokens = append(tokens, string(r))
			}
		}
	}
	if buf.Len() > 0 {
		tokens = append(tokens, buf.String())
	}
	return tokens
}
