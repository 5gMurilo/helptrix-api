package strategies

import (
	"context"
	"fmt"
	"log"
	"strings"

	storageinterfaces "github.com/5gMurilo/helptrix-api/core/interfaces/storage"
	userinterfaces "github.com/5gMurilo/helptrix-api/core/interfaces/user"
	"github.com/5gMurilo/helptrix-api/core/utils"
	"github.com/google/uuid"
)

type ProfileImageStrategy struct {
	storageSvc storageinterfaces.IStorageService
	userRepo   userinterfaces.IUserRepository
	bucketName string
}

func NewProfileImageStrategy(
	storageSvc storageinterfaces.IStorageService,
	userRepo userinterfaces.IUserRepository,
	bucketName string,
) *ProfileImageStrategy {
	return &ProfileImageStrategy{
		storageSvc: storageSvc,
		userRepo:   userRepo,
		bucketName: bucketName,
	}
}

func (s *ProfileImageStrategy) Upload(
	ctx context.Context,
	requesterID uuid.UUID,
	ownerID uuid.UUID,
	filename string,
	data []byte,
	contentType string,
) (string, error) {
	if requesterID != ownerID {
		return "", utils.ErrNotOwner
	}

	currentURL, err := s.userRepo.GetProfilePicture(ownerID)
	if err != nil {
		return "", err
	}

	log.Println("currentURL", currentURL)
	log.Println("contentType", contentType)
	log.Println("bucketName", s.bucketName)

	if currentURL != "" {
		prefix := fmt.Sprintf("https://storage.googleapis.com/%s/", s.bucketName)
		objectPath := strings.TrimPrefix(currentURL, prefix)
		if err := s.storageSvc.DeleteFile(ctx, objectPath); err != nil {
			return "", err
		}
	}

	newURL, err := s.storageSvc.UploadFile(ctx, "profile-images", ownerID.String(), filename, data, contentType)
	if err != nil {
		log.Println("error uploading file", err)
		return "", err
	}

	if err := s.userRepo.UpdateProfilePicture(ownerID, newURL); err != nil {
		return "", err
	}

	return newURL, nil
}
