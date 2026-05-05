package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// Сделать кастомную waitGroup
// Подсказка: структура состоит из атомика, done канала и мьютекса
// Add(), Done(), Wait()

type cuWaitGruop struct {
	size atomic.Int64
	ch   chan struct{}
	mu   sync.Mutex
}

func New() *cuWaitGruop {
	return &cuWaitGruop{
		size: atomic.Int64{},
		ch:   make(chan struct{}),
		mu:   sync.Mutex{},
	}
}

func (w *cuWaitGruop) Add(count int64) {
	newSize := w.size.Add(count)
	if newSize < 0 {
		panic("size wgGroup less zero")
	}
}

func (w *cuWaitGruop) Done() {
	w.mu.Lock()
	defer w.mu.Unlock()

	newSize := w.size.Add(-1)
	if newSize < 0 {
		panic("size wgGroup less zero")
	}
	if newSize == 0 {
		close(w.ch)
		return
	}
}

func (w *cuWaitGruop) Wait() {
	<-w.ch
}

func main() {
	wg := New()

	for i := 0; i < 10; i++ {
		wg.Add(1)

		go func(d int) {
			defer wg.Done()
			fmt.Printf("Hi there from goroutine №%d\n", d)
		}(i)
	}

	wg.Wait()

	fmt.Println("success")
}
