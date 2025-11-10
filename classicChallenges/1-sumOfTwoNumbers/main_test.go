package main

import (
	"testing"
)

func TestSum(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"Positive numbers", 2, 3, 5},
		{"Zero values", 0, 0, 0},
		{"Negative numbers", -2, -3, -5},
		{"Mixed signs", -5, 10, 5},
		{"Large numbers", 1000000000, 1000000000, 2000000000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Sum(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Sum(%d, %d) = %d; want %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}
