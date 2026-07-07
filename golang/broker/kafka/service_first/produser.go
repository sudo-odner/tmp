package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	writer := &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Topic:    "test",
		Balancer: &kafka.Hash{},
	}
	defer writer.Close()

	fmt.Println("starting first service")
	var wg sync.WaitGroup
	workerCount := 5

	wg.Add(workerCount)
	for i := range workerCount {
		go func(idx int) {
			defer wg.Done()
			fmt.Printf("[worker %d] worker starting\n", idx)
			for range 100 {
				userID := fmt.Sprintf("user_%d", rand.Intn(3))
				amount := rand.Intn(1000)

				payload := fmt.Sprintf("{\"amount\": %d, \"worker\": %d}", amount, idx)

				if err := writer.WriteMessages(context.Background(), kafka.Message{
					Key:   []byte(userID),
					Value: []byte(payload),
				}); err != nil {
					fmt.Printf("[worker %d] error writer in kafka, worker stoped\n", idx)
					return
				}
				fmt.Printf("[worker %d] send transaction for %s: %s\n", idx, userID, payload)
				time.Sleep(time.Duration(rand.Intn(300)) * time.Millisecond)
			}
			fmt.Printf("[worker %d] worker stoped\n", idx)
		}(i)
	}

	wg.Wait()
}
