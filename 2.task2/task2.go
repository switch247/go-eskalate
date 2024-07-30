package main

import (
	"fmt"
	"strings"
	"unicode"
)

func main() {
	fmt.Println(WordFrequencyCount("Hello  World! "))
	fmt.Println(IsPalindrome("??3tes et3!   "))
}

// helper
func clean(input string) string {
	var result string
	for _, char := range input {
		if unicode.IsLetter(char) || unicode.IsNumber(char) {
			result += string(char)
		}
	}
	return strings.ToLower(result)
}

// task 1

func WordFrequencyCount(input string) map[string]int {
	// thi s skip spaces
	frequency := make(map[string]int)
	parts := strings.Split(input, " ")
	for _, part := range parts {
		part = clean(part)
		if part != "" {
			fmt.Println(part)
			frequency[part]++
		}
	}
	return frequency
}

// Task 2

// (ignoring spaces, punctuation, and capitalization).

func IsPalindrome(input string) bool {
	input = clean(input)
	fmt.Println(input)
	length := len(input)
	for i := 0; i < length/2; i++ {
		if input[i] != input[length-i-1] {
			return false
		}
	}
	return true
}
