package in

import (
	"context"
	"github.com/aakashloyar/beats/config"
	"time"
)

type ListTracksInput struct {
	Title     string
	ArtistIDs []string
	AlbumID   string
	Limit     string
	Offset    string
}

type ListTracksOutput struct {
	ID            string
	Title         string
	ArtistIDs     []string
	AlbumID       *string
	CoverImageURL *string
	DurationMS    int64
	Language      []config.Language
	ReleasedAt    *time.Time
	CreatedAt     time.Time
}

type ListTracksService interface {
	Execute(ctx context.Context, input ListTracksInput) ([]ListTracksOutput, error)
}
