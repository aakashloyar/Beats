package service

import (
	"context"
	"github.com/aakashloyar/beats/track/internal/application/ports/in/track"
	"github.com/aakashloyar/beats/track/internal/application/ports/out"
)

type GetTrackService struct {
	trackRepo out.TrackRepository
}

func NewGetTrackService(trackRepo out.TrackRepository) in.GetTrackService {
	return &GetTrackService{
		trackRepo: trackRepo,
	}
}
func (s *GetTrackService) Execute(ctx context.Context, input in.GetTrackInput) (in.GetTrackOutput, error) {
	x, err := s.trackRepo.FindByID(input.TrackID)
	if err != nil {
		return in.GetTrackOutput{}, err
	}
	track := in.GetTrackOutput{
		ID:            x.ID,
		ArtistID:      x.ArtistID,
		AlbumID:       x.AlbumID,
		CoverImageURL: x.CoverImageURL,
		DurationMS:    x.DurationMS,
		Language:      x.Language,
		ReleaseDate:   x.ReleaseDate,
		CreatedAt:     x.CreatedAt,
	}
	return track, nil
}
