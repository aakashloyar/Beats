package out
import (
	"github.com/aakashloyar/beats/track/internal/domain"
	"github.com/aakashloyar/beats/track/internal/application/ports/in"
)

type TrackRepository interface {
	Save(track domain.Track) error
	FindById(id string) (domain.Track, error)
	ListTracks(input in.ListTracksInput) ([]domain.Track, error)
}

