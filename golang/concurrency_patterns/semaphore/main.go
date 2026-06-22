package main

import (
	"fmt"
	"sync"
	"time"
)

type Semaphore struct {
	tickets chan struct{}
}

func NewSemaphore(num int) *Semaphore {
	return &Semaphore{
		tickets: make(chan struct{}, num),
	}
}

func (s *Semaphore) Acquire() {
	s.tickets <- struct{}{}
}

func (s *Semaphore) Realase() {
	<-s.tickets
}

func main() {
	var wg sync.WaitGroup

	semaphore := NewSemaphore(5)
	wg.Add(7)
	for range 7 {
		go func() {
			defer func() {
				wg.Done()
				semaphore.Realase()
			}()
			semaphore.Acquire()

			fmt.Println("working...")
			time.Sleep(2 * time.Second)
			fmt.Printf("exiting...")
		}()
	}

	wg.Wait()
}
