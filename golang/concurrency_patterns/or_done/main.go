package main

import (
	"fmt"
	"time"
)

func orDone[T any](in <-chan T, done <-chan struct{}) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			default:
			}

			select {
			case value, opened := <-in:
				if !opened {
					return
				}

				out <- value
			case <-done:
				return
			}
		}
	}()
	return out
}

func main() {
	ch := make(chan string)
	chDone := make(chan struct{})
	go func() {
		defer close(ch)
		for i := range 1000 {
			ch <- fmt.Sprintf("%d", i)
		}
	}()

	go func() {
		time.Sleep(1 * time.Millisecond)
		close(chDone)
	}()

	for data := range orDone(ch, chDone) {
		fmt.Println(data)
	}
}
