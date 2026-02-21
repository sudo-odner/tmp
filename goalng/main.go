package main

import (
	"fmt"
)

func main() {
	slices := []int{1, 2, 3}
	slices = append(slices, 3, 4, 5, 6)
	fmt.Println(slices, len(slices), cap(slices))
}
