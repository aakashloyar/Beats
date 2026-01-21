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
    clock     out.Clock
    idGen     out.IDGenerator
}


func (s CreateTrackService) Execute(ctx context.Context, input *in.CreateTrackInput) (*in.CreateTrackOutput,error){
    if input.Title == "" {
        return nil, errors.New("Title is required")
    }
    track:=&domain.Track{
        ID: s.idGen.NewID(),
        Title: input.Title,
        ArtistID: input.ArtistID,
        AlbumID: input.AlbumID,
        CoverImageURL: input.CoverImageURL,
        DurationMs: input.DurationMs,
        Language: input.Language,
        ReleaseDate: input.ReleaseDate,
        CreatedAt: s.clock.Now(),
    }
    if err := s.trackRepo.Save(track) ; err != nil {
        return nil, err
    }
    return &in.CreateTrackOutput{TrackID: track.ID}, nil

}
