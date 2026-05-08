package main

import "testing"

func TestSolution(t *testing.T) {
	tests := []struct {
		s   string
		out string
	}{
		{
			s:   "babad",
			out: "bab",
		},
		{
			s:   "cbbd",
			out: "bb",
		},
		{
			s:   "ccc",
			out: "ccc",
		},
		{
			s:   "abcda",
			out: "a",
		},
	}

	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			out := longestPalindrome(tt.s)
			if out != tt.out {
				t.Errorf("longestPalindrome() = %s, want %s", out, tt.out)
			}
		})
	}
}
