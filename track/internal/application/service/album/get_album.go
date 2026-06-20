package service

import (
	"context"
	"github.com/aakashloyar/beats/track/internal/application/ports/in/album"
	"github.com/aakashloyar/beats/track/internal/application/ports/out"
)

type GetAlbumService struct {
	albumRepo out.AlbumRepository
}

func NewGetAlbumService(albumRepo out.AlbumRepository) in.GetAlbumService {
	return &GetAlbumService{
		albumRepo: albumRepo,
	}
}

func (s *GetAlbumService) Execute(ctx context.Context, input in.GetAlbumInput) (in.GetAlbumOutput, error) {
	x, err := s.albumRepo.FindByID(input.AlbumID)
	if err != nil {
		return in.GetAlbumOutput{}, nil
	}
	album := in.GetAlbumOutput{
		ID:              x.ID,
		Title:           x.Title,
		CoverImageURL:   x.CoverImageURL,
		ReleaseDate:     x.ReleaseDate,
		CreatedAt:       x.CreatedAt,
	}
	return album, nil
}
