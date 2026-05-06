package main

import "fmt"

// This file is for solving task on the topic of two pointer by method "from both side"

// Дан отсторированниый массив и число turget, нужно вернуть позиции 2 чисел,
// которые даеют в сумме target.
// Array: [-2, 1, 6, 9, 12, 21]
// Target: 18
//
//   |                |
//   v                v
// [-2, 1, 6, 9, 12, 21]
// Для этого мы испольуем два указателя, ставим их в начало и конец.
// После чего суммируем числа (currSum) и сравниваем с target.
// Если currSum > target -> дваигаем правый указатель, если currSum < target двигаем
// левый указатль. Если currSum == target -> выводим его, если указатеи пересеклись
// -> нет решения
//
// Время: O(n)
// Память: O(1)

func two_sum(nums []int, target int) []int {
	l := 0
	r := len(nums) - 1

	for l < r {
		currSum := nums[l] + nums[r]
		if currSum == target {
			return []int{l, r}
		}
		if currSum > target {
			r--
		}
		if currSum < target {
			l++
		}
	}
	return []int{-1, -1}
}

func main() {
	nums := []int{-2, 1, 6, 9, 12, 21}
	target := 18
	answer := two_sum(nums, target)
	fmt.Println(answer)
}
