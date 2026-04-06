package uploader

import (
	"context"

	"github.com/google/uuid"

	uploaderinterfaces "github.com/5gMurilo/helptrix-api/core/interfaces/uploader"
	"github.com/5gMurilo/helptrix-api/core/utils"
)

type UploaderService struct {
	strategies map[string]uploaderinterfaces.IImageUploadStrategy
}

func NewUploaderService(strategies map[string]uploaderinterfaces.IImageUploadStrategy) uploaderinterfaces.IUploaderService {
	return &UploaderService{strategies: strategies}
}

func (s *UploaderService) Upload(ctx context.Context, imageType string, requesterID uuid.UUID, ownerID uuid.UUID, filename string, data []byte, contentType string) (string, error) {
	strategy, ok := s.strategies[imageType]
	if !ok {
		return "", utils.ErrInvalidImageType
	}

	return strategy.Upload(ctx, requesterID, ownerID, filename, data, contentType)
}
