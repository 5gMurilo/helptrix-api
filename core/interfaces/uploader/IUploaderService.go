package uploaderinterfaces

import (
	"context"

	"github.com/google/uuid"
)

type IUploaderService interface {
	Upload(ctx context.Context, imageType string, requesterID uuid.UUID, ownerID uuid.UUID, filename string, data []byte, contentType string) (string, error)
}
