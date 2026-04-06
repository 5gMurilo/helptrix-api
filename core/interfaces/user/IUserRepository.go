package userinterfaces

import (
	"github.com/5gMurilo/helptrix-api/core/domain"
	"github.com/google/uuid"
)

type IUserRepository interface {
	GetProfile(userID uuid.UUID, filters domain.ProfileFilters) (domain.GetProfileResponseDTO, error)
	UpdateProfile(userID uuid.UUID, dto domain.UpdateProfileRequestDTO) error
	DeleteProfile(userID uuid.UUID) error
	GetProfilePicture(userID uuid.UUID) (string, error)
	UpdateProfilePicture(userID uuid.UUID, url string) error
}
