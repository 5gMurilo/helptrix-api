package reviewinterfaces

import (
	"github.com/5gMurilo/helptrix-api/core/domain"
	"github.com/google/uuid"
)

type IReviewService interface {
	CreateReview(businessID uuid.UUID, dto domain.CreateReviewRequestDTO) error
	ListBusinessReviews(businessID uuid.UUID, p domain.PaginationParams) ([]domain.ReviewListResponseDTO, error)
	ListHelperReviews(helperID uuid.UUID, p domain.PaginationParams) ([]domain.ReviewListResponseDTO, error)
}
