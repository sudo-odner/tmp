package main

import (
	"fmt"
	"unsafe"
)

func accessToArrayElement1() {
	data := [3]int{1, 2, 3}

	idx := 4
	fmt.Println(data[idx]) // panic: index out of range (non-constant)

	fmt.Println(data[4]) // compilation error: out of bounds (constant compiller veiw 3 < 4)
}

func accessToArrayElement2() {
	data := [3]int{1, 2, 3}

	idx := -1
	fmt.Println(data[idx]) // panic: index out of range

	fmt.Println(data[-1]) // compilation error: index must not negative
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

	// An array is a comparable type; you can compare values with it, but the <, > and <=, >= operation are disable for arrays
	fmt.Println(first == second) // true
	fmt.Println(first != second) // false

	fmt.Println(first < second) // comilation error: invalid operation
}

func emptyArray() {
	var data [10]byte
	fmt.Println(unsafe.Sizeof(data)) // 10

	fmt.Println(data == nil) // compilation error: invalid operation: mismatched type [10]byte and nil
}

func zeroArray() {
	var data [0]int
	fmt.Println(unsafe.Sizeof(data)) // 0
}

func negativeArray() {
	var data [-1]int // compilation error: invalid array lenght
	_ = data
}

func arrayCreation() {
	length1 := 100
	var data1 [length1]int // compilation error: invalid array lenght
	_ = length1
	_ = data1

	const length2 = 100
	var data2 [length2]int
	_ = data2
}

func makeArray() {
	_ = make([10]int, 10) // compilation error: invalid arguments - cannot make type must be slice, map or channel
}

func appendToArray() {
	_ = append([10]int{}, 10) // compilation error: invalid append - argument must be a slice
}

func accessToSliceElement1() {
	data := make([]int, 3)
	fmt.Println(data[4]) // panic: index out of range
}

func accessToSliceElement2() {
	data := make([]int, 3, 6)
	fmt.Println(data[4]) // panic: index out of range
}

func accessToElement3() {
	data := make([]int, 3, 6)
	_ = data[-1] // compilation error: invalid array lenght
}

func accessToNilSlice1() {
	var data []int
	_ = data[0] // panic: index out of range
}

func accessToNilSlice2() {
	var data []int
	data[0] = 10 // panic: index out of range
}

func appendToNilSlice() {
	var data []int
	data = append(data, 10) // make new slice
}

func rangeByNilSlice() {
	var data []int
	for range data { // Выполниться 0 раз, тоесть сразу выйдет так как len(data) == 0
	}
}

func makeZeroSlice() {
	data := make([]int, 0)
	fmt.Println(len(data)) // 0
	fmt.Println(cap(data)) // 0
}

func makeSlice() {
	_ = make([]int, -5)    // compilation error: index must not negative
	_ = make([]int, 10, 5) // compilation error: lenght and capacity swapped (len must be less cap)

	size := -5
	_ = make([]int, size) // panic: len out of range

	size = 5
	_ = make([]int, size*2, size) // panic: cap out of range
}

func sliceMoreThanSize() {
	data := make([]int, 2, 6)
	slice1 := data[1:6] // ok
	_ = slice1
}

func sliceWithIncorrectIndeces() {
	data := make([]int, 2, 6)

	slice2 := data[1:7] // panic: slice bounds of range
	_ = slice2

	slice3 := data[2:1] // compilation error: invalid slice indices (because constants)
	_ = slice3

	left := 2
	right := 1
	slice4 := data[left:right] // panic: slice bounds of range
	_ = slice4
}

func sliceWithNilSlice() {
	var data []int

	slice := data[:]  // ok
	slice = data[0:0] // ok
	slice = data[0:1] // panic: slice bounds of range
	_ = slice
}

func increaseCapacity() {
	data := make([]int, 0, 10)
	data = data[:10:100] // panic: slcie bounds of range
}
