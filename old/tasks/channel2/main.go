package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

// Имеется функция котторая работает неопределенно долго (до 100 секунд)
func randomTimeWork() {
	time.Sleep(time.Duration(rand.Intn(100)) * time.Second)
}

// Написать функцию обертку для этой функции, которая будет перерывать выполнение, если
// фцнкция работает больше 3 секунд, и возврощять ошибку

func predictableTimeWork(ctx context.Context) error {
	ch := make(chan struct{})

	go func() {
		randomTimeWork()
		close(ch)
	}()

	select {
	case <-ch:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("context is done")
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := predictableTimeWork(ctx); err != nil {
		fmt.Println("ERROR: context time is done, work not completed")
	}
	fmt.Println("INFO: work is done")
}
