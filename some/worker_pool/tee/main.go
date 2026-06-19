package main

import (
	"fmt"
	"sync"
)

func TeeBlocked[T any](inputCh <-chan T, num int) []<-chan T {
	outputCh := make([]chan T, num)
	for i := range num {
		outputCh[i] = make(chan T)
	}

	go func() {
		// Bocked operation, because if some goroutine non read channel, all goroutine stop,
		// where channel is read
		// Бликорующая операция(гарантия доставки), если одна из горутин не прочитает значание,
		// другие будут остановлены
		for data := range inputCh {
			for _, ch := range outputCh {
				ch <- data
			}
		}

		for _, ch := range outputCh {
			close(ch)
		}
	}()

	outputChFormat := make([]<-chan T, num)
	for i := range num {
		outputChFormat[i] = outputCh[i]
	}
	return outputChFormat
}

func TeeStream[T any](inputCh <-chan T, num int) []<-chan T {
	outputCh := make([]chan T, num)
	for i := range num {
		outputCh[i] = make(chan T)
	}
	outputChClose := make(chan struct{})

	blockedCh := make([]byte, num)
	go func() {
		for data := range inputCh {
			for idx := range blockedCh {
				if blockedCh[idx] == 0 {
					blockedCh[idx] = 1
					go func() {
						defer func() {
							blockedCh[idx] = 0
						}()

						select {
						case <-outputChClose:
							return
						default:
						}

						select {
						case <-outputChClose:
							return
						case outputCh[idx] <- data:
						}
					}()
				}
			}
		}

		close(outputChClose)
		for _, ch := range outputCh {
			close(ch)
		}
	}()

	outputChFormat := make([]<-chan T, num)
	for i := range num {
		outputChFormat[i] = outputCh[i]
	}
	return outputChFormat
}

func main() {
	inputCh := make(chan string)
	go func() {
		defer close(inputCh)
		for i := range 1000 {
			inputCh <- fmt.Sprintf("%d", i)
		}
	}()

	splits := TeeNonBlocked(inputCh, 2)
	var wg sync.WaitGroup
	wg.Add(len(splits))

	for i, ch := range splits {
		go func() {
			defer wg.Done()
			counter := 0
			for data := range ch {
				fmt.Printf("%d channel get: %s\n", i, data)
				counter++
				// 	if i == 1 {
				// 		time.Sleep(1 * time.Millisecond)
				// 	}
			}
			fmt.Printf("%d channel done, get %d data\n", i, counter)
		}()
	}

	wg.Wait()
}
