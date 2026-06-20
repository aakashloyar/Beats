package out

import (
	"context"
	"time"
	"github.com/aakashloyar/beats/config"
)	

type UploadCompletedEvent struct {
	UploadID      string          `json:"upload_id"`
	Title         string          `json:"title"`
	ArtistIDs     []string        `json:"artist_ids"`
	AlbumID       *string         `json:"album_id,omitempty"`
	CoverImageURL *string         `json:"cover_image_url,omitempty"`
	ReleasedAt    *time.Time      `json:"released_at,omitempty"`
	Language      config.Language `json:"language"`
	StorageKey    string          `json:"storage_key"`
}

type EventPublisher interface {
	PublishUploadCompleted(ctx context.Context, event UploadCompletedEvent) error
}