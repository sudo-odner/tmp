package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	fmt.Println("starting second service")
	var wg sync.WaitGroup

	countReader := 5
	wg.Add(countReader)
	for i := range countReader {
		go func(idx int) {
			defer wg.Done()

			reader := kafka.NewReader(kafka.ReaderConfig{
				Brokers: []string{"localhost:9092"},
				GroupID: "concurrent-processors",
				Topic:   "test",
			})
			defer reader.Close()

			fmt.Printf(
				"[worker %d][time %s] worker init and starting\n",
				idx, time.Now().Format("2006-01-02 15:04:05"),
			)
			ctx := context.Background()
			for {
				msg, err := reader.FetchMessage(ctx)
				if err != nil {
					fmt.Printf(
						"[worker %d][time %s] failed read msg, error: %s\n",
						idx, time.Now().Format("2006-01-02 15:04:05"), err.Error(),
					)
					time.Sleep(1 * time.Second)
					continue
				}

				fmt.Printf(
					"[worker %d][time %s] read msg:\n\tpartition=%d\n\toffset=%d\n\tkey=%s\n\tvalue=%s\n",
					idx, time.Now().Format("2006-01-02 15:04:05"), msg.Partition, msg.Offset, msg.Key, msg.Value,
				)
				err = reader.CommitMessages(ctx, msg)
				if err != nil {
					fmt.Printf(
						"[worker %d][time %s] failed commit messsage: %s\n",
						idx, time.Now().Format("2006-01-02 15:04:05"), err.Error(),
					)
				}
			}
		}(i)
	}
	wg.Wait()
}
