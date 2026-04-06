package adapterstorage

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go/v4"
	storageinterfaces "github.com/5gMurilo/helptrix-api/core/interfaces/storage"
	"google.golang.org/api/option"
)

// FirebaseStorageClient wraps a GCS BucketHandle and implements IStorageService.
type FirebaseStorageClient struct {
	bucket     *storage.BucketHandle
	bucketName string
}

// NewFirebaseStorageClient initialises a Firebase Storage client.
//
// Authentication:
//   - If FIREBASE_CREDENTIALS_PATH is set, that JSON key file is used.
//   - Otherwise, Application Default Credentials (ADC) are used.
//
// FIREBASE_STORAGE_BUCKET must be set; an error is returned if it is empty.
func NewFirebaseStorageClient(ctx context.Context) (storageinterfaces.IStorageService, error) {
	bucketName := strings.TrimPrefix(os.Getenv("FIREBASE_STORAGE_BUCKET"), "gs://")
	if bucketName == "" {
		return nil, errors.New("FIREBASE_STORAGE_BUCKET environment variable is required")
	}

	var opts []option.ClientOption

	if credJSON := os.Getenv("FIREBASE_CREDENTIALS"); credJSON != "" {
		opts = append(opts, option.WithCredentialsJSON([]byte(credJSON)))
	}

	cfg := &firebase.Config{StorageBucket: bucketName}

	app, err := firebase.NewApp(ctx, cfg, opts...)
	if err != nil {
		return nil, fmt.Errorf("error initialising Firebase app: %w", err)
	}

	client, err := app.Storage(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting Firebase Storage client: %w", err)
	}

	bucket, err := client.DefaultBucket()
	if err != nil {
		return nil, fmt.Errorf("error getting default bucket: %w", err)
	}

	return &FirebaseStorageClient{
		bucket:     bucket,
		bucketName: bucketName,
	}, nil
}

// UploadFile uploads data to <folder>/<ownerID>/<filename>, sets the object ACL to
// public-read, and returns the canonical GCS public URL.
func (c *FirebaseStorageClient) UploadFile(
	ctx context.Context,
	folder string,
	ownerID string,
	filename string,
	data []byte,
	contentType string,
) (string, error) {
	objectPath := fmt.Sprintf("%s/%s/%s", folder, ownerID, filename)

	obj := c.bucket.Object(objectPath)
	writer := obj.NewWriter(ctx)
	writer.ContentType = contentType

	if _, err := writer.Write(data); err != nil {
		_ = writer.Close()
		return "", fmt.Errorf("error writing object %q: %w", objectPath, err)
	}

	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("error closing writer for object %q: %w", objectPath, err)
	}

	if err := obj.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return "", fmt.Errorf("error setting public-read ACL on object %q: %w", objectPath, err)
	}

	url := fmt.Sprintf("https://storage.googleapis.com/%s/%s", c.bucketName, objectPath)
	return url, nil
}

// DeleteFile removes the object at objectPath from the bucket.
// If the object does not exist the operation is treated as a no-op (idempotent).
func (c *FirebaseStorageClient) DeleteFile(ctx context.Context, objectPath string) error {
	if err := c.bucket.Object(objectPath).Delete(ctx); err != nil {
		if errors.Is(err, storage.ErrObjectNotExist) {
			return nil
		}
		return fmt.Errorf("error deleting object %q: %w", objectPath, err)
	}
	return nil
}
