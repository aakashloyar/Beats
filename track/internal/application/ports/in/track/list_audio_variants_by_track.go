package in

import (
	"context"
	"github.com/aakashloyar/beats/track/internal/domain"
	"time"
)

type ListAudioVariantsByTrackInput struct {
	TrackID string
}

type ListAudioVariantsByTrackOutput struct {
	ID           string
	TrackID      string
	Codec        domain.Codec
	BitrateKbps  int
	SampleRateHz int
	Channels     int
	DurationMs   int64
	FileURL      string
	CreatedAt    time.Time
}

type ListAudioVariantsByTrackService interface {
	Execute(ctx context.Context, input ListAudioVariantsByTrackInput) ([]ListAudioVariantsByTrackOutput, error)
}
