package kafkaconsumer

import (
	"time"
	"github.com/aakashloyar/beats/config"
	"context"
	"encoding/json"
	"log"
	"github.com/aakashloyar/beats/track/internal/application/ports/in/track"
)

type SaveEncodedEventRequest struct {
	ID            string          `json:"upload_id"`
	Title         string          `json:"title"`
	ArtistIDs     []string        `json:"artist_ids"`
	AlbumID       *string         `json:"album_id,omitempty"`
	CoverImageURL *string         `json:"cover_image_url,omitempty"`
	ReleasedAt    *time.Time      `json:"released_at,omitempty"`
	Language      config.Language `json:"language"`
	StreamKey     string          `json:"storage_key"`
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

			var event SaveEncodedEventRequest

			if err := json.Unmarshal(record.Value, &event); err != nil {
				log.Printf("unmarshal error: %v", err)
				continue
			}

			var track in.CreateTrackInput = in.CreateTrackInput{
				ID:            event.ID,
				Title:         event.Title,
				ArtistIDs:     event.ArtistIDs,
				AlbumID:       event.AlbumID,
				CoverImageURL: event.CoverImageURL,
				ReleasedAt:    event.ReleasedAt,
			}

			if err := c.createTrackService.Execute(ctx, track); err != nil {
				log.Printf("create Track failed: %v", err)
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