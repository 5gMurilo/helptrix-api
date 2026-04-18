package review

import (
	"testing"

	"github.com/5gMurilo/helptrix-api/core/domain"
	"github.com/5gMurilo/helptrix-api/core/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockReviewRepository struct {
	mock.Mock
}

func (m *MockReviewRepository) Create(review *domain.Review) error {
	args := m.Called(review)
	return args.Error(0)
}

func (m *MockReviewRepository) ListByBusiness(businessID uuid.UUID) ([]domain.Review, error) {
	args := m.Called(businessID)
	return args.Get(0).([]domain.Review), args.Error(1)
}

func (m *MockReviewRepository) ListByHelper(helperID uuid.UUID) ([]domain.Review, error) {
	args := m.Called(helperID)
	return args.Get(0).([]domain.Review), args.Error(1)
}

func (m *MockReviewRepository) GetByBusinessAndHelper(businessID, helperID uuid.UUID) (*domain.Review, error) {
	args := m.Called(businessID, helperID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Review), args.Error(1)
}

func TestReviewService_CreateReview_Success(t *testing.T) {
	repo := new(MockReviewRepository)
	svc := NewReviewService(repo)

	businessID := uuid.New()
	helperID := uuid.New()

	dto := domain.CreateReviewRequestDTO{
		ProposalID:  uuid.New().String(),
		HelperID:    helperID.String(),
		Rate:        5,
		Review:      "Great service!",
		ServiceType: "Plumbing",
	}

	repo.On("Create", mock.Anything).Return(nil)

	err := svc.CreateReview(businessID, dto)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestReviewService_CreateReview_SelfReview(t *testing.T) {
	repo := new(MockReviewRepository)
	svc := NewReviewService(repo)

	userID := uuid.New()

	dto := domain.CreateReviewRequestDTO{
		ProposalID:  uuid.New().String(),
		HelperID:    userID.String(),
		Rate:        5,
		Review:      "Self review",
		ServiceType: "Plumbing",
	}

	err := svc.CreateReview(userID, dto)

	assert.Error(t, err)
	assert.Equal(t, utils.ErrCannotReviewSelf, err)
}

func TestReviewService_CreateReview_InvalidUUID(t *testing.T) {
	repo := new(MockReviewRepository)
	svc := NewReviewService(repo)

	businessID := uuid.New()

	dto := domain.CreateReviewRequestDTO{
		ProposalID:  uuid.New().String(),
		HelperID:    "invalid-uuid",
		Rate:        5,
		Review:      "Test",
		ServiceType: "Plumbing",
	}

	err := svc.CreateReview(businessID, dto)

	assert.Error(t, err)
}

func TestReviewService_CreateReview_InvalidProposalUUID(t *testing.T) {
	repo := new(MockReviewRepository)
	svc := NewReviewService(repo)

	businessID := uuid.New()

	dto := domain.CreateReviewRequestDTO{
		ProposalID:  "invalid-uuid",
		HelperID:    uuid.New().String(),
		Rate:        5,
		Review:      "Test",
		ServiceType: "Plumbing",
	}

	err := svc.CreateReview(businessID, dto)

	assert.Error(t, err)
}

func TestReviewService_ListBusinessReviews_Success(t *testing.T) {
	repo := new(MockReviewRepository)
	svc := NewReviewService(repo)

	businessID := uuid.New()

	reviews := []domain.Review{
		{
			Rate:        5,
			ServiceType: "Plumbing",
			CreatedAt:   domain.Review{}.CreatedAt,
		},
	}

	repo.On("ListByBusiness", businessID).Return(reviews, nil)

	result, err := svc.ListBusinessReviews(businessID)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, 5, result[0].Rate)
	assert.Equal(t, "Plumbing", result[0].ServiceType)
	assert.Empty(t, result[0].Review) // Business doesn't see review text in list
	repo.AssertExpectations(t)
}

func TestReviewService_ListHelperReviews_Success(t *testing.T) {
	repo := new(MockReviewRepository)
	svc := NewReviewService(repo)

	helperID := uuid.New()

	reviews := []domain.Review{
		{
			Rate:        4,
			Review:      "Good work",
			ServiceType: "Electrical",
			CreatedAt:   domain.Review{}.CreatedAt,
		},
	}

	repo.On("ListByHelper", helperID).Return(reviews, nil)

	result, err := svc.ListHelperReviews(helperID)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, 4, result[0].Rate)
	assert.Equal(t, "Good work", result[0].Review)
	assert.Equal(t, "Electrical", result[0].ServiceType)
	repo.AssertExpectations(t)
}
