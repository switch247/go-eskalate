package task1

import (
	"fmt"
	"main/utils/cleaner"
	"strings"
)

// task 1

func WordFrequencyCount(input string) map[string]int {
	// thi s skip spaces
	frequency := make(map[string]int)
	parts := strings.Split(input, " ")
	for _, part := range parts {
		part = cleaner.Clean(part)
		if part != "" {
			fmt.Println(part)
			frequency[part]++
		}
	}
	return frequency
}
