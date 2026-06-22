package main

import (
	"fmt"
	"math/rand"
)

type Future[T any] struct {
	ch chan T
}

func NewFuture[T any](fn func() T) Future[T] {
	out := make(chan T)
	go func() {
		defer close(out)
		out <- fn()
	}()
	return Future[T]{
		ch: out,
	}
}

func (f *Future[T]) Get() T {
	return <-f.ch
}

func main() {
	asyncJob := func() int {
		return rand.Intn(10)
	}

	future := NewFuture(asyncJob)
	fmt.Printf("get answer on asyncJob: %d\n", future.Get()) // blocking operation
}
