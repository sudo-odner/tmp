package main

import "fmt"

func Transformer[T any](in <-chan T, fn func(T) T) <-chan T {
	out := make(chan T)
	go func() {
		for data := range in {
			out <- fn(data)
		}
		close(out)
	}()

	return out
}

func main() {
	ch := make(chan string)
	go func() {
		for i := range 100 {
			ch <- fmt.Sprintf("%d", i)
		}
		close(ch)
	}()

	transform := func(in string) string {
		return fmt.Sprintf("md - %s", in)
	}

	for data := range Transformer(ch, transform) {
		fmt.Println(data)
	}
}
