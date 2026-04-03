package serviceinterfaces

import (
	"github.com/5gMurilo/helptrix-api/core/domain"
	"github.com/google/uuid"
)

type IServiceService interface {
	Create(userID uuid.UUID, userType string, dto domain.CreateServiceRequestDTO) (domain.ServiceResponseDTO, error)
	List(userID uuid.UUID, userType string) ([]domain.ServiceResponseDTO, error)
	GetByID(serviceID uuid.UUID, userID uuid.UUID, userType string) (domain.ServiceResponseDTO, error)
	Update(serviceID uuid.UUID, userID uuid.UUID, userType string, dto domain.UpdateServiceRequestDTO) (domain.ServiceResponseDTO, error)
	Delete(serviceID uuid.UUID, userID uuid.UUID, userType string) error
}
