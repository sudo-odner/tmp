package main

import (
	"sort"
)

// Given an integer array nums, return all the triplets [nums[i], nums[j], nums[k]] such that i != j, i != k, and j != k, and nums[i] + nums[j] + nums[k] == 0.
// Notice that the solution set must not contain duplicate triplets.
//
//
// Example 1:
//
// Input: nums = [-1,0,1,2,-1,-4]
// Output: [[-1,-1,2],[-1,0,1]]
// Explanation:
// nums[0] + nums[1] + nums[2] = (-1) + 0 + 1 = 0.
// nums[1] + nums[2] + nums[4] = 0 + 1 + (-1) = 0.
// nums[0] + nums[3] + nums[4] = (-1) + 2 + (-1) = 0.
// The distinct triplets are [-1,0,1] and [-1,-1,2].
// Notice that the order of the output and the order of the triplets does not matter.
//
// Example 2:
//
// Input: nums = [0,1,1]
// Output: []
// Explanation: The only possible triplet does not sum up to 0.
//
// Example 3:
//
// Input: nums = [0,0,0]
// Output: [[0,0,0]]
// Explanation: The only possible triplet sums up to 0.
//
//
// Constraints:
//     3 <= nums.length <= 3000
//     -105 <= nums[i] <= 105

func threeSum(nums []int) [][]int {
	var answer [][]int

	// sort nums array
	sort.Ints(nums)
	n := len(nums)

	for l := range n - 2 {
		if nums[l] > 0 {
			break
		}
		if l > 0 && nums[l] == nums[l-1] {
			continue
		}

		m, r := l+1, n-1
		for m < r {
			sum := nums[l] + nums[m] + nums[r]
			if sum == 0 {
				answer = append(answer, []int{nums[l], nums[m], nums[r]})
				for m < r && nums[m] == nums[m+1] {
					m++
				}
				for m < r && nums[r] == nums[r-1] {
					r--
				}
				m++
				r--
			} else if sum < 0 {
				m++
			} else {
				r--
			}
		}
	}
	return answer
}
