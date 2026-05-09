package main

import (
	"reflect"
	"testing"
)

func TestSolution(t *testing.T) {
	tests := []struct {
		name   string
		nums   []int
		output [][]int
	}{
		{
			name:   "1 Task",
			nums:   []int{-1, 0, 1, 2, -1, -4},
			output: [][]int{{-1, -1, 2}, {-1, 0, 1}},
		},
		{
			name:   "2 Task",
			nums:   []int{-1, 0, 1, 0},
			output: [][]int{{-1, 0, 1}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := threeSum(tt.nums)
			if !reflect.DeepEqual(out, tt.output) {
				t.Errorf("\t got %v,\n\t want %v\n", out, tt.output)
			}
		})
	}
}
