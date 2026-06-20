package main

import "fmt"

func Filter[T any](in <-chan T, predicate func(T) bool) <-chan T {
	out := make(chan T)
	go func() {
		for data := range in {
			if predicate(data) {
				out <- data
			}
		}
		close(out)
	}()

	return out
}

func main() {
	ch := make(chan int)
	go func() {
		for i := range 100 {
			ch <- i
		}
		close(ch)
	}()

	filter := func(in int) bool {
		return in%2 == 0
	}

	for data := range Filter(ch, filter) {
		fmt.Println(data)
	}
}
