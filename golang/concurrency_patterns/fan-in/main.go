package main

import (
	"fmt"
	"sync"
)

func FanIn[T any](channels ...<-chan T) <-chan T {
	out := make(chan T)
	var wg sync.WaitGroup
	wg.Add(len(channels))

	for _, ch := range channels {
		go func(c <-chan T) {
			defer wg.Done()
			for c := range ch {
				out <- c
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		defer close(ch1)
		for i := range 100 {
			ch1 <- fmt.Sprintf("ch1: %d", i)
		}
	}()

	go func() {
		defer close(ch2)
		for i := range 100 {
			ch2 <- fmt.Sprintf("ch2: %d", i)
		}
	}()

	for data := range FanIn(ch1, ch2) {
		fmt.Println(data)
	}
}
