package out

import (
	"github.com/aakashloyar/beats/track/internal/domain"
)

type TrackRepository interface {
	Save(track domain.Track) error
	FindByID(trackID string) (domain.Track, error)
	ListTracks(input domain.TrackFilter) ([]domain.Track, error)
	ListAudioVariantsByTrack(trackID string) ([]domain.AudioVariant, error)
}
