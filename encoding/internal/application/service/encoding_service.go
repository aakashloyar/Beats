package service

import (
	"context"

	"github.com/aakashloyar/beats/encoding/config"
	"github.com/aakashloyar/beats/encoding/internal/application/ports/in"
	"github.com/aakashloyar/beats/encoding/internal/application/ports/out"
	"path/filepath"
	"os"
)
type EncodingService struct {
	storage   out.Storage
	probe     out.MediaProbe
	validator out.Validtor
	encoder   out.Encoder
	publisher out.EventPublisher
}

func NewEncodingService(storage out.Storage, probe out.MediaProbe, validator out.Validtor, encoder out.Encoder, publisher out.EventPublisher) in.EncodingService {
	return &EncodingService{
		storage: storage,
		probe: probe,
		encoder: encoder,
		publisher: publisher,
	}
}

func (s *EncodingService) Execute(ctx context.Context, event in.UploadCompletedEvent) error {

	// 1. Download file from S3
	filePath, err := s.storage.DownloadFile(ctx, event.StorageKey)
	if err != nil {
		return err
	}

	//2. Extract Metadata 
	meta, err := s.probe.Probe(ctx, filePath)

	if err != nil {
		return err
	}

    // 3. Validate file
	err = s.validator.ValidateAudio(meta)
	if err != nil {
		return err
	}

	// 4. Encode variants
	jobID := event.UploadID
	outDir := filepath.Join(config.OutPutDirForTranscodedFiles, jobID)
	err = s.encoder.EncodeVariants(ctx, filePath, outDir)
	if err != nil {
		return err
	}

	// 5. Upload encoded files
	storagePath := filepath.Join(
		config.StoragePathForTranscodedFiles,
		jobID,
	)

	err = s.storage.UploadDirectory(ctx, outDir, storagePath)
	if err != nil {
		return err
	}

	streamKey := filepath.Join(
		storagePath,
		"master.m3u8",
	)
	defer os.RemoveAll(outDir)

	// 6. Publish event to core service
	err = s.publisher.PublishUploadCompleted(ctx, out.EncodeCompletedEvent{
		ID:            event.UploadID,
		Title:         event.Title,
		ArtistIDs:      event.ArtistIDs,
		AlbumID:       event.AlbumID,
		CoverImageURL: event.CoverImageURL,
		ReleasedAt:    event.ReleasedAt,
		Language:      event.Language,
		StreamKey:     streamKey,
	})
	return err
}