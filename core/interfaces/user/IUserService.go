package userinterfaces

import (
	"github.com/5gMurilo/helptrix-api/core/domain"
	"github.com/google/uuid"
)

type IUserService interface {
	GetProfile(requesterID uuid.UUID, requesterType string, targetID uuid.UUID, filters domain.ProfileFilters) (domain.GetProfileResponseDTO, error)
	UpdateProfile(requesterID uuid.UUID, targetID uuid.UUID, dto domain.UpdateProfileRequestDTO) error
	DeleteProfile(requesterID uuid.UUID, targetID uuid.UUID) error
}
