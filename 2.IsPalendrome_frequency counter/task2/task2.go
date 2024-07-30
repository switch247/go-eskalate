package task2

import (
	"fmt"
	"main/utils/cleaner"
)

// Task 2

// (ignoring spaces, punctuation, and capitalization).

func IsPalindrome(input string) bool {
	input = cleaner.Clean(input)
	fmt.Println(input)
	length := len(input)
	for i := 0; i < length/2; i++ {
		if input[i] != input[length-i-1] {
			return false
		}
	}
	return true
}
