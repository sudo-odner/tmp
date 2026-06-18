package main

import (
	"fmt"
	"sync"
	"time"
)

// Task:
// - Есть стукрура WorkerPool
// - Пулл иницииируется с фиксированным колличеством воркеров
// - Канал для задач должен быть небуферизованным
// - Пул должен уметь безопасно принимать задачи.
//   Если все воркеры заняты, и новая задача не может быть принята в течение 500 миллисекунд,
//   пул должен отказаться от неё (не блокируя всю программу намертво)
//   и выдать ошибку/сообщение в консоль.

type WorkerPool struct {
	workCount int
	workChan  chan func()

	closeCh     chan struct{}
	closeDoneCh chan struct{}
}

func NewWorkerPool(workerCount int) (*WorkerPool, error) {
	if workerCount <= 0 {
		return nil, fmt.Errorf("worker count muts be upper zero")
	}

	var wg sync.WaitGroup

	workChan := make(chan func())

	closeCh := make(chan struct{})
	closeDoneCh := make(chan struct{})

	wg.Add(workerCount)
	for range workerCount {
		go func() {
			defer wg.Done()
			for {
				select {
				case <-closeCh:
					return
				default:
				}

				select {
				case <-closeCh:
					return
				case task := <-workChan:
					task()
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(workChan)
		close(closeDoneCh)
	}()

	return &WorkerPool{
		workCount: workerCount,
		workChan:  workChan,

		closeCh:     closeCh,
		closeDoneCh: closeDoneCh,
	}, nil
}

func (wp *WorkerPool) AddTask(task func()) error {
	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()

	select {
	case <-wp.closeCh:
		return fmt.Errorf("worker poll is shutdown, task not started")
	default:
	}

	select {
	case wp.workChan <- task:
		return nil
	case <-wp.closeCh:
		return fmt.Errorf("worker poll is shutdown, task not started")
	case <-t.C:
		return fmt.Errorf("timeout start task 500ms")
	}
}

func (wp *WorkerPool) Shutdown() {
	close(wp.closeCh)
	<-wp.closeDoneCh
}

func main() {
	// Создаю worker pool
	wPool, err := NewWorkerPool(3)
	if err != nil {
		panic(err)
	}

	// Иметриую работу и добавляю задачи в worker pool
	go func() {
		for i := range 100 {
			if err := wPool.AddTask(func() {
				time.Sleep(1 * time.Second)
				fmt.Printf("Task number %d done\n", i)
			}); err != nil {
				fmt.Println("Falied add task:", err.Error())
			}
		}
	}()

	time.Sleep(5 * time.Second)
	wPool.Shutdown()
}
