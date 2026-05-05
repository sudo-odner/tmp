package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func processData(val int) int {
	time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
	return val * 2
}

func main() {
	in := make(chan int)
	out := make(chan int)

	go func() {
		for i := range 100 {
			in <- i
		}
		close(in)
	}()

	now := time.Now()
	processParallel(in, out, 5)

	for val := range out {
		fmt.Println(val)
	}
	fmt.Println(time.Since(now))
}

// Операция не должна выполняться больше 5 секунд

func processParallel(in <-chan int, out chan<- int, numWorkers int) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	var wg sync.WaitGroup
	wg.Add(numWorkers)

	for range numWorkers {
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case data, ok := <-in:
					if !ok {
						return
					}
					ch := make(chan int)
					go func() { ch <- processData(data) }()
					select {
					case <-ctx.Done():
						return
					case newData := <-ch:
						out <- newData
					}
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		cancel()
		close(out)
	}()
}
