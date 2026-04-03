package service

import (
	"github.com/5gMurilo/helptrix-api/core/domain"
	serviceinterfaces "github.com/5gMurilo/helptrix-api/core/interfaces/service"
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
	if userType != utils.UserTypeHelper {
		return domain.ServiceResponseDTO{}, utils.ErrHelperOnly
	}

	value, err := decimal.NewFromString(dto.Value)
	if err != nil {
		return domain.ServiceResponseDTO{}, utils.ErrInvalidValueFormat
	}

	if !value.IsPositive() {
		return domain.ServiceResponseDTO{}, utils.ErrValueNotPositive
	}

	if !utils.TimeHHMMRegex.MatchString(dto.StartTime) {
		return domain.ServiceResponseDTO{}, utils.ErrInvalidStartTimeFormat
	}

	if !utils.TimeHHMMRegex.MatchString(dto.EndTime) {
		return domain.ServiceResponseDTO{}, utils.ErrInvalidEndTimeFormat
	}

	hasCategory, err := s.repo.UserHasCategory(userID, dto.CategoryID)
	if err != nil {
		return domain.ServiceResponseDTO{}, err
	}
	if !hasCategory {
		return domain.ServiceResponseDTO{}, utils.ErrCategoryNotAssignedToUser
	}

	exists, err := s.repo.ExistsByNameAndUser(dto.Name, userID)
	if err != nil {
		return domain.ServiceResponseDTO{}, err
	}
	if exists {
		return domain.ServiceResponseDTO{}, utils.ErrServiceNameNotUnique
	}

	return s.repo.Create(userID, dto)
}

func (s *ServiceService) List(userID uuid.UUID, userType string) ([]domain.ServiceResponseDTO, error) {
	if userType != utils.UserTypeHelper {
		return nil, utils.ErrHelperOnly
	}
	return s.repo.List(userID)
}

func (s *ServiceService) GetByID(serviceID uuid.UUID, userID uuid.UUID, userType string) (domain.ServiceResponseDTO, error) {
	if userType != utils.UserTypeHelper {
		return domain.ServiceResponseDTO{}, utils.ErrHelperOnly
	}
	return s.repo.GetByID(serviceID, userID)
}

func (s *ServiceService) Update(serviceID uuid.UUID, userID uuid.UUID, userType string, dto domain.UpdateServiceRequestDTO) (domain.ServiceResponseDTO, error) {
	if userType != utils.UserTypeHelper {
		return domain.ServiceResponseDTO{}, utils.ErrHelperOnly
	}

	if dto.Value != nil {
		val, err := decimal.NewFromString(*dto.Value)
		if err != nil {
			return domain.ServiceResponseDTO{}, utils.ErrInvalidValueFormat
		}
		if !val.IsPositive() {
			return domain.ServiceResponseDTO{}, utils.ErrValueNotPositive
		}
	}

	if dto.StartTime != nil && !utils.TimeHHMMRegex.MatchString(*dto.StartTime) {
		return domain.ServiceResponseDTO{}, utils.ErrInvalidStartTimeFormat
	}

	if dto.EndTime != nil && !utils.TimeHHMMRegex.MatchString(*dto.EndTime) {
		return domain.ServiceResponseDTO{}, utils.ErrInvalidEndTimeFormat
	}

	if dto.CategoryID != nil {
		hasCategory, err := s.repo.UserHasCategory(userID, *dto.CategoryID)
		if err != nil {
			return domain.ServiceResponseDTO{}, err
		}
		if !hasCategory {
			return domain.ServiceResponseDTO{}, utils.ErrCategoryNotAssignedToUser
		}
	}

	if dto.Name != nil {
		exists, err := s.repo.ExistsByNameAndUserExcluding(*dto.Name, userID, serviceID)
		if err != nil {
			return domain.ServiceResponseDTO{}, err
		}
		if exists {
			return domain.ServiceResponseDTO{}, utils.ErrServiceNameNotUnique
		}
	}

	return s.repo.Update(serviceID, userID, dto)
}

func (s *ServiceService) Delete(serviceID uuid.UUID, userID uuid.UUID, userType string) error {
	if userType != utils.UserTypeHelper {
		return utils.ErrHelperOnly
	}
	return s.repo.Delete(serviceID, userID)
}
