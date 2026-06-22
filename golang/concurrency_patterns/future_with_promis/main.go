package main

import (
	"fmt"
	"time"
)

type Future[T any] struct {
	ch <-chan T
}

func NewFuture[T any](in <-chan T) *Future[T] {
	return &Future[T]{
		ch: in,
	}
}

func (f *Future[T]) Get() T {
	return <-f.ch
}

type Promis[T any] struct {
	ch chan T
}

func NewPromis[T any]() Promis[T] {
	return Promis[T]{
		ch: make(chan T),
	}
}

func (p *Promis[T]) Set(value T) {
	p.ch <- value
	close(p.ch)
}

func (p *Promis[T]) GetFuture() *Future[T] {
	return NewFuture(p.ch)
}

func main() {
	// Созадею обещание что когда выполню то уведомлю
	promis := NewPromis[string]()

	go func() {
		time.Sleep(2 * time.Second)
		promis.Set("i am done")
	}()

	// Получаю объект будущего который когда то отработает в фоне и когда понадобится попрошу
	future := promis.GetFuture()
	// Жду future
	fmt.Printf("get future: %s\n", future.Get())
}
