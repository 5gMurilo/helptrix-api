package uploaderinterfaces

import (
	"context"

	"github.com/google/uuid"
)

type IImageUploadStrategy interface {
	Upload(ctx context.Context, requesterID uuid.UUID, ownerID uuid.UUID, filename string, data []byte, contentType string) (string, error)
}
