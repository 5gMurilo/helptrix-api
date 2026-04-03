package serviceinterfaces

import (
	"github.com/5gMurilo/helptrix-api/core/domain"
	"github.com/google/uuid"
)

type IServiceRepository interface {
	Create(userID uuid.UUID, dto domain.CreateServiceRequestDTO) (domain.ServiceResponseDTO, error)
	ExistsByNameAndUser(name string, userID uuid.UUID) (bool, error)
	ExistsByNameAndUserExcluding(name string, userID uuid.UUID, excludeID uuid.UUID) (bool, error)
	UserHasCategory(userID uuid.UUID, categoryID uint) (bool, error)
	List(userID uuid.UUID) ([]domain.ServiceResponseDTO, error)
	GetByID(serviceID uuid.UUID, userID uuid.UUID) (domain.ServiceResponseDTO, error)
	Update(serviceID uuid.UUID, userID uuid.UUID, dto domain.UpdateServiceRequestDTO) (domain.ServiceResponseDTO, error)
	Delete(serviceID uuid.UUID, userID uuid.UUID) error
}
