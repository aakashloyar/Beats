package s3

import (
	"io"
	"os"
	"context"
	//"time"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	//"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aakashloyar/beats/encoding/config"
	"path/filepath"
)

type S3Storage struct {
	client     *s3.Client
	presigner  *s3.PresignClient
	bucketName string
}

func NewS3Storage(client *s3.Client, bucket string) *S3Storage {
	return &S3Storage{
		client:     client,
		presigner:  s3.NewPresignClient(client),
		bucketName: bucket,
	}
}

func (s *S3Storage) DownloadFile(ctx context.Context, key string) (string, error) {
	result, err:= s.client.GetObject(ctx,&s3.GetObjectInput{
		Bucket: &s.bucketName,
		Key: &key,
	})
	if err != nil {
		return "", err 
	}
	defer result.Body.Close()
	tmpFile, err := os.CreateTemp(config.UploadLocalPathForEncoding,"*.flac")

	if err != nil {
		return "", err 
	}
	defer tmpFile.Close()
	_, err = io.Copy(tmpFile, result.Body)
	if err != nil {
		return "", err 
	} 
	return tmpFile.Name(),nil 
}

func (s *S3Storage) UploadDirectory(ctx context.Context, localDir string,storagePath string) error {
	return filepath.Walk(localDir, 
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// skip directories
			if info.IsDir() {
				return nil
			}

			// relative path from localDir
			relPath, err := filepath.Rel(
				localDir,
				path,
			)
			if err != nil {
				return err
			}

			// build S3 key
			s3Key := filepath.Join(
				storagePath,
				relPath,
			)

			// upload file
			err = s.uploadSingleFile(
				ctx,
				path,
				s3Key,
			)
			if err != nil {
				return err
			}

			return nil
		},
	)
}

func (s *S3Storage) uploadSingleFile(ctx context.Context, localPath string, s3Key string) error {

	file, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = s.client.PutObject(
		ctx,
		&s3.PutObjectInput{
			Bucket:      &s.bucketName,
			Key:         &s3Key,
			Body:        file,
		},
	)

	return err
}