package in

import (
	"github.com/aakashloyar/beats/track/internal/domain"
	"time"
)
type CreateTrackInput struct {
    Title         string
    ArtistID      string
    AlbumID       *string
    CoverImageURL *string
    DurationMs    int64
    Language      domain.Language
    ReleaseDate   *time.Time
}

type CreateTrackOutput struct {
    TrackID string
}
