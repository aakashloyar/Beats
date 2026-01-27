package in

import (
	"github.com/aakashloyar/beats/track/internal/domain"
	"context"
)

type GetTrackService interface {
	Execute(ctx context.Context, trackId string) (domain.Track, error)
}