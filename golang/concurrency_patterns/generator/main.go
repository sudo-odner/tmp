package main

import (
	"fmt"
)

func GeneratorWithChannel(start, end int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for number := start; number < end; number++ {
			out <- number
		}
	}()
	return out
}

func main() {
	for number := range GeneratorWithChannel(100, 200) {
		fmt.Println(number)
	}
}
