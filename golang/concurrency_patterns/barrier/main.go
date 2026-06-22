package main

import (
	"fmt"
	"sync"
)

type Barrier struct {
	mutex sync.Mutex
	count int
	size  int

	beforeCh chan struct{}
	afterCh  chan struct{}
}

func NewBarrier(size int) *Barrier {
	return &Barrier{
		size:     size,
		beforeCh: make(chan struct{}, size),
		afterCh:  make(chan struct{}, size),
	}
}

func (b *Barrier) Before() {
	b.mutex.Lock()
	b.count++
	if b.count == b.size {
		for range b.size {
			b.beforeCh <- struct{}{}
		}
	}
	b.mutex.Unlock()
	<-b.beforeCh
}

func (b *Barrier) After() {
	b.mutex.Lock()
	b.count--
	if b.count == 0 {
		for range b.size {
			b.afterCh <- struct{}{}
		}
	}
	b.mutex.Unlock()
	<-b.afterCh
}

func main() {
	var wg sync.WaitGroup
	wg.Add(3)

	bootstrap := func() {
		fmt.Println("bootstrap")
	}

	work := func() {
		fmt.Println("work")
	}

	count := 3
	barrier := NewBarrier(count)
	for range count {
		go func() {
			for range count {
				defer wg.Done()
				barrier.Before()
				bootstrap()

				barrier.After()
				work()
			}
		}()
	}

	wg.Wait()
}
