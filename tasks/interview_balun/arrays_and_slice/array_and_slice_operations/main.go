package main

import (
	"fmt"
	"unsafe"
)

func accessToArrayElement1() {
	data := [3]int{1, 2, 3}

	idx := 4
	fmt.Println(data[idx]) // panic: index out of range (non-constant)

	// fmt.Println(data[4]) // comilation error: out of bounds (constant compiller veiw 3 < 4)
}

func accessToArrayElement2() {
	data := [3]int{1, 2, 3}

	idx := -1
	fmt.Println(data[idx]) // panic: index out of range

	// fmt.Println(data[-1]) // comilation error: ndex must not negative
}

func arrayLen() {
	data := [10]int{}
	fmt.Println(len(data)) // 10
}

func capArray() {
	var data [10]int
	fmt.Println(cap(data)) // 10
}

// For an array, len and cap are always equal to its fixed size

func arraysComparison() {
	first := [...]int{1, 2, 3}
	second := [...]int{1, 2, 3}

	fmt.Println(first == second) // true
	fmt.Println(first != second) // false

	// fmt.Println(first < second) // comilation error: invalid operation
}

// An array is a comparable type; you can compare values with it, but the <, > and <=, >= operation are disable for arrays

func emptyArray() {
	var data [10]byte
	fmt.Println(unsafe.Sizeof(data)) // 10

	// fmt.Println(data == nil) // comilation error: invalid operation: mismatched type [10]byte and nil
}
