package out

import "context"

type Storage interface {
	DownloadFile(ctx context.Context, key string) (string, error)
	UploadDirectory(ctx context.Context, localPath string, storagePath string) error
}