package service

import (
	"context"
	"errors"
	"github.com/aakashloyar/beats/track/internal/application/ports/in/artist"
	"github.com/aakashloyar/beats/track/internal/application/ports/out"
	"github.com/aakashloyar/beats/track/internal/domain"
)

type CreateArtistService struct {
	artistRepo out.ArtistRepository
	idGen      domain.IDGenerator
	clock      domain.Clock
}

func NewCreateTrackService(artistRepo out.ArtistRepository, idGen domain.IDGenerator, clock domain.Clock) in.CreateArtistService {
	return &CreateArtistService{
		artistRepo: artistRepo,
		idGen:      idGen,
		clock:      clock,
	}
}

func (s *CreateArtistService) Execute(ctx context.Context, input in.CreateArtistInput) (in.CreateArtistOutput, error) {
	if input.Name == "" {
		return in.CreateArtistOutput{}, errors.New("Name is required")
	}
	artist := domain.Artist{
		ID:              s.idGen.NewID(),
		Name:            input.Name,
		Bio:             input.Bio,
		ProfileImageURL: input.ProfileImageURL,
		CreatedAt:       s.clock.Now(),
	}
	if err := s.artistRepo.Save(artist); err != nil {
		return in.CreateArtistOutput{}, err
	}
	return in.CreateArtistOutput{ArtistID: artist.ID}, nil
}
