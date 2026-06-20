package main

import "fmt"

func main() {
	a := make([]int, 3, 5)     // a: len=3, cap=5 [0,0,0,0,0]
	a[0], a[1], a[2] = 1, 2, 3 // a: len=3, cap=5 [1,2,3,0,0]

	b := append(a, 4) // b(*a); len=4 cap=5 [1,2,3,4,0]
	c := append(a, 5) // c(*a); len=4 cap=5 [1,2,3,5,0]

	fmt.Println("a:", a) // [1,2,3,5,0]
	fmt.Println("b:", b) // [1,2,3,5,0]
	fmt.Println("c:", c) // [1,2,3,5,0]

	b[0] = 99
	fmt.Println("after modifying b[0]:")
	fmt.Println("a:", a) // [99,2,3,5,0]
	fmt.Println("b:", b) // [99,2,3,5,0]
	fmt.Println("c:", c) // [99,2,3,5,0]
}
