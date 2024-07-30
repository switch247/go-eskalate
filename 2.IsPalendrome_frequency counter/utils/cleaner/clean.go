package cleaner

import (
	"strings"
	"unicode"
)

// helper
func Clean(input string) string {
	var result string
	for _, char := range input {
		if unicode.IsLetter(char) || unicode.IsNumber(char) {
			result += string(char)
		}
	}
	return strings.ToLower(result)
}
