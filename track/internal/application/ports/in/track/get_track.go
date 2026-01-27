package in

import (
	"github.com/aakashloyar/beats/track/internal/domain"
	"context"
)

type GetTrackService interface {
	Execute(ctx context.Context, trackID string) (domain.Track, error)
}