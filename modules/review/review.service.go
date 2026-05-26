package review

import (
	"log/slog"

	"github.com/5gMurilo/helptrix-api/core/domain"
	reviewinterfaces "github.com/5gMurilo/helptrix-api/core/interfaces/review"
	"github.com/5gMurilo/helptrix-api/core/logger"
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
	log := logger.Get().With(
		slog.String("layer", "service"),
		slog.String("module", "review"),
		slog.String("operation", "CreateReview"),
		slog.String("business_id", businessID.String()),
	)

	proposalID, err := uuid.Parse(dto.ProposalID)
	if err != nil {
		log.Warn("review creation rejected: invalid proposal_id format", slog.String("proposal_id", dto.ProposalID))
		return err
	}

	helperID, err := uuid.Parse(dto.HelperID)
	if err != nil {
		log.Warn("review creation rejected: invalid helper_id format", slog.String("helper_id", dto.HelperID))
		return err
	}

	if businessID == helperID {
		log.Warn("review creation rejected: user cannot review themselves")
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

	if err := s.repo.Create(review); err != nil {
		log.Error("failed to create review",
			slog.String("error", err.Error()),
			slog.String("helper_id", helperID.String()),
			slog.String("proposal_id", proposalID.String()),
		)
		return err
	}

	log.Info("review created successfully",
		slog.String("helper_id", helperID.String()),
		slog.String("proposal_id", proposalID.String()),
	)
	return nil
}

func (s *ReviewService) ListBusinessReviews(businessID uuid.UUID) ([]domain.ReviewListResponseDTO, error) {
	log := logger.Get().With(
		slog.String("layer", "service"),
		slog.String("module", "review"),
		slog.String("operation", "ListBusinessReviews"),
		slog.String("business_id", businessID.String()),
	)

	reviews, err := s.repo.ListByBusiness(businessID)
	if err != nil {
		log.Error("failed to list business reviews", slog.String("error", err.Error()))
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

func (s *ReviewService) ListHelperReviews(helperID uuid.UUID) ([]domain.ReviewListResponseDTO, error) {
	log := logger.Get().With(
		slog.String("layer", "service"),
		slog.String("module", "review"),
		slog.String("operation", "ListHelperReviews"),
		slog.String("helper_id", helperID.String()),
	)

	reviews, err := s.repo.ListByHelper(helperID)
	if err != nil {
		log.Error("failed to list helper reviews", slog.String("error", err.Error()))
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
