package common

import (
	"strings"
	"unicode"
)

func Tokenize(expression string) []string {
	var tokens []string
	var numBuffer strings.Builder

	for _, char := range expression {
		if unicode.IsDigit(char) || char == '.' {
			numBuffer.WriteRune(char)
		} else if char == ' ' {
			if numBuffer.Len() > 0 {
				tokens = append(tokens, numBuffer.String())
				numBuffer.Reset()
			}
		} else {
			if numBuffer.Len() > 0 {
				tokens = append(tokens, numBuffer.String())
				numBuffer.Reset()
			}
			tokens = append(tokens, string(char))
		}
	}
	if numBuffer.Len() > 0 {
		tokens = append(tokens, numBuffer.String())
	}
	return tokens
}
