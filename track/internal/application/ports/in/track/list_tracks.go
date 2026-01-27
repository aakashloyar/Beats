package in

import (
	"github.com/aakashloyar/beats/track/internal/domain"
	"context"
)

type ListTracksInput struct {
	Title    string
    ArtistID string
	AlbumID  string
	Limit    string
	Offset   string 
}

type ListTracksService interface {
    Execute(ctx context.Context, input ListTracksInput) ([]domain.Track, error)
}

