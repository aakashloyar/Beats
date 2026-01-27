package track

import (
	"github.com/aakashloyar/beats/track/internal/application/ports/out"
	"github.com/aakashloyar/beats/track/internal/application/ports/in/track"
    "github.com/aakashloyar/beats/track/internal/domain"
    "context"
)

type GetTrackService struct {
	trackRepo    out.TrackRepository
}

func NewGetTrackService(trackRepo out.TrackRepository) in.GetTrackService {
	return &GetTrackService{
		trackRepo: trackRepo,
	}
}
func (s *GetTrackService) Execute(ctx context.Context, trackID string) (domain.Track, error) {
	x,err:=s.trackRepo.FindByID(trackID)
	if err!=nil {
		return domain.Track{},err
	}
	return x,nil
}

