package main

import "testing"

func TestAtCoder(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "single parentheses",
			input:    "OMRON Corporation Programming Contest 2025 #2 (AtCoder Beginner Contest 432)",
			expected: "AtCoder Beginner Contest 432",
		},
		{
			name:     "no parentheses",
			input:    "AtCoder Beginner Contest 430",
			expected: "AtCoder Beginner Contest 430",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := validateAtCoder(tc.input)
			if actual != tc.expected {
				t.Errorf("Input: %v\nExpect: %v\nActual: %v", tc.input, tc.expected, actual)
			}
		})
	}
}

func TestCodeforces(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "single parentheses",
			input:    "Educational Codeforces Round 182 (Rated for Div. 2)",
			expected: "Educational Codeforces Round 182",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := validateCodeforces(tc.input)
			if actual != tc.expected {
				t.Errorf("Input: %v\nExpected: %v\nnActual: %v\n", tc.input, tc.expected, actual)
			}
		})
	}
}
