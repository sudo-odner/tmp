package main

import (
	"fmt"
	"sync"
)

type Counter struct {
	counter int64
	mut     sync.Mutex
}

func (c *Counter) Inc() {
	c.mut.Lock()
	defer c.mut.Unlock()

	c.counter++
}

func (c *Counter) Dec() {
	c.mut.Lock()
	defer c.mut.Unlock()

	c.counter--
}

func (c *Counter) Get() int64 {
	c.mut.Lock()
	defer c.mut.Unlock()

	return c.counter
}

func (c *Counter) Reset() {
	c.mut.Lock()
	defer c.mut.Unlock()

	c.counter = 0
}

func main() {
	var counter Counter

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				counter.Inc()
			}
		}()
	}

	wg.Wait()
	fmt.Println("Final count:", counter.Get())

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				counter.Dec()
			}
		}()
	}

	wg.Wait()
	fmt.Println("After decrements, counter:", counter.Get())

	counter.Reset()
	fmt.Println("After reset, counter:", counter.Get())
}
