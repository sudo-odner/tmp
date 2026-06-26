package main

import (
	"fmt"
	"net/http"
	"sync"
)

// Написать программу, которая параллельно отправляет НТТР-запросы к двум URL и выводит статус-коды ответов.
// Требования: запросы действительно параллельны (горутины + WaitGroup или каналы), ошибки корректно обработаны.

type responce struct {
	url string

	resp *http.Response
	err  error
}

func main() {
	urls := []string{
		"https://google.com",
		"https://yandexl.ru",
	}
	var wg sync.WaitGroup
	output := make(chan responce)

	wg.Add(len(urls))
	for _, url := range urls {
		go func() {
			defer wg.Done()
			out := responce{
				url: url,
			}
			out.resp, out.err = http.Get(url)

			output <- out
		}()
	}
	go func() {
		wg.Wait()
		close(output)
	}()

	for answer := range output {
		fmt.Println(answer)
		if answer.err != nil {
			fmt.Println("failed get url:", answer.url, "-", answer.err)
		} else {
			fmt.Println("responce url:", answer.url, ", status:", answer.resp.Status)
			answer.resp.Body.Close()
		}
	}
}
