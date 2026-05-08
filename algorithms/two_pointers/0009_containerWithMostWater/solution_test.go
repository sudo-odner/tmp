package main

import "testing"

func TestSolution(t *testing.T) {
	tests := []struct {
		name   string
		height []int
		output int
	}{
		{
			name:   "1 Task",
			height: []int{1, 1},
			output: 1,
		},
		{
			name:   "2 Task",
			height: []int{1, 8, 6, 2, 5, 4, 8, 3, 7},
			output: 49,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := maxArea(tt.height)
			if out != tt.output {
				t.Errorf("maxArea() is %v, want %v", out, tt.output)
			}
		})
	}
}
