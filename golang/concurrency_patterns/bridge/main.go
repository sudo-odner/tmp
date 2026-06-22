package main

import "fmt"

func Bridge[T any](in chan chan T) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		for channel := range in {
			for value := range channel {
				out <- value
			}
		}
	}()
	return out
}

func main() {
	chch := make(chan chan string)
	go func() {
		defer close(chch)
		ch1 := make(chan string, 5)
		for i := 0; i < 10; i += 2 {
			ch1 <- fmt.Sprintf("ch1: %d", i)
		}
		close(ch1)

		ch2 := make(chan string, 5)
		for i := 1; i < 10; i += 2 {
			ch2 <- fmt.Sprintf("ch2: %d", i)
		}
		close(ch2)

		chch <- ch1
		chch <- ch2
	}()

	for value := range Bridge(chch) {
		fmt.Println(value)
	}
}
