package track

import (
	"github.com/aakashloyar/beats/track/internal/application/ports/out"
    "github.com/aakashloyar/beats/track/internal/application/ports/in"
    "github.com/aakashloyar/beats/track/internal/domain"
    "errors"
    "context"
)

type CreateTrackService struct {
    trackRepo out.TrackRepository
    clock     domain.Clock
    idGen     domain.IDGenerator
}

func NewCreateTrackService( trackRepo out.TrackRepository, clock domain.Clock, idGen domain.IDGenerator) in.CreateTrackService {
    return &CreateTrackService{
        trackRepo: trackRepo,
        clock : clock,
        idGen : idGen,
    }
}

func (s *CreateTrackService) Execute(ctx context.Context, input in.CreateTrackInput) (in.CreateTrackOutput,error){
    if input.Title == "" {
        return in.CreateTrackOutput{}, errors.New("Title is required")
    }
    track:=domain.Track{
        ID: s.idGen.NewID(),
        Title: input.Title,
        ArtistID: input.ArtistID,
        AlbumID: input.AlbumID,
        CoverImageURL: input.CoverImageURL,
        DurationMS: input.DurationMS,
        Language: input.Language,
        ReleaseDate: input.ReleaseDate,
        CreatedAt: s.clock.Now(),
    }
    if err := s.trackRepo.Save(track) ; err != nil {
        return in.CreateTrackOutput{}, err
    }
    return in.CreateTrackOutput{TrackID: track.ID}, nil

}
