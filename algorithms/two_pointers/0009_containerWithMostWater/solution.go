package main

func maxArea(height []int) int {
	p1 := 0
	p2 := len(height) - 1
	maxStore := 0
	var newMax int
	for p1 != p2 {
		if height[p1] > height[p2] {
			newMax = height[p2] * (p2 - p1)
			if newMax > maxStore {
				maxStore = newMax
			}
			p2--
		} else {
			newMax = height[p1] * (p2 - p1)
			if newMax > maxStore {
				maxStore = newMax
			}
			p1++
		}
	}

	return maxStore
}
