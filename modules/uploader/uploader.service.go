package uploader

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	uploaderinterfaces "github.com/5gMurilo/helptrix-api/core/interfaces/uploader"
	"github.com/5gMurilo/helptrix-api/core/logger"
	"github.com/5gMurilo/helptrix-api/core/utils"
)

type UploaderService struct {
	strategies map[string]uploaderinterfaces.IImageUploadStrategy
}

func NewUploaderService(strategies map[string]uploaderinterfaces.IImageUploadStrategy) uploaderinterfaces.IUploaderService {
	return &UploaderService{strategies: strategies}
}

func (s *UploaderService) Upload(ctx context.Context, imageType string, requesterID uuid.UUID, ownerID uuid.UUID, filename string, data []byte, contentType string) (string, error) {
	log := logger.Get().With(
		slog.String("layer", "service"),
		slog.String("module", "uploader"),
		slog.String("operation", "Upload"),
		slog.String("requester_id", requesterID.String()),
		slog.String("owner_id", ownerID.String()),
		slog.String("image_type", imageType),
	)

	strategy, ok := s.strategies[imageType]
	if !ok {
		log.Warn("upload rejected: unsupported image type", slog.String("image_type", imageType))
		return "", utils.ErrInvalidImageType
	}

	url, err := strategy.Upload(ctx, requesterID, ownerID, filename, data, contentType)
	if err != nil {
		log.Error("image upload failed", slog.String("error", err.Error()), slog.String("filename", filename))
		return "", err
	}

	log.Info("image uploaded successfully", slog.String("filename", filename))
	return url, nil
}
