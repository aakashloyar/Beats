package service

import (
	"context"
	"errors"
	"github.com/aakashloyar/beats/track/internal/application/ports/in/track"
	"github.com/aakashloyar/beats/track/internal/application/ports/out"
	"github.com/aakashloyar/beats/track/internal/domain"
)

type CreateTrackService struct {
	trackRepo out.TrackRepository
	idGen     domain.IDGenerator
	clock     domain.Clock
}

func NewCreateTrackService(trackRepo out.TrackRepository, idGen domain.IDGenerator, clock domain.Clock) in.CreateTrackService {
	return &CreateTrackService{
		trackRepo: trackRepo,
		clock:     clock,
		idGen:     idGen,
	}
}

func (s *CreateTrackService) Execute(ctx context.Context, input in.CreateTrackInput) error {
	if input.Title == "" {
		return errors.New("Title is required")
	}
	track := domain.Track{
		ID:            input.ID,
		Title:         input.Title,
		ArtistID:      input.ArtistID,
		AlbumID:       input.AlbumID,
		CoverImageURL: input.CoverImageURL,
		DurationMS:    input.DurationMS,
		Language:      input.Language,
		ReleasedAt:    input.ReleasedAt,
		CreatedAt:     s.clock.Now(),
	}
	if err := s.trackRepo.Save(track); err != nil {
		return err
	}
	return nil

}
