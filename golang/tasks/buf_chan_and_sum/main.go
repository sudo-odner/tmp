package main

import (
	"fmt"
	"sync"
)

func main() {
	c := make(chan int, 1000)
	var wg sync.WaitGroup

	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			foo(c)
		}()
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	sum := 0
	for r := range c {
		sum += r
	}

	fmt.Println(sum)
}

func foo(c chan int) {
	r := 10
	for i := 0; i < r; i++ {
		c <- r
	}
}
