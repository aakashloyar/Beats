package ffmpeg

import (
	"context"
	"encoding/json"
	"os/exec"
	"strconv"
	"errors"
	"github.com/aakashloyar/beats/encoding/internal/application/ports/out"
)

type FFprobeAdapter struct{}

func NewFFprobeAdapter() *FFprobeAdapter {
	return &FFprobeAdapter{}
}


type ffprobeOutput struct {
	//array of stream because
	//suppose first 1 contains audio 
	//2nd one contain video or some other audio
	Streams []struct {
		CodecName string `json:"codec_name"`
		CodecType string `json:"codec_type"`
		Channels  int    `json:"channels"`
		SampleRate string `json:"sample_rate"`
	} `json:"streams"`

	Format struct {
		Duration string `json:"duration"`
		BitRate  string `json:"bit_rate"`
	} `json:"format"`
}

func (f *FFprobeAdapter) Probe(ctx context.Context, path string) (out.AudioMetadata, error) {

	cmd := exec.CommandContext(
		ctx,
		"ffprobe",
		"-v", "quiet",
		"-print_format", "json",
		"-show_streams",
		"-show_format",
		path,
	)

	output, err := cmd.Output()
	if err != nil {
		return out.AudioMetadata{}, err
	}

	var result ffprobeOutput
	if err := json.Unmarshal(output, &result); err != nil {
		return out.AudioMetadata{}, err
	}
	if len(result.Streams) == 0 {
		return out.AudioMetadata{}, errors.New("no streams found")
	}
	audioIndex := -1

	for i, stream := range result.Streams {
		if stream.CodecType == "audio" {
			audioIndex = i
			break
		}
	}

	if audioIndex == -1 {
		return out.AudioMetadata{}, errors.New("no audio stream found")
	}

	duration, err := strconv.ParseFloat(
		result.Format.Duration,
		64,
	)
	if err != nil {
		return out.AudioMetadata{}, err
	}

	bitrate, err := strconv.ParseInt(
		result.Format.BitRate,
		10,
		64,
	)
	if err != nil {
		return out.AudioMetadata{}, err 
	}

	sampleRate, err := strconv.Atoi(
		result.Streams[audioIndex].SampleRate,
	)
	if err != nil  {
		return out.AudioMetadata{}, err
	}


	return out.AudioMetadata{
		Duration:   duration,
		Codec:      result.Streams[audioIndex].CodecName,
		SampleRate: sampleRate,
		Channels:   result.Streams[audioIndex].Channels,
		Bitrate:    bitrate,
	}, nil
}