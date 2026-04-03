package user

import (
	"github.com/5gMurilo/helptrix-api/core/domain"
	userinterfaces "github.com/5gMurilo/helptrix-api/core/interfaces/user"
	"github.com/5gMurilo/helptrix-api/core/utils"
	"github.com/google/uuid"
)

type UserService struct {
	repo userinterfaces.IUserRepository
}

func NewUserService(repo userinterfaces.IUserRepository) userinterfaces.IUserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetProfile(requesterID uuid.UUID, requesterType string, targetID uuid.UUID, filters domain.ProfileFilters) (domain.GetProfileResponseDTO, error) {
	if requesterType != utils.UserTypeBusiness {
		filters = domain.ProfileFilters{}
	}
	return s.repo.GetProfile(targetID, filters)
}

func (s *UserService) UpdateProfile(requesterID uuid.UUID, targetID uuid.UUID, dto domain.UpdateProfileRequestDTO) error {
	if requesterID != targetID {
		return utils.ErrNotOwner
	}
	return s.repo.UpdateProfile(targetID, dto)
}

func (s *UserService) DeleteProfile(requesterID uuid.UUID, targetID uuid.UUID) error {
	if requesterID != targetID {
		return utils.ErrNotOwner
	}
	return s.repo.DeleteProfile(targetID)
}
