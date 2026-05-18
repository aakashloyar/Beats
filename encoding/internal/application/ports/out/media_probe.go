package out

import (
	"context"
)

type AudioMetadata struct {
	Duration   float64
	Codec      string
	SampleRate int
	Channels   int
	Bitrate    int64
}

type MediaProbe interface {
	Probe(ctx context.Context, path string) (AudioMetadata, error)
}