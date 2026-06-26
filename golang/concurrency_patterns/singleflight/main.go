package main

import (
	"fmt"
	"sync"
	"time"
)

type call struct {
	value interface{}
	err   error

	done chan struct{}
}

type SingleFlight struct {
	mu    sync.Mutex
	calls map[string]*call
}

func NewSingleFlight() *SingleFlight {
	return &SingleFlight{
		calls: make(map[string]*call),
	}
}

func (sg *SingleFlight) Do(hash string, fn func() (interface{}, error)) (interface{}, error) {
	sg.mu.Lock()
	if call, ok := sg.calls[hash]; ok {
		sg.mu.Unlock()
		return sg.wait(call)
	}

	newCall := call{
		done: make(chan struct{}),
	}
	sg.calls[hash] = &newCall
	sg.mu.Unlock()
	go func() {
		defer func() {
			sg.mu.Lock()
			close(newCall.done)
			delete(sg.calls, hash)
			sg.mu.Unlock()
		}()
		newCall.value, newCall.err = fn()
	}()

	return sg.wait(&newCall)
}

func (sg *SingleFlight) wait(call *call) (interface{}, error) {
	<-call.done
	return call.value, call.err
}

func main() {
	const inFlightRequests = 10
	var wg sync.WaitGroup
	wg.Add(inFlightRequests)

	singleFlight := NewSingleFlight()
	for i := range inFlightRequests {
		go func() {
			defer wg.Done()
			value, err := singleFlight.Do("insert", func() (interface{}, error) {
				fmt.Println("singleFlight")
				time.Sleep(1 * time.Second)
				return "some", nil
			})

			fmt.Println(i, "=", value, ";", err)
		}()
	}
	wg.Wait()
}
