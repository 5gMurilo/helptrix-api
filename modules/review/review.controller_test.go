package review

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/5gMurilo/helptrix-api/adapter/auth"
	"github.com/5gMurilo/helptrix-api/core/domain"
	"github.com/5gMurilo/helptrix-api/core/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockReviewService struct {
	mock.Mock
}

func (m *MockReviewService) CreateReview(businessID uuid.UUID, dto domain.CreateReviewRequestDTO) error {
	args := m.Called(businessID, dto)
	return args.Error(0)
}

func (m *MockReviewService) ListBusinessReviews(businessID uuid.UUID) ([]domain.ReviewListResponseDTO, error) {
	args := m.Called(businessID)
	return args.Get(0).([]domain.ReviewListResponseDTO), args.Error(1)
}

func (m *MockReviewService) ListHelperReviews(helperID uuid.UUID) ([]domain.ReviewListResponseDTO, error) {
	args := m.Called(helperID)
	return args.Get(0).([]domain.ReviewListResponseDTO), args.Error(1)
}

func TestReviewController_Create_NonBusinessUser(t *testing.T) {
	svc := new(MockReviewService)

	payload := &auth.Payload{
		UserID:   uuid.New().String(),
		UserType: utils.UserTypeHelper,
	}

	req := httptest.NewRequest(http.MethodPost, "/review", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("authorization_payload", payload)

	ctrl := NewReviewController(svc)
	ctrl.Create(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestReviewController_Create_InvalidBody(t *testing.T) {
	svc := new(MockReviewService)

	payload := &auth.Payload{
		UserID:   uuid.New().String(),
		UserType: utils.UserTypeBusiness,
	}

	body := bytes.NewBufferString(`{"invalid": "json"`)
	req := httptest.NewRequest(http.MethodPost, "/review", body)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("authorization_payload", payload)

	ctrl := NewReviewController(svc)
	ctrl.Create(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestReviewController_Create_ReviewAlreadyExists(t *testing.T) {
	svc := new(MockReviewService)

	businessID := uuid.New()
	payload := &auth.Payload{
		UserID:   businessID.String(),
		UserType: utils.UserTypeBusiness,
	}

	dto := domain.CreateReviewRequestDTO{
		HelperID:    uuid.New().String(),
		Rate:        5,
		Review:      "Test",
		ServiceType: "Test",
	}

	body, _ := json.Marshal(dto)
	req := httptest.NewRequest(http.MethodPost, "/review", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("authorization_payload", payload)

	svc.On("CreateReview", mock.Anything, mock.Anything).Return(utils.ErrReviewAlreadyExists)

	ctrl := NewReviewController(svc)
	ctrl.Create(c)

	assert.Equal(t, http.StatusConflict, w.Code)
	assert.Contains(t, w.Body.String(), utils.ErrReviewAlreadyExists.Error())
}

func TestReviewController_Create_Success(t *testing.T) {
	svc := new(MockReviewService)

	businessID := uuid.New()
	payload := &auth.Payload{
		UserID:   businessID.String(),
		UserType: utils.UserTypeBusiness,
	}

	dto := domain.CreateReviewRequestDTO{
		HelperID:    uuid.New().String(),
		Rate:        5,
		Review:      "Great service!",
		ServiceType: "Plumbing",
	}

	body, _ := json.Marshal(dto)
	req := httptest.NewRequest(http.MethodPost, "/review", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("authorization_payload", payload)

	svc.On("CreateReview", mock.Anything, mock.Anything).Return(nil)

	ctrl := NewReviewController(svc)
	ctrl.Create(c)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestReviewController_ListBusiness_NonBusinessUser(t *testing.T) {
	svc := new(MockReviewService)

	payload := &auth.Payload{
		UserID:   uuid.New().String(),
		UserType: utils.UserTypeHelper,
	}

	req := httptest.NewRequest(http.MethodGet, "/review/business", nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("authorization_payload", payload)

	ctrl := NewReviewController(svc)
	ctrl.ListBusiness(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestReviewController_ListBusiness_Success(t *testing.T) {
	svc := new(MockReviewService)

	businessID := uuid.New()
	payload := &auth.Payload{
		UserID:   businessID.String(),
		UserType: utils.UserTypeBusiness,
	}

	reviews := []domain.ReviewListResponseDTO{
		{Rate: 5, ServiceType: "Plumbing"},
	}

	req := httptest.NewRequest(http.MethodGet, "/review/business", nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("authorization_payload", payload)

	svc.On("ListBusinessReviews", businessID).Return(reviews, nil)

	ctrl := NewReviewController(svc)
	ctrl.ListBusiness(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestReviewController_ListHelper_Success(t *testing.T) {
	svc := new(MockReviewService)

	helperID := uuid.New()
	payload := &auth.Payload{
		UserID:   helperID.String(),
		UserType: utils.UserTypeHelper,
	}

	reviews := []domain.ReviewListResponseDTO{
		{Rate: 4, Review: "Good work", ServiceType: "Electrical"},
	}

	req := httptest.NewRequest(http.MethodGet, "/review/helper", nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("authorization_payload", payload)

	svc.On("ListHelperReviews", helperID).Return(reviews, nil)

	ctrl := NewReviewController(svc)
	ctrl.ListHelper(c)

	assert.Equal(t, http.StatusOK, w.Code)
}
