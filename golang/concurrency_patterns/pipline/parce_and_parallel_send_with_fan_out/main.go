package main

import (
	"fmt"
	"sync"
)

// Pipline
//                        	       / -> split1 -> send   |   - \
// inputChannel -> parce -> split -                      |    -> outputChannel
//                                 \ -> split2 -> send   |   - /
//

func Split[T any](in <-chan T, num int) []<-chan T {
	out := make([]chan T, num)
	for i := range num {
		out[i] = make(chan T)
	}

	go func() {
		idx := 0
		for value := range in {
			out[idx] <- value
			idx = (idx + 1) % num
		}

		for _, ch := range out {
			close(ch)
		}
	}()

	outF := make([]<-chan T, num)
	for i := range out {
		outF[i] = out[i]
	}

	return outF
}

func Parce(in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for value := range in {
			out <- fmt.Sprintf("parce - %s", value)
		}
	}()
	return out
}

func Send(in <-chan string, num int) <-chan string {
	out := make(chan string)

	var wg sync.WaitGroup

	wg.Add(num)
	for i, ch := range Split(in, num) {
		go func() {
			defer wg.Done()
			for value := range ch {
				out <- fmt.Sprintf("split%d send - %s", i, value)
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
		for i := range 100 {
			ch <- fmt.Sprintf("%d", i)
		}
	}()

	for value := range Send(Parce(ch), 2) {
		fmt.Println(value)
	}
}
