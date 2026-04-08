package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Напиши функцию которая:
//
//	параллельно загружает все URL по конфигурации http.Client и context.Background(),
//	использует не более concurrency горутин одновременно,
//	возвращает срез ошибок в том же порядке, что и urls (если urls[i] завершился успешно — errors[i] == nil, если был ошибка — errors[i] != nil).
//
// Условия:
//
//	Не используй WaitGroup напрямую, сделай через каналы или context + Mutex/Channels.
//	Сигнализация об ошибке — через error, не panic.
type task struct {
	idx int
	err error
}

func workerWithTimeOut(ctx context.Context, client *http.Client, url string) (err error) {
	defer func() {
		if recover() != nil {
			err = fmt.Errorf("PANIC GET url: %s", url)
		}
	}()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Printf("GET status code: %d\n", resp.StatusCode)
	return nil
}

func parallelDownload(urls []string, concurrency int) []error {
	if len(urls) == 0 {
		return nil
	}
	client := http.DefaultClient
	ctx := context.Background()
	errorUrls := make([]error, len(urls))

	chWorkerAnswer := make(chan task, len(urls))
	limitG := make(chan struct{}, concurrency)
	for idx, url := range urls {
		limitG <- struct{}{}
		go func(i int, u string) {
			taskCtx, cancel := context.WithTimeout(ctx, time.Duration(10*time.Second))
			defer cancel()

			chWorkerAnswer <- task{
				idx: i,
				err: workerWithTimeOut(taskCtx, client, u),
			}
			_ = <-limitG
		}(idx, url)
	}
	for range urls {
		t := <-chWorkerAnswer
		errorUrls[t.idx] = t.err
	}

	return errorUrls
}

func main() {
	urls := []string{"https://google.com", "asdf", "https://yandex.ru/", "https://google.com", "https://google.com", "https://google.com", "https://google.com", "https://google.com", "https://google.com", "https://google.com", "https://google.com", "https://google.com", "https://google.com", "https://google.com", "adf", "https://google.com", "https://google.com"}
	fmt.Println(parallelDownload(urls, 4))
}
