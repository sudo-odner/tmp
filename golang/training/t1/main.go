package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func ProcessTasks(ctx context.Context, tasks []func() error, poolSize int) []error {
	var wg sync.WaitGroup
	// Test

	ch := make(chan struct{}, poolSize) // Канал колличеста активных вокреров
	outErr := make([]error, len(tasks)) // Канал/список ответов задач

	func() {
		for i, t := range tasks {
			select {
			case ch <- struct{}{}:
				wg.Add(1)
				go func() {
					defer wg.Done()
					defer func() { _ = <-ch }()
					outErr[i] = t()
				}()
			case <-ctx.Done():
				for j := i; j < len(tasks); j++ {
					outErr[j] = ctx.Err()
				}
				return
			}
		}
	}()

	wg.Wait()

	return outErr
}

func main() {
	tasks := make([]func() error, 0, 15)
	for i := range 15 {
		tasks = append(tasks, func() error {
			time.Sleep(time.Duration(time.Millisecond * 1))
			return fmt.Errorf("done %d task", i)
		})
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Millisecond*2))
	defer cancel()

	output := ProcessTasks(ctx, tasks, 5)
	for _, o := range output {
		fmt.Printf("%s\n", o)
	}
}
