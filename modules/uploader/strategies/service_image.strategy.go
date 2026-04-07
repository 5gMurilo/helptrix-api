package strategies

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/5gMurilo/helptrix-api/core/domain"
	serviceinterfaces "github.com/5gMurilo/helptrix-api/core/interfaces/service"
	storageinterfaces "github.com/5gMurilo/helptrix-api/core/interfaces/storage"
	"github.com/google/uuid"
)

type ServiceImageStrategy struct {
	storageSvc  storageinterfaces.IStorageService
	serviceRepo serviceinterfaces.IServiceRepository
	bucketName  string
}

func NewServiceImageStrategy(
	storageSvc storageinterfaces.IStorageService,
	serviceRepo serviceinterfaces.IServiceRepository,
	bucketName string,
) *ServiceImageStrategy {
	return &ServiceImageStrategy{
		storageSvc:  storageSvc,
		serviceRepo: serviceRepo,
		bucketName:  bucketName,
	}
}

func (s *ServiceImageStrategy) Upload(
	ctx context.Context,
	requesterID uuid.UUID,
	serviceID uuid.UUID,
	filename string,
	data []byte,
	contentType string,
) (string, error) {
	svcDTO, err := s.serviceRepo.GetByID(serviceID, requesterID)
	if err != nil {
		return "", err
	}

	ext := filepath.Ext(filename)
	if ext == "" {
		ext = ".jpg"
	}
	storageFilename := uuid.New().String() + ext

	newURL, err := s.storageSvc.UploadFile(ctx, "service-images", serviceID.String(), storageFilename, data, contentType)
	if err != nil {
		return "", err
	}

	photos := append(append([]string(nil), svcDTO.Photos...), newURL)
	_, err = s.serviceRepo.Update(serviceID, requesterID, domain.UpdateServiceRequestDTO{Photos: photos})
	if err != nil {
		prefix := fmt.Sprintf("https://storage.googleapis.com/%s/", s.bucketName)
		objectPath := strings.TrimPrefix(newURL, prefix)
		_ = s.storageSvc.DeleteFile(ctx, objectPath)
		return "", err
	}

	return newURL, nil
}
