package main

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
	// sorting nums array
	for range len(nums) {
		for i := len(nums) - 1; i > 0; i-- {
			if nums[i] < nums[i-1] {
				nums[i], nums[i-1] = nums[i-1], nums[i]
			}
		}
	}

	var answer [][]int
	seen := make(map[[3]int]bool)

	for m := range len(nums) {
		l := 0
		r := len(nums) - 1
		for l < r {
			if l == m {
				l++
			} else if r == m {
				r--
			} else if nums[l]+nums[r] == 0-nums[m] {
				add := [3]int{nums[l], nums[m], nums[r]}
				if m < l {
					add[0], add[1] = add[1], add[0]
				}
				if m > r {
					add[1], add[2] = add[2], add[1]
				}
				if !seen[add] {
					answer = append(answer, add[:])
					seen[add] = true
				}
				r--
			} else if nums[l]+nums[r] < 0-nums[m] {
				l++
			} else {
				r--
			}
		}
	}
	return answer
}
