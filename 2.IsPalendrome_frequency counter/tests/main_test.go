package tests

import (
	"main/task1"
	"main/task2"
	"main/utils/cleaner"
	"reflect"
	"testing"
)

func TestClean(t *testing.T) {
	tests := []struct {
		input  string
		output string
	}{
		{"Hello, World!", "helloworld"},
		{"123abc", "123abc"},
		{"!@#$%^&*()_+-=", ""},
		{"", ""},
	}

	for _, test := range tests {
		if result := cleaner.Clean(test.input); result != test.output {
			t.Errorf("Expected %q, but got %q", test.output, result)
		} else {
			t.Logf("Test passed for input: %q", test.input)
		}
	}
}

func TestWordFrequencyCount(t *testing.T) {
	tests := []struct {
		input  string
		output map[string]int
	}{
		{"Hello, World! Hello", map[string]int{"hello": 2, "world": 1}},
		{"123abc 456def", map[string]int{"123abc": 1, "456def": 1}},
		{"!@#$%^&*()_+-=", map[string]int{}},
	}

	for _, test := range tests {
		if frequency := task1.WordFrequencyCount(test.input); !reflect.DeepEqual(frequency, test.output) {
			t.Errorf("Expected %v, but got %v", test.output, frequency)
		}
	}
}

func TestIsPalindrome(t *testing.T) {
	tests := []struct {
		input  string
		output bool
	}{
		{"madam", true},
		{"hello", false},
		{"", true},
		{"a", true},
		{"abc", false},
	}

	for _, test := range tests {
		if result := task2.IsPalindrome(test.input); result != test.output {
			t.Errorf("Expected %t, but got %t", test.output, result)
		}
	}
}
