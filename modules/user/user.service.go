package user

import (
	"log/slog"

	"github.com/5gMurilo/helptrix-api/core/domain"
	userinterfaces "github.com/5gMurilo/helptrix-api/core/interfaces/user"
	"github.com/5gMurilo/helptrix-api/core/logger"
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
	log := logger.Get().With(
		slog.String("layer", "service"),
		slog.String("module", "user"),
		slog.String("operation", "GetProfile"),
		slog.String("requester_id", requesterID.String()),
		slog.String("target_id", targetID.String()),
	)

	if requesterType != utils.UserTypeBusiness {
		filters = domain.ProfileFilters{}
	}

	result, err := s.repo.GetProfile(targetID, filters)
	if err != nil {
		log.Error("failed to get user profile", slog.String("error", err.Error()))
		return domain.GetProfileResponseDTO{}, err
	}

	return result, nil
}

func (s *UserService) UpdateProfile(requesterID uuid.UUID, targetID uuid.UUID, dto domain.UpdateProfileRequestDTO) error {
	log := logger.Get().With(
		slog.String("layer", "service"),
		slog.String("module", "user"),
		slog.String("operation", "UpdateProfile"),
		slog.String("requester_id", requesterID.String()),
		slog.String("target_id", targetID.String()),
	)

	if requesterID != targetID {
		log.Warn("update profile rejected: requester is not the profile owner")
		return utils.ErrNotOwner
	}

	if err := s.repo.UpdateProfile(targetID, dto); err != nil {
		log.Error("failed to update user profile", slog.String("error", err.Error()))
		return err
	}

	log.Info("user profile updated successfully")
	return nil
}

func (s *UserService) DeleteProfile(requesterID uuid.UUID, targetID uuid.UUID) error {
	log := logger.Get().With(
		slog.String("layer", "service"),
		slog.String("module", "user"),
		slog.String("operation", "DeleteProfile"),
		slog.String("requester_id", requesterID.String()),
		slog.String("target_id", targetID.String()),
	)

	if requesterID != targetID {
		log.Warn("delete profile rejected: requester is not the profile owner")
		return utils.ErrNotOwner
	}

	if err := s.repo.DeleteProfile(targetID); err != nil {
		log.Error("failed to delete user profile", slog.String("error", err.Error()))
		return err
	}

	log.Info("user profile deleted successfully")
	return nil
}
