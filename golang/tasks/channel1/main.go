package main

import (
	"fmt"
	"time"
)

// Написать 3 функции:
// writer - генирирует числа от 1 до 10
// doubler - умножает числа на 2, иметируя работу (500ms)
// reader - читает и выводит на экран
// ----- Реализуй код -----
func writer() <-chan int {
	ch := make(chan int)

	go func() {
		for data := range 10 {
			ch <- data
		}
		close(ch)
	}()

	return ch
}

func dubbler(chIn <-chan int) <-chan int {
	ch := make(chan int)

	go func() {
		for data := range chIn {
			time.Sleep(500 * time.Millisecond)
			ch <- (data * 2)
		}
		close(ch)
	}()

	return ch
}

func reader(chIn <-chan int) {
	for data := range chIn {
		fmt.Printf("new data = %d\n", data)
	}
}

// ------------------------
func main() {
	reader(dubbler(writer()))
}
