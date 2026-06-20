package main

import "fmt"

func process[T any](in <-chan T) <-chan T {
	out := make(chan T)

	go func() {
		defer close(out)
		for {
			select {
			case <-in:
				return
			default:
				// processing
			}
		}
	}()

	return out
}

func main() {
	chClose := make(chan struct{})
	chCloseDone := process(chClose)

	close(chClose)
	<-chCloseDone

	fmt.Println("terminated")
}
