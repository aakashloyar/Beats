package track

import (
	"github.com/aakashloyar/beats/track/internal/application/ports/out"
    "github.com/aakashloyar/beats/track/internal/domain"
    "context"
)

type GetTrackService struct {
	trackRepo    out.TrackRepository
}

func (s GetTrackService) Execute(ctx context.Context, trackId string) (*domain.Track, error) {
	x,err:=s.trackRepo.FindById(trackId)
	if err!=nil {
		return nil,err
	}
	return x,nil
}
