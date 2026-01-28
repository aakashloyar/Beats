package service

import (
	"context"

	"github.com/aakashloyar/beats/track/internal/application/ports/in/artist"
	"github.com/aakashloyar/beats/track/internal/application/ports/out"
)

type GetArtistService struct {
	artistRepo out.ArtistRepository
}

func NewGetArtistService(artistRepo out.ArtistRepository) in.GetArtistService {
	return &GetArtistService{
		artistRepo: artistRepo,
	}
}

func (s *GetArtistService) Execute(ctx context.Context, input in.GetArtistInput) (in.GetArtistOutput, error) {
	x, err := s.artistRepo.FindByID(input.ArtistID)
	if err != nil {
		return in.GetArtistOutput{}, nil
	}
	artist := in.GetArtistOutput{
		ID:              x.ID,
		Name:            x.Name,
		Bio:             x.Bio,
		ProfileImageURL: x.ProfileImageURL,
		CreatedAt:       x.CreatedAt,
	}
	return artist, nil
}
