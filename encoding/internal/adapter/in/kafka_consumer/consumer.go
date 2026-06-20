package kafkaconsumer

import (
	"time"
	"context"
	"log"
	"github.com/aakashloyar/beats/encoding/internal/application/ports/in"
	"encoding/json"
	"github.com/aakashloyar/beats/config"
)

type EncodeUploadedEventRequest struct {
	UploadID      string          `json:"upload_id"`
	Title         string          `json:"title"`
	ArtistIDs     []string        `json:"artist_ids"`
	AlbumID       *string         `json:"album_id,omitempty"`
	CoverImageURL *string         `json:"cover_image_url,omitempty"`
	ReleasedAt    *time.Time      `json:"released_at,omitempty"`
	Language      config.Language `json:"language"`
	StorageKey    string          `json:"storage_key"`
}


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

			var event EncodeUploadedEventRequest

			if err := json.Unmarshal(record.Value, &event); err != nil {
				log.Printf("unmarshal error: %v", err)
				continue
			}

			var track = in.UploadCompletedEvent{
				UploadID:      event.UploadID,
				Title:         event.Title,
				ArtistIDs:     event.ArtistIDs,
				AlbumID:       event.AlbumID,
				CoverImageURL: event.CoverImageURL,
				ReleasedAt:    event.ReleasedAt,
				Language:      event.Language,
				StorageKey:    event.StorageKey,
			}

			if err := c.encodingService.Execute(ctx, track); err != nil {
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