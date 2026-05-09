package main

import (
	"fmt"
	"sort"
)

// Given an integer array nums of length n and an integer target, find three integers at distinct indices in nums such that the sum is closest to target.
//
// Return the sum of the three integers.
//
// You may assume that each input would have exactly one solution.
//
// Example 1:
//
// Input: nums = [-1,2,1,-4], target = 1
// Output: 2
// Explanation: The sum that is closest to the target is 2. (-1 + 2 + 1 = 2).
//
// Example 2:
//
// Input: nums = [0,0,0], target = 1
// Output: 0
// Explanation: The sum that is closest to the target is 0. (0 + 0 + 0 = 0).
//
// Constraints:
//
//	3 <= nums.length <= 500
//	-1000 <= nums[i] <= 1000
//	-104 <= target <= 104
func abs(i int) uint {
	if i < 0 {
		return uint(-1 * i)
	}
	return uint(i)
}

func threeSumClosest(nums []int, target int) int {
	var a int
	sort.Ints(nums)
	n := len(nums)
	md := ^uint(0)

	for l := range n - 2 {
		m, r := l+1, n-1
		for m < r {
			sum := nums[l] + nums[m] + nums[r]
			diff := sum - target
			fmt.Printf("%v %v %v   %v: best diff %v\n", nums[l], nums[m], nums[r], sum, diff)
			if diff == 0 {
				return target
			} else if diff < 0 {
				m++
			} else {
				r--
			}
			if abs(diff) < md {
				a = sum
				md = abs(diff)
			}
		}
	}
	return a
}
