package main

import (
	"fmt"
	"sync"
)

// Pipline
//                       / -> send   |   - \
// inputChannel -> parce -           |      -> outputChannel
//                       \ -> send   |   - /
//

func parce(in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		for value := range in {
			out <- fmt.Sprintf("parce - %s", value)
		}
		close(out)
	}()
	return out
}

func send(in <-chan string, num int) <-chan string {
	out := make(chan string)
	var wg sync.WaitGroup
	wg.Add(num)

	for i := range num {
		go func() {
			defer wg.Done()
			for value := range in {
				out <- fmt.Sprintf("send %d - %s", i, value)
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	ch := make(chan string)
	go func() {
		defer close(ch)
		for i := range 10 {
			ch <- fmt.Sprintf("%d", i)
		}
	}()

	for value := range send(parce(ch), 2) {
		fmt.Println(value)
	}
}
