package main

import (
	"fmt"
	"time"
)

func or[T any](channels ...<-chan T) <-chan T {
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}

	out := make(chan T)
	go func() {
		defer close(out)
		var data T
		select {
		case data = <-channels[0]:
		case data = <-channels[1]:
		case data = <-or(channels[2:]...):
		}
		fmt.Println(data)
		out <- data
	}()
	return out
}

func main() {
	// 3 goroutine start (with main g=4)
	<-or(
		time.After(200*time.Millisecond),
		time.After(230*time.Millisecond),
		time.After(260*time.Millisecond),
		time.After(410*time.Millisecond),
		time.After(450*time.Millisecond),
		time.After(590*time.Millisecond),
		time.After(600*time.Millisecond),
	)
	// first goroutine done (but sub func(recursion)
	// goroutine leak (deadlock) (write in non-read channel))
	time.Sleep(1 * time.Second)

	select {} // can see all deadlock(leek 2 goroutine)
}
