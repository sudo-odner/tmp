package main

import (
	"fmt"
	"runtime"
	"time"
)

// first goroutine done (but sub func(recursion)
// goroutine leak (deadlock) (write in non-read channel))
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
		case data = <-or(channels[1:]...):
		}
		fmt.Println(data)
		out <- data
	}()
	return out
}

// Signal helps runs without blocked goroutine, but if some time tick long time
// Помогет частично убрать deadlock, но если time или сигналы будут долго обрабатываться
// таже булет утечка
func orSignal[T any](channels ...<-chan T) <-chan T {
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}
	doneCh := make(chan T)

	go func() {
		defer close(doneCh)
		select {
		case <-channels[0]:
		case <-orSignal(channels[1:]...):
		}
	}()

	return doneCh
}

func main() {
	// 3 goroutine start (with main g=4)
	fmt.Println("num g:", runtime.NumGoroutine())
	start := time.Now()
	<-orSignal(
		time.After(200*time.Millisecond),
		time.After(230*time.Millisecond),
		time.After(260*time.Millisecond),
		time.After(410*time.Millisecond),
		time.After(450*time.Millisecond),
		time.After(590*time.Millisecond),
		time.After(600*time.Millisecond),
	)
	fmt.Printf("Called after: %s\n", time.Since(start))
	fmt.Println("num g:", runtime.NumGoroutine())
	time.Sleep(200 * time.Millisecond)
	fmt.Println("num g:", runtime.NumGoroutine())
	// select {} // for see all deadlock(leek 2 goroutine)
}
