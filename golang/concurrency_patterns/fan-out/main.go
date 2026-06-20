package main

import (
	"fmt"
	"sync"
)

func FanOutBlocked[T any](inputCh <-chan T, num int) []<-chan T {
	outputCh := make([]chan T, num)
	for i := range outputCh {
		outputCh[i] = make(chan T)
	}

	go func() {
		idx := 0
		for data := range inputCh {
			outputCh[idx] <- data // with blocked operation (but with balance)
			idx = (idx + 1) % num
		}

		for _, ch := range outputCh {
			close(ch)
		}
	}()

	outputChFormat := make([]<-chan T, num)
	for i := 0; i < num; i++ {
		outputChFormat[i] = outputCh[i]
	}
	return outputChFormat
}

func FanOutNonBlocked[T any](inputCh <-chan T, num int) []<-chan T {
	outputCh := make([]chan T, num)
	for i := range num {
		outputCh[i] = make(chan T)
	}

	for _, ch := range outputCh {
		go func() {
			defer close(ch)
			for data := range inputCh {
				ch <- data // with out blocked, but no balance
			}
		}()
	}

	outputChFormat := make([]<-chan T, num)
	for i := range num {
		outputChFormat[i] = outputCh[i]
	}
	return outputChFormat
}

func main() {
	mainCh := make(chan string)
	go func() {
		defer close(mainCh)
		for i := range 200 {
			mainCh <- fmt.Sprintf("%d", i)
		}
	}()

	splits := FanOutNonBlocked(mainCh, 2)
	var wg sync.WaitGroup

	wg.Add(2)
	for i, splitCh := range splits {
		go func() {
			defer wg.Done()
			counter := 0
			for data := range splitCh {
				fmt.Printf("%d channel get: %s\n", i, data)
				counter++
			}
			fmt.Printf("%d channel read: %d data\n", i, counter)
		}()
	}
	wg.Wait()
}
