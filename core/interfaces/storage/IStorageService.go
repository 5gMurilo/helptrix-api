package storageinterfaces

import "context"

type IStorageService interface {
	UploadFile(ctx context.Context, folder string, ownerID string, filename string, data []byte, contentType string) (string, error)
	DeleteFile(ctx context.Context, objectPath string) error
}
