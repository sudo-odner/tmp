package main

import "fmt"

// This file is for solving task on the topic of two pointer by method "from both side"
// Green flags "from both side"
// - по условию дан один отсторированниый массив
// - задача на проверку полиндрома (возможны усложнения)
// - ответ формируется за счет сужения области с двух сторон

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
	// 1. Иницилизация указателей
	l := 0
	r := len(nums) - 1

	// 2. Цикл
	for l < r {
		// 3. Логика движения указателя
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

// func main() {
// 	nums := []int{-2, 1, 6, 9, 12, 21}
// 	target := 18
// 	answer := two_sum(nums, target)
// 	fmt.Println(answer)
// }

// Дано слово s. Нужно проверить, является ли оно палиндромом. Палиндром - слово,
// которе читается слева на право и справва на левао одинаково. Например, "шалаш" -
// полиндром, а "кот" - нет
//
// Ввод:
// s = "шалаш"
// Вывод:
// true

func palindrom(s string) bool {
	sRune := []rune(s)
	l := 0
	r := len(sRune) - 1

	for l < r {
		if sRune[l] != sRune[r] {
			return false
		}
		l++
		r--
	}
	return true
}

func main() {
	shalash := "шалаш"
	cat := "кот"
	fmt.Println(palindrom(shalash), palindrom(cat))
}
