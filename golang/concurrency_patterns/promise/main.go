package main

import (
	"fmt"
	"math/rand"
	"time"
)

type result[T any] struct {
	val T
	err error
}

type Promis[T any] struct {
	ch chan result[T]
}

func NewPromis[T any](asyncJob func() (T, error)) *Promis[T] {
	out := make(chan result[T])
	go func() {
		defer close(out)
		value, err := asyncJob()
		out <- result[T]{
			val: value,
			err: err,
		}
	}()
	return &Promis[T]{
		ch: out,
	}
}

func (p *Promis[T]) Then(succesFn func(T), errorFn func(error)) {
	go func() {
		answer := <-p.ch
		if answer.err != nil {
			errorFn(answer.err)
		} else {
			succesFn(answer.val)
		}
	}()
}

func main() {
	asyncFn := func() (int, error) {
		value := rand.Intn(10)
		if value%2 != 0 {
			return 0, fmt.Errorf("falid gen value, is it: %d", value)
		}
		return value, nil
	}

	pr := NewPromis(asyncFn)
	pr.Then(
		func(value int) {
			fmt.Println("Success function: ", value)
		},
		func(err error) {
			fmt.Println("Error function: ", err)
		},
	)

	time.Sleep(2 * time.Second)

	// return Them or success function or error function
}
