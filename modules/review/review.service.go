package review

import (
	"github.com/5gMurilo/helptrix-api/core/domain"
	reviewinterfaces "github.com/5gMurilo/helptrix-api/core/interfaces/review"
	"github.com/5gMurilo/helptrix-api/core/utils"
	"github.com/google/uuid"
)

type ReviewService struct {
	repo reviewinterfaces.IReviewRepository
}

func NewReviewService(repo reviewinterfaces.IReviewRepository) reviewinterfaces.IReviewService {
	return &ReviewService{repo: repo}
}

func (s *ReviewService) CreateReview(businessID uuid.UUID, dto domain.CreateReviewRequestDTO) error {
	proposalID, err := uuid.Parse(dto.ProposalID)
	if err != nil {
		return err
	}

	helperID, err := uuid.Parse(dto.HelperID)
	if err != nil {
		return err
	}

	if businessID == helperID {
		return utils.ErrCannotReviewSelf
	}

	review := &domain.Review{
		ProposalID:  proposalID,
		BusinessID:  businessID,
		HelperID:    helperID,
		Rate:        dto.Rate,
		Review:      dto.Review,
		ServiceType: dto.ServiceType,
	}

	return s.repo.Create(review)
}

func (s *ReviewService) ListBusinessReviews(businessID uuid.UUID, p domain.PaginationParams) ([]domain.ReviewListResponseDTO, error) {
	reviews, err := s.repo.ListByBusiness(businessID, p)
	if err != nil {
		return nil, err
	}

	result := make([]domain.ReviewListResponseDTO, 0, len(reviews))
	for _, r := range reviews {
		result = append(result, domain.ReviewListResponseDTO{
			Rate:        r.Rate,
			Review:      "",
			ServiceType: r.ServiceType,
			CreatedAt:   r.CreatedAt,
		})
	}

	return result, nil
}

func (s *ReviewService) ListHelperReviews(helperID uuid.UUID, p domain.PaginationParams) ([]domain.ReviewListResponseDTO, error) {
	reviews, err := s.repo.ListByHelper(helperID, p)
	if err != nil {
		return nil, err
	}

	result := make([]domain.ReviewListResponseDTO, 0, len(reviews))
	for _, r := range reviews {
		result = append(result, domain.ReviewListResponseDTO{
			Rate:        r.Rate,
			Review:      r.Review,
			ServiceType: r.ServiceType,
			CreatedAt:   r.CreatedAt,
		})
	}

	return result, nil
}
