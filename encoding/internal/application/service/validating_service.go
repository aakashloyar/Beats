package service

import (
	"errors"
	"fmt"

	"github.com/aakashloyar/beats/encoding/internal/application/ports/out"
	"github.com/aakashloyar/beats/encoding/config"
)

type Validator struct {

}

func NewValidator() out.Validtor {
	return &Validator{}
}
func (v *Validator)ValidateAudio(meta out.AudioMetadata) error {
	// 1. Duration validation
	if meta.Duration <= 0 {
		return errors.New("invalid audio duration")
	}

	if meta.Duration < 5 {
		return errors.New("audio too short (minimum 5 seconds)")
	}

	if meta.Duration > 2*60*60 {
		return errors.New("audio too long (maximum 2 hours)")
	}

	// 2. Codec validation
	if !config.SupportedCodecs[meta.Codec] {
		return fmt.Errorf(
			"unsupported codec: %s",
			meta.Codec,
		)
	}

	
	//3. SampleRate Validation
	if !config.SupportedSampleRates[meta.SampleRate] {
		return fmt.Errorf(
			"unsupported sample rate: %d",
			meta.SampleRate,
		)
	}

	// 4. Channel validation
	if meta.Channels != 1 && meta.Channels != 2 {
		return fmt.Errorf(
			"unsupported channel count: %d",
			meta.Channels,
		)
	}

	// 5. Bitrate sanity check
	if meta.Bitrate <= 0 {
		return errors.New("invalid bitrate")
	}

	return nil
}