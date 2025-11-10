package main

import (
	"testing"
)

func TestReverseString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Simple word", "hello", "olleh"},
		{"Sentence with spaces", "Go is fun!", "!nuf si oG"},
		{"Empty string", "", ""},
		{"Palindrome", "madam", "madam"},
		{"Special characters", "12345!@#$%", "%$#@!54321"},
		{"Mixed case", "GoLang", "gnaLoG"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ReverseString(tt.input)
			if result != tt.expected {
				t.Errorf("ReverseString(%q) = %q; want %q", tt.input, result, tt.expected)
			}
		})
	}
}
