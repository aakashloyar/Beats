package in

import (
	"context"
	"github.com/aakashloyar/beats/config"
	"time"
)

type CreateTrackInput struct {
	ID            string
	Title         string
	ArtistIDs     []string
	AlbumID       *string
	CoverImageURL *string
	DurationMS    int64
	Language      []config.Language
	ReleasedAt    *time.Time
}

type CreateTrackService interface {
	Execute(ctx context.Context, input CreateTrackInput) error
}
