package out

import "github.com/aakashloyar/beats/track/internal/domain"

type AudioVariantRepository interface {
	Save(variant domain.AudioVariant) error
	FindByTrackID(trackID string) ([]domain.AudioVariant, error)
}
