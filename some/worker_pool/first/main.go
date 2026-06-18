package main

import (
	"fmt"
	"sync"
	"time"
)

// Task 1: Есть код, нужно реализовать простой worker pool
// type Job func()
//
// func main() {
// 	numWorkers := 3
// 	numJobs := 10
//
// 	jobs := make(chan Job, numJobs)
//
// 	// 1. ЗАДАЧА: Запустить воркеров в цикле.
// 	// Каждый воркер должен читать задачи из канала jobs, пока канал не закроется.
//
// 	// 2. ЗАДАЧА: Отправить 10 задач в канал.
// 	// Каждая задача — это простая функция, например: func() { fmt.Println("делаю работу...") }
//
// 	// 3. ЗАДАЧА: Правильно закрыть канал и дождаться выполнения всех воркеров.
// }

// Answer:
// Примичание данный код валиден на вресиях 1.22 и выше. Так как for делает снимки при каждом цыкле,
// В предыдущих версиях for мог отработать быстрее чем горутины и все последующие goroutine
// прочитали бы i == numWorkers
type Job func()

func main() {
	numWorkers := 3
	numJobs := 10

	jobs := make(chan Job, numJobs)
	var wg sync.WaitGroup

	// Запускаю workers, где с помощю range проожу buff cahn (range читает канал пока тот открыт)
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func() {
			defer wg.Done()
			for job := range jobs {
				fmt.Printf("Goroutine number %d start job\n", i)
				job()
			}
		}()
	}

	// Имитирую работу и нагрузку на workers pool
	for i := 0; i < 10; i++ {
		jobs <- func() { time.Sleep(1 * time.Second) }
	}

	// Закрываю канал и жду пока все workers доработают
	close(jobs)
	wg.Wait()
}
