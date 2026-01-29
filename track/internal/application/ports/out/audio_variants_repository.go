package out

import "github.com/aakashloyar/beats/track/internal/domain"

type AudioVariantRepository interface {
	Save(variant domain.AudioVariant) error
}
