package main

import (
	"context"
	"sync"
	"time"
)

/*
Контекст:
Представь, что у нас есть высоконагруженная система, которая генерирует события (например, логи или метрики).
Мы не хотим отправлять каждое событие в базу данных по отдельности, так как это слишком дорого.
Нам нужно реализовать компонент, который собирает события в пачки (batches) и сохраняет их.
Твоя задача:
Реализовать структуру Batcher, которая:
1) Принимает строки (события) через метод Add.
2) "Сбрасывает" (flush) батч, если:
- Набралось maxBatchSize элементов.
- Прошло flushInterval времени с момента последнего сброса (таймаут).
3) Должна быть потокобезопасной.
4)Должна корректно завершать работу (graceful shutdown) при отмене контекста, сбрасывая оставшиеся в памяти данные.
*/

type Batcher struct {
	dataChan     chan string
	capBatchSize int
	timeLimit    time.Duration
	saveFn       func([]string)
}

func NewBatcher(maxBatchSize int, flushInterval time.Duration, saveFn func([]string)) *Batcher {
	return &Batcher{
		dataChan:     make(chan string, maxBatchSize),
		capBatchSize: maxBatchSize,
		timeLimit:    flushInterval,
		saveFn:       saveFn,
	}
}

func (b *Batcher) Add(event string) {
	b.dataChan <- event
}

func (b *Batcher) Run(ctx context.Context) {
	dataBatch := make([]string, 0, b.capBatchSize)
	ticker := time.NewTicker(b.timeLimit)
	defer ticker.Stop()

	flush := func() {
		if len(dataBatch) > 0 {
			copyDataBatch := make([]string, 0, len(dataBatch))
			copy(copyDataBatch, dataBatch)
			b.saveFn(copyDataBatch)
			dataBatch = dataBatch[:0]
		}
	}

	for {
		select {
		case data := <-b.dataChan:
			dataBatch = append(dataBatch, data)
			if len(dataBatch) >= b.capBatchSize {
				flush()
				ticker.Reset(b.timeLimit)
			}
		case <-ticker.C:
			flush()
		case <-ctx.Done():
			flush()
			return
		}
	}

}

func main() {
	// Это пример использования, писать его не обязательно,
	// но можешь использовать для самопроверки
	save := func(batch []string) {
		println("Saved batch of size:", len(batch))
		println(batch)
	}

	b := NewBatcher(5, 2*time.Second, save)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	wg.Go(func() { b.Run(ctx) })

	b.Add("event 1")
	b.Add("event 2")
	b.Add("event 3")
	b.Add("event 4")
	b.Add("event 5")
	b.Add("event 6")
	b.Add("event 7")
	b.Add("event 8")
	b.Add("event 9")
	b.Add("event 10")
	b.Add("event 11")
	// и так далее
	wg.Wait()
}
