package kafkaconsumer

import (
	"context"
	"log"
	"github.com/aakashloyar/beats/encoding/internal/application/ports/in"
	"encoding/json"
)

func (c *Consumer) Start(ctx context.Context) error {
	for {
		fetches := c.kafkaClient.client.PollFetches(ctx)

		if errs := fetches.Errors(); len(errs) > 0 {
			for _, err := range errs {
				log.Printf("fetch error: %v", err)
			}
			continue
		}

		iter := fetches.RecordIter()
		for !iter.Done() {
			record := iter.Next()

			var event in.UploadCompletedEvent

			if err := json.Unmarshal(record.Value, &event); err != nil {
				log.Printf("unmarshal error: %v", err)
				continue
			}

			if err := c.encodingService.Execute(ctx, event); err != nil {
				log.Printf("encoding failed: %v", err)
				continue
			}

			if err := c.kafkaClient.client.CommitRecords(ctx, record); err != nil {
				log.Printf("commit failed: %v", err)
			}
		}
	}
}

func (c *Consumer) Close() {
	c.kafkaClient.client.Close()
}