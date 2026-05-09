package main

import "testing"

func TestSolution(t *testing.T) {
	tests := []struct {
		name   string
		nums   []int
		target int
		output int
	}{
		{
			name:   "1 Task",
			nums:   []int{-1, 2, 1, -4},
			target: 1,
			output: 2,
		},
		{
			name:   "2 Task",
			nums:   []int{0, 0, 0},
			target: 1,
			output: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := threeSumClosest(tt.nums, tt.target)
			if out != tt.output {
				t.Errorf("got %v, want %v", out, tt.output)
			}
		})
	}
}
