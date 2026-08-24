package main

import (
	"context"
	"fmt"
	"time"
)

// Что произойдет при выполении программы?
//
// Answer: Deadlock, main функция заблокируется навсегда, потому что нет читателя
// у буферизированного канала
func buffChan() {
	tch := make(chan int)
	tch <- 42
	fmt.Println(<-tch)
}

// Что произойдет в данном случае?
//
// Answer: Deadlock(main goroutine), не дожидаясь выполнение другой горутины
// (в теории, если она выполнялась бы быстрее чем 2сек, она могла бы отработать быстрее чем main
// поток, но планировщик Go не детерминирован, мы не можем предугадать в какой последовательности
// будут выполнены горутины)
func buffChanWithGoroutine() {
	tch := make(chan int)
	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println("Goroutine is done")
	}()
	tch <- 10
}

// Перед тобой код, коллега утвержает что здесь есть утечка горути (goroutine leak) или
// частичный deadlock, если context будет отменен. Прав ли он?
//
// Answer: Нет, код написан правильно
func contextWorker(ctx context.Context, jobs <-chan int) {
	for {
		select {
		case <-ctx.Done():
			return
		case job := <-jobs:
			fmt.Println("Get a job:", job)
		}
	}
}

func main() {
}
