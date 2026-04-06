package adapterstorage

import (
	"context"
	"os"
	"testing"
)

func TestNewFirebaseStorageClient_MissingBucket(t *testing.T) {
	// Ensure FIREBASE_STORAGE_BUCKET is not set for this test.
	original := os.Getenv("FIREBASE_STORAGE_BUCKET")
	os.Unsetenv("FIREBASE_STORAGE_BUCKET")
	defer func() {
		if original != "" {
			os.Setenv("FIREBASE_STORAGE_BUCKET", original)
		}
	}()

	client, err := NewFirebaseStorageClient(context.Background())

	if err == nil {
		t.Fatal("expected error when FIREBASE_STORAGE_BUCKET is empty, got nil")
	}

	if client != nil {
		t.Error("expected nil client when FIREBASE_STORAGE_BUCKET is empty, got non-nil")
	}
}
