package main

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type ErrorGroup struct {
	err  error
	wg   sync.WaitGroup
	once sync.Once

	done chan struct{}
}

func NewErrorGroup() (*ErrorGroup, <-chan struct{}) {
	doneCh := make(chan struct{})
	return &ErrorGroup{
		done: doneCh,
	}, doneCh
}

func (eg *ErrorGroup) Go(fn func() error) {
	eg.wg.Add(1)
	go func() {
		defer eg.wg.Done()
		select {
		case <-eg.done:
			return
		default:
			if err := fn(); err != nil {
				eg.once.Do(func() {
					eg.err = err
					close(eg.done)
				})
			}
		}
	}()
}

func (eg *ErrorGroup) Wait() error {
	eg.wg.Wait()
	return eg.err
}

func main() {
	group, groupDone := NewErrorGroup()
	for range 5 {
		group.Go(func() error {
			timeout := time.Second * time.Duration(rand.Intn(10))
			timer := time.NewTimer(timeout)

			select {
			case <-groupDone:
				return errors.New("error done")
			case <-timer.C:
				return errors.New("error timer")
			}
		})
	}

	if err := group.Wait(); err != nil {
		fmt.Println(err.Error())
	}
}
