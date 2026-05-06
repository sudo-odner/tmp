package main

import "fmt"

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

func main() {
	nums1 := []int{1, 3, 5, 7}
	nums2 := []int{2, 4, 6, 8}
	// nums1 := []int{1, 3, 5, 7}
	// nums2 := []int{2, 4, 6, 8, 9, 10, 11}

	fmt.Println(inOneSortedArray(nums1, nums2))
}
