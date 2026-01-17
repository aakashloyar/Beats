package domain

import "time"

type Track struct {
	ID            string
	Title         string
	ArtistID      string
	AlbumID       *string
	CoverImageURL *string
	DurationMs    int64
	Language      Language
	ReleaseDate   *time.Time
	CreatedAt     time.Time
}
