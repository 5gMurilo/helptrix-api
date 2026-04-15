package reviewinterfaces

import (
	"github.com/5gMurilo/helptrix-api/core/domain"
	"github.com/google/uuid"
)

type IReviewRepository interface {
	Create(review *domain.Review) error
	ListByBusiness(businessID uuid.UUID) ([]domain.Review, error)
	ListByHelper(helperID uuid.UUID) ([]domain.Review, error)
	GetByBusinessAndHelper(businessID, helperID uuid.UUID) (*domain.Review, error)
}
