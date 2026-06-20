package main

import (
	"fmt"
	"time"
)

type Worker struct {
	closeCh     chan struct{}
	closeDoneCh chan struct{}
}

func NewWorker() Worker {
	worker := Worker{
		closeCh:     make(chan struct{}),
		closeDoneCh: make(chan struct{}),
	}

	go func() {
		tiker := time.NewTicker(200 * time.Millisecond)
		defer func() {
			tiker.Stop()
			close(worker.closeDoneCh)
		}()

		for {
			select {
			case <-worker.closeCh:
				return
			default:
			}

			select {
			case <-worker.closeCh:
				return
			case <-tiker.C:
				fmt.Println("get and do some work")
			}
		}
	}()
	return worker
}

func (w *Worker) Shutdown() {
	close(w.closeCh)
	<-w.closeDoneCh
}

func main() {
	w := NewWorker()
	// imitation work
	time.Sleep(1 * time.Second)
	w.Shutdown()
}
