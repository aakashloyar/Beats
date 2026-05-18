package domain

import (
	"time"
	"github.com/aakashloyar/beats/config"
)

type Track struct {
	ID            string
	Title         string
	ArtistID      string
	AlbumID       *string
	CoverImageURL *string
	DurationMS    int64
	Language      []config.Language
	ReleasedAt    *time.Time
	CreatedAt     time.Time
}

type TrackFilter struct {
	Title    *string
	ArtistID *string
	AlbumID  *string
	Language *[]config.Language
	Limit    *string
	Offset   *string
}
