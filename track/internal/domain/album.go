package domain

import "time"

type Album struct {
	ID            string
	Title         string
	CoverImageURL *string
	ReleaseDate   *time.Time
	CreatedAt     time.Time
}
