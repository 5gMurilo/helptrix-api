package service

import (
	"log/slog"

	"github.com/5gMurilo/helptrix-api/core/domain"
	serviceinterfaces "github.com/5gMurilo/helptrix-api/core/interfaces/service"
	"github.com/5gMurilo/helptrix-api/core/logger"
	"github.com/5gMurilo/helptrix-api/core/utils"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ServiceService struct {
	repo serviceinterfaces.IServiceRepository
}

func NewServiceService(repo serviceinterfaces.IServiceRepository) serviceinterfaces.IServiceService {
	return &ServiceService{repo: repo}
}

func (s *ServiceService) Create(userID uuid.UUID, userType string, dto domain.CreateServiceRequestDTO) (domain.ServiceResponseDTO, error) {
	log := logger.Get().With(
		slog.String("layer", "service"),
		slog.String("module", "service"),
		slog.String("operation", "Create"),
		slog.String("user_id", userID.String()),
	)

	if userType != utils.UserTypeHelper {
		log.Warn("service creation rejected: only helpers can create services", slog.String("user_type", userType))
		return domain.ServiceResponseDTO{}, utils.ErrHelperOnly
	}

	value, err := decimal.NewFromString(dto.Value)
	if err != nil {
		log.Warn("service creation rejected: invalid value format", slog.String("value", dto.Value))
		return domain.ServiceResponseDTO{}, utils.ErrInvalidValueFormat
	}

	if !value.IsPositive() {
		log.Warn("service creation rejected: value must be positive", slog.String("value", dto.Value))
		return domain.ServiceResponseDTO{}, utils.ErrValueNotPositive
	}

	if !utils.TimeHHMMRegex.MatchString(dto.StartTime) {
		log.Warn("service creation rejected: invalid start_time format", slog.String("start_time", dto.StartTime))
		return domain.ServiceResponseDTO{}, utils.ErrInvalidStartTimeFormat
	}

	if !utils.TimeHHMMRegex.MatchString(dto.EndTime) {
		log.Warn("service creation rejected: invalid end_time format", slog.String("end_time", dto.EndTime))
		return domain.ServiceResponseDTO{}, utils.ErrInvalidEndTimeFormat
	}

	hasCategory, err := s.repo.UserHasCategory(userID, dto.CategoryID)
	if err != nil {
		log.Error("failed to check user category", slog.String("error", err.Error()))
		return domain.ServiceResponseDTO{}, err
	}
	if !hasCategory {
		log.Warn("service creation rejected: category not assigned to user", slog.Any("category_id", dto.CategoryID))
		return domain.ServiceResponseDTO{}, utils.ErrCategoryNotAssignedToUser
	}

	exists, err := s.repo.ExistsByNameAndUser(dto.Name, userID)
	if err != nil {
		log.Error("failed to check service name uniqueness", slog.String("error", err.Error()))
		return domain.ServiceResponseDTO{}, err
	}
	if exists {
		log.Warn("service creation rejected: service name already in use", slog.String("name", dto.Name))
		return domain.ServiceResponseDTO{}, utils.ErrServiceNameNotUnique
	}

	result, err := s.repo.Create(userID, dto)
	if err != nil {
		log.Error("failed to create service", slog.String("error", err.Error()))
		return domain.ServiceResponseDTO{}, err
	}

	log.Info("service created successfully", slog.String("service_id", result.ID.String()))
	return result, nil
}

func (s *ServiceService) List(userID uuid.UUID, userType string) ([]domain.ServiceResponseDTO, error) {
	log := logger.Get().With(
		slog.String("layer", "service"),
		slog.String("module", "service"),
		slog.String("operation", "List"),
		slog.String("user_id", userID.String()),
	)

	if userType != utils.UserTypeHelper {
		log.Warn("service list rejected: only helpers can list their services", slog.String("user_type", userType))
		return nil, utils.ErrHelperOnly
	}

	result, err := s.repo.List(userID)
	if err != nil {
		log.Error("failed to list services", slog.String("error", err.Error()))
		return nil, err
	}

	return result, nil
}

func (s *ServiceService) GetByID(serviceID uuid.UUID, userID uuid.UUID, userType string) (domain.ServiceResponseDTO, error) {
	log := logger.Get().With(
		slog.String("layer", "service"),
		slog.String("module", "service"),
		slog.String("operation", "GetByID"),
		slog.String("service_id", serviceID.String()),
		slog.String("user_id", userID.String()),
	)

	if userType != utils.UserTypeHelper {
		log.Warn("service get rejected: only helpers can access their services", slog.String("user_type", userType))
		return domain.ServiceResponseDTO{}, utils.ErrHelperOnly
	}

	result, err := s.repo.GetByID(serviceID, userID)
	if err != nil {
		log.Error("failed to get service by ID", slog.String("error", err.Error()))
		return domain.ServiceResponseDTO{}, err
	}

	return result, nil
}

func (s *ServiceService) Update(serviceID uuid.UUID, userID uuid.UUID, userType string, dto domain.UpdateServiceRequestDTO) (domain.ServiceResponseDTO, error) {
	log := logger.Get().With(
		slog.String("layer", "service"),
		slog.String("module", "service"),
		slog.String("operation", "Update"),
		slog.String("service_id", serviceID.String()),
		slog.String("user_id", userID.String()),
	)

	if userType != utils.UserTypeHelper {
		log.Warn("service update rejected: only helpers can update their services", slog.String("user_type", userType))
		return domain.ServiceResponseDTO{}, utils.ErrHelperOnly
	}

	if dto.Value != nil {
		val, err := decimal.NewFromString(*dto.Value)
		if err != nil {
			log.Warn("service update rejected: invalid value format", slog.String("value", *dto.Value))
			return domain.ServiceResponseDTO{}, utils.ErrInvalidValueFormat
		}
		if !val.IsPositive() {
			log.Warn("service update rejected: value must be positive", slog.String("value", *dto.Value))
			return domain.ServiceResponseDTO{}, utils.ErrValueNotPositive
		}
	}

	if dto.StartTime != nil && !utils.TimeHHMMRegex.MatchString(*dto.StartTime) {
		log.Warn("service update rejected: invalid start_time format", slog.String("start_time", *dto.StartTime))
		return domain.ServiceResponseDTO{}, utils.ErrInvalidStartTimeFormat
	}

	if dto.EndTime != nil && !utils.TimeHHMMRegex.MatchString(*dto.EndTime) {
		log.Warn("service update rejected: invalid end_time format", slog.String("end_time", *dto.EndTime))
		return domain.ServiceResponseDTO{}, utils.ErrInvalidEndTimeFormat
	}

	if dto.CategoryID != nil {
		hasCategory, err := s.repo.UserHasCategory(userID, *dto.CategoryID)
		if err != nil {
			log.Error("failed to check user category on update", slog.String("error", err.Error()))
			return domain.ServiceResponseDTO{}, err
		}
		if !hasCategory {
			log.Warn("service update rejected: category not assigned to user")
			return domain.ServiceResponseDTO{}, utils.ErrCategoryNotAssignedToUser
		}
	}

	if dto.Name != nil {
		exists, err := s.repo.ExistsByNameAndUserExcluding(*dto.Name, userID, serviceID)
		if err != nil {
			log.Error("failed to check service name uniqueness on update", slog.String("error", err.Error()))
			return domain.ServiceResponseDTO{}, err
		}
		if exists {
			log.Warn("service update rejected: service name already in use", slog.String("name", *dto.Name))
			return domain.ServiceResponseDTO{}, utils.ErrServiceNameNotUnique
		}
	}

	result, err := s.repo.Update(serviceID, userID, dto)
	if err != nil {
		log.Error("failed to update service", slog.String("error", err.Error()))
		return domain.ServiceResponseDTO{}, err
	}

	log.Info("service updated successfully")
	return result, nil
}

func (s *ServiceService) Delete(serviceID uuid.UUID, userID uuid.UUID, userType string) error {
	log := logger.Get().With(
		slog.String("layer", "service"),
		slog.String("module", "service"),
		slog.String("operation", "Delete"),
		slog.String("service_id", serviceID.String()),
		slog.String("user_id", userID.String()),
	)

	if userType != utils.UserTypeHelper {
		log.Warn("service delete rejected: only helpers can delete their services", slog.String("user_type", userType))
		return utils.ErrHelperOnly
	}

	if err := s.repo.Delete(serviceID, userID); err != nil {
		log.Error("failed to delete service", slog.String("error", err.Error()))
		return err
	}

	log.Info("service deleted successfully")
	return nil
}
