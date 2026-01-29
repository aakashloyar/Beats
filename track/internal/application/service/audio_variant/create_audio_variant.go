package audio_variant

import (
	"context"
	"errors"
	"github.com/aakashloyar/beats/track/internal/application/ports/in/audio_variant"
	"github.com/aakashloyar/beats/track/internal/application/ports/out"
	"github.com/aakashloyar/beats/track/internal/domain"
)

type CreateAudioVariantService struct {
	repo  out.AudioVariantRepository
	idGen out.IDGenerator
	clock out.Clock
}

func NewCreateAudioVariantService(repo out.AudioVariantRepository, idGen out.IDGenerator, clock out.Clock) in.CreateAudioVariantService {
	return &CreateAudioVariantService{
		repo:  repo,
		idGen: idGen,
		clock: clock,
	}
}

func (s *CreateAudioVariantService) Execute(ctx context.Context, input in.CreateAudioVariantInput) (in.CreateAudioVariantOutput, error) {

	if input.TrackID == "" {
		return in.CreateAudioVariantOutput{}, errors.New("track_id is required")
	}

	v := domain.AudioVariant{
		ID:           s.idGen.NewID(),
		TrackID:      input.TrackID,
		Codec:        domain.Codec(input.Codec),
		BitrateKbps:  input.BitrateKbps,
		SampleRateHz: input.SampleRateHz,
		Channels:     input.Channels,
		DurationMs:   input.DurationMs,
		FileURL:      input.FileURL,
		CreatedAt:    s.clock.Now(),
	}

	if err := s.repo.Save(v); err != nil {
		return in.CreateAudioVariantOutput{}, err
	}

	return in.CreateAudioVariantOutput{AudioVariantID: v.ID}, nil
}
