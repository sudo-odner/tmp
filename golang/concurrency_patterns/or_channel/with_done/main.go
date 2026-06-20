package main

import (
	"fmt"
	"runtime"
	"time"
)

func OrSignal[T any](channels ...<-chan T) <-chan T {
	out := make(chan T)
	done := make(chan struct{})
	go func() {
		defer func() {
			close(done)
			close(out)
		}()
		for value := range orSignal(done, channels...) {
			out <- value
		}
	}()
	return out
}

func orSignal[T any](done <-chan struct{}, channels ...<-chan T) <-chan T {
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}

	out := make(chan T)

	go func() {
		defer close(out)
		select {
		case <-channels[0]:
		case <-channels[2]:
		case <-orSignal(done, channels[2:]...):
		case <-done:
		}
	}()

	return out
}

func main() {
	fmt.Println("num g:", runtime.NumGoroutine())
	start := time.Now()
	<-OrSignal(
		time.After(100*time.Millisecond),
		time.After(550*time.Millisecond),
		time.After(999*time.Millisecond),
		time.After(1*time.Second),
		time.After(1*time.Hour),
		time.After(12*time.Hour),
		time.After(24*time.Hour),
	)
	fmt.Printf("Called after: %s\n", time.Since(start))
	fmt.Println("num g:", runtime.NumGoroutine())
	time.Sleep(200 * time.Millisecond)
	fmt.Println("num g:", runtime.NumGoroutine())
	// select {} // for check deadlocks
}
