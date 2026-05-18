package ffmpeg

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"github.com/aakashloyar/beats/encoding/config"
)

type EncoderAdapter struct{}

func NewEncoderAdapter() *EncoderAdapter {
	return &EncoderAdapter{}
}

func (e *EncoderAdapter) EncodeVariants(ctx context.Context, inputPath string, outputDir string) error {

	for _, profile := range config.Profiles {

		// Example:
		// outputDir/64k/
		profileDir := filepath.Join(
			outputDir,
			profile.Bitrate,
		)

		if err := os.MkdirAll(profileDir, 0755); err != nil {
			return err
		}

		playlistPath := filepath.Join(
			profileDir,
			"playlist.m3u8",
		)

		audioPath := filepath.Join(
			profileDir,
			"audio.mp4",
		)

		cmd := exec.CommandContext(
			ctx,
			"ffmpeg",

			// input
			"-i", inputPath,

			// audio codec
			"-c:a", profile.FFmpegCodec,

			// bitrate
			"-b:a", profile.Bitrate,

			// CMAF fragmented MP4 flags
			"-movflags", "frag_keyframe+empty_moov+default_base_moof",

			// HLS using fMP4
			"-hls_segment_type", "fmp4",

			// single-file mode
			"-hls_flags", "single_file",

			// segment duration
			"-hls_time", "6",

			// media file path
			"-hls_segment_filename", audioPath,

			// overwrite -> without this if the file exist then it will ask permission to overwrite or not
			"-y",

			// output playlist
			playlistPath,
		)

		if err := cmd.Run(); err != nil {
			return fmt.Errorf(
				"failed encoding %s: %w",
				profile.Bitrate,
				err,
			)
		}
	}
	// Create master playlist
	err := e.createMasterPlaylist(outputDir)
	if err != nil {
		return err
	}

	return nil
}


func (e *EncoderAdapter) createMasterPlaylist(outputDir string) error {

	masterPath := filepath.Join(
		outputDir,
		"master.m3u8",
	)

	file, err := os.Create(masterPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString("#EXTM3U\n\n")
	if err != nil {
		return err
	}

	for _, profile := range config.Profiles {

		entry := fmt.Sprintf(
			"#EXT-X-STREAM-INF:BANDWIDTH=%d\n%s/playlist.m3u8\n\n",
			profile.Bandwidth,
			profile.Bitrate,
		)

		_, err := file.WriteString(entry)
		if err != nil {
			return err
		}
	}

	return nil
}