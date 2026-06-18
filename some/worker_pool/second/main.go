package main

import (
	"fmt"
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
}

func NewWorkerPool(workerCount int) (*WorkerPool, error) {
	if workerCount <= 0 {
		return nil, fmt.Errorf("worker count muts be upper zero")
	}
	workChan := make(chan func())

	for range workerCount {
		go func() {
			for task := range workChan {
				task()
			}
		}()
	}

	return &WorkerPool{
		workCount: workerCount,
		workChan:  workChan,
	}, nil
}

func (wp *WorkerPool) AddTask(task func()) error {
	t := time.NewTicker(500 * time.Microsecond)
	defer t.Stop()
	select {
	case wp.workChan <- task:
		return nil
	case <-t.C:
		return fmt.Errorf("timeout start task 500ms")
	}
}

func main() {
	wPool, err := NewWorkerPool(3)
	if err != nil {
		panic(err)
	}
	for i := range 100 {
		if err := wPool.AddTask(func() {
			time.Sleep(1 * time.Second)
			fmt.Printf("Task number %d done\n", i)
		}); err != nil {
			fmt.Println("Falied add task", err.Error())
		}
	}

	time.Sleep(5 * time.Second)
}
