package main

import (
	"fmt"
)

// This file is for solving task on the topic of two pointer by method "every one by pointer"
// Green flags "every one by pointer"
// - Даны несколько массивов или строк
// - Нужно искать объединение/пересичение и т.д. у двух последовательностей

// Даны два отсортированных массива, нужно вернуть все их общие элементы.
// Array one: [0, 2, 4, 8, 8]
// Array two: [1, 2, 2, 7, 8, 8, 8]
//
//	| p1
//	v
// [0, 2, 4, 8, 8]
//
//	| p2
//	v
// [1, 2, 2, 7, 8, 8, 8]
//
// Стваим 2 указателя p1 и p2 на начало двух массивов. После чего начинаем сравнивать
// p1 == p2. Если p1 == p2 -> добавляем в ответ и двигаем оба указателя,
// если p1 < p2 -> двигаем p1, если p1 > p2 -> двигаем p2. Повторяем пока один из
// указателей не выйдет за границы массива.
//
// Время: O(n+m)
// Память: O(min(n, m))

func bothElements(arr1 []int, arr2 []int) []int {
	// 1. Иницизилация
	p1 := 0
	p2 := 0
	answer := []int{}

	// 2. Цикл
	for p1 < len(arr1) && p2 < len(arr2) {
		if arr1[p1] == arr2[p2] {
			answer = append(answer, arr1[p1])
			p1++
			p2++
		} else if arr1[p1] < arr2[p2] {
			p1++
		} else if arr1[p1] > arr2[p2] {
			p2++
		}
	}
	return answer
}

// func main() {
// 	arr1 := []int{0, 2, 4, 8, 8}
// 	arr2 := []int{1, 2, 2, 7, 8, 8, 8}
// 	arr3 := []int{1, 3, 4, 7}
// 	arr4 := []int{0, 2, 5, 6}
// 	fmt.Println(bothElements(arr1, arr2), bothElements(arr3, arr4))
// }
// Даны два отсортированных массива nums1 и nums2. Нужно объединить их в один отсортированный массив.
//
// Ввод:
// nums1 = [1, 3, 5, 7], nums2 = [2, 4, 6, 8]
// Вывод:
// [1, 2, 3, 4, 5, 6, 7, 8]

func inOneSortedArray(nums1, nums2 []int) []int {
	p1 := 0
	p2 := 0
	answer := []int{}

	for p1 < len(nums1) && p2 < len(nums2) {
		if nums1[p1] == nums1[p2] {
			answer = append(answer, nums1[p1], nums2[p2])
			p1++
			p2++
		} else if nums1[p1] > nums2[p2] {
			answer = append(answer, nums2[p2])
			p2++
		} else if nums1[p1] < nums2[p2] {
			answer = append(answer, nums1[p1])
			p1++
		}
	}
	if len(nums1) > len(nums2) {
		answer = append(answer, nums1[len(nums2):]...)
	} else if len(nums1) < len(nums2) {
		answer = append(answer, nums2[len(nums1):]...)
	}
	return answer
}

// func main() {
// 	nums1 := []int{1, 3, 5, 7}
// 	nums2 := []int{2, 4, 6, 8}
// 	// nums1 := []int{1, 3, 5, 7}
// 	// nums2 := []int{2, 4, 6, 8, 9, 10, 11}
//
// 	fmt.Println(inOneSortedArray(nums1, nums2))
// }

// You are given two integer arrays nums1 and nums2, sorted in non-decreasing order,
// and two integers m and n, representing the number of elements in nums1 and nums2 respectively.
//
// Merge nums1 and nums2 into a single array sorted in non-decreasing order.
//
// The final sorted array should not be returned by the function, but instead be stored inside
// the array nums1. To accommodate this, nums1 has a length of m + n, where the first m elements
// denote the elements that should be merged, and the last n elements are set to 0 and should be
// ignored. nums2 has a length of n.
//
//
//
// Example 1:
//
// Input: nums1 = [1,2,3,0,0,0], m = 3, nums2 = [2,5,6], n = 3
// Output: [1,2,2,3,5,6]
// Explanation: The arrays we are merging are [1,2,3] and [2,5,6].
// The result of the merge is [1,2,2,3,5,6] with the underlined elements coming from nums1.
//
// Example 2:
//
// Input: nums1 = [1], m = 1, nums2 = [], n = 0
// Output: [1]
// Explanation: The arrays we are merging are [1] and [].
// The result of the merge is [1].
//
// Example 3:
//
// Input: nums1 = [0], m = 0, nums2 = [1], n = 1
// Output: [1]
// Explanation: The arrays we are merging are [] and [1].
// The result of the merge is [1].
// Note that because m = 0, there are no elements in nums1. The 0 is only there to ensure the merge result
// can fit in nums1.
// Constraints:
// nums1.length == m + n
// nums2.length == n
// 0 <= m, n <= 200
// 1 <= m + n <= 200
// -109 <= nums1[i], nums2[j] <= 109

func merge(nums1 []int, m int, nums2 []int, n int) {
	p1 := 0
	p2 := 0

	for len(nums1) > m+n && len(nums2) > n {
		fmt.Println(nums1, nums2, m, n, p1, p2)
		if nums1[p1] == 0 {
			nums1[p1], nums2[p2] = nums1[p2], nums2[p1]
			p1++
			p2++
		} else if nums1[p1] > nums2[p2] {
			nums1[p1], nums2[p2] = nums1[p2], nums2[p1]
			p2++
		} else if nums1[p1] == nums2[p2] {
			p1++
		} else if nums1[p1] < nums2[p2] {
			p2--
		}
	}
}

func main() {
	// nums1 := []int{1, 2, 3, 0, 0, 0}
	// nums2 := []int{2, 5, 6}
	// nums1 := []int{4, 5, 6, 0, 0, 0}
	// nums2 := []int{1, 2, 3}
	nums1 := []int{4, 0, 0, 0, 0, 0}
	nums2 := []int{1, 2, 3, 5, 6}
	fmt.Println(nums1)
	merge(nums1, len(nums1)-len(nums2), nums2, len(nums2))
	fmt.Println(nums1)
}
