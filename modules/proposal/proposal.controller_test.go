package proposal_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/5gMurilo/helptrix-api/adapter/auth"
	"github.com/5gMurilo/helptrix-api/core/domain"
	proposalinterfaces "github.com/5gMurilo/helptrix-api/core/interfaces/proposal"
	proposalmodule "github.com/5gMurilo/helptrix-api/modules/proposal"
	"github.com/5gMurilo/helptrix-api/core/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// mockProposalService implements IProposalService for controller tests.
type mockProposalService struct {
	CreateFn       func(dto domain.CreateProposalRequestDTO, userID uuid.UUID) (domain.ProposalResponseDTO, error)
	GetByIDFn      func(proposalID uuid.UUID, requesterID uuid.UUID) (domain.ProposalResponseDTO, error)
	UpdateStatusFn func(proposalID uuid.UUID, dto domain.UpdateProposalStatusRequestDTO, requesterID uuid.UUID, requesterType string) (domain.ProposalResponseDTO, error)
	ListFn         func(requesterID uuid.UUID, requesterType string, statusFilter string) ([]domain.ProposalResponseDTO, error)
}

func (m *mockProposalService) Create(dto domain.CreateProposalRequestDTO, userID uuid.UUID) (domain.ProposalResponseDTO, error) {
	return m.CreateFn(dto, userID)
}

func (m *mockProposalService) GetByID(proposalID uuid.UUID, requesterID uuid.UUID) (domain.ProposalResponseDTO, error) {
	return m.GetByIDFn(proposalID, requesterID)
}

func (m *mockProposalService) UpdateStatus(proposalID uuid.UUID, dto domain.UpdateProposalStatusRequestDTO, requesterID uuid.UUID, requesterType string) (domain.ProposalResponseDTO, error) {
	return m.UpdateStatusFn(proposalID, dto, requesterID, requesterType)
}

func (m *mockProposalService) List(requesterID uuid.UUID, requesterType string, statusFilter string) ([]domain.ProposalResponseDTO, error) {
	return m.ListFn(requesterID, requesterType, statusFilter)
}

var _ proposalinterfaces.IProposalService = (*mockProposalService)(nil)

// defaultMockService returns a mock service where all methods succeed.
func defaultMockService(userID, helperID uuid.UUID) *mockProposalService {
	resp := sampleResponseDTO(userID, helperID)
	return &mockProposalService{
		CreateFn: func(dto domain.CreateProposalRequestDTO, uid uuid.UUID) (domain.ProposalResponseDTO, error) {
			return resp, nil
		},
		GetByIDFn: func(proposalID uuid.UUID, requesterID uuid.UUID) (domain.ProposalResponseDTO, error) {
			return resp, nil
		},
		UpdateStatusFn: func(proposalID uuid.UUID, dto domain.UpdateProposalStatusRequestDTO, requesterID uuid.UUID, requesterType string) (domain.ProposalResponseDTO, error) {
			updated := resp
			updated.Status = dto.Status
			return updated, nil
		},
		ListFn: func(requesterID uuid.UUID, requesterType string, statusFilter string) ([]domain.ProposalResponseDTO, error) {
			return []domain.ProposalResponseDTO{resp}, nil
		},
	}
}

func sampleResponseDTO(userID, helperID uuid.UUID) domain.ProposalResponseDTO {
	return domain.ProposalResponseDTO{
		ID:          uuid.New(),
		UserID:      userID,
		HelperID:    helperID,
		CategoryID:  1,
		Description: "Preciso de ajuda",
		Value:       100.00,
		Status:      utils.ProposalStatusPending,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// setupProposalRouter creates a gin engine with the proposal controller wired up.
func setupProposalRouter(svc *mockProposalService, userID string, userType string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	ctrl := proposalmodule.NewProposalController(svc)

	injectPayload := func(c *gin.Context) {
		c.Set("authorization_payload", &auth.Payload{
			UserID:   userID,
			UserType: userType,
		})
		c.Next()
	}

	r.POST("/proposal", injectPayload, ctrl.Create)
	r.GET("/proposal", injectPayload, ctrl.List)
	r.GET("/proposal/:id", injectPayload, ctrl.GetByID)
	r.PATCH("/proposal/:id/status", injectPayload, ctrl.UpdateStatus)

	return r
}

func validCreateProposalBody(helperID uuid.UUID) []byte {
	body := map[string]interface{}{
		"helper_id":   helperID.String(),
		"category_id": 1,
		"description": "Preciso de ajuda com encanamento",
		"value":       150.00,
	}
	b, _ := json.Marshal(body)
	return b
}

func validUpdateStatusBody(status string) []byte {
	b, _ := json.Marshal(map[string]string{"status": status})
	return b
}

// --- Create handler tests ---

func TestProposalController_Create_201Created(t *testing.T) {
	userID := uuid.New()
	helperID := uuid.New()
	svc := defaultMockService(userID, helperID)

	router := setupProposalRouter(svc, userID.String(), utils.UserTypeBusiness)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/proposal", bytes.NewBuffer(validCreateProposalBody(helperID)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("esperado status 201, obteve: %d | body: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("erro ao decodificar resposta: %v", err)
	}
	if _, ok := resp["id"]; !ok {
		t.Error("resposta deve conter a chave 'id'")
	}
}

func TestProposalController_Create_403NonBusinessUser(t *testing.T) {
	userID := uuid.New()
	helperID := uuid.New()
	svc := defaultMockService(userID, helperID)

	router := setupProposalRouter(svc, userID.String(), utils.UserTypeHelper)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/proposal", bytes.NewBuffer(validCreateProposalBody(helperID)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("esperado status 403, obteve: %d", w.Code)
	}
}

func TestProposalController_Create_400BadBody(t *testing.T) {
	userID := uuid.New()
	helperID := uuid.New()
	svc := defaultMockService(userID, helperID)

	router := setupProposalRouter(svc, userID.String(), utils.UserTypeBusiness)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/proposal", bytes.NewBuffer([]byte(`{"description":"only"}`)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado status 400, obteve: %d", w.Code)
	}
}

func TestProposalController_Create_409AlreadyActive(t *testing.T) {
	userID := uuid.New()
	helperID := uuid.New()
	svc := defaultMockService(userID, helperID)
	svc.CreateFn = func(dto domain.CreateProposalRequestDTO, uid uuid.UUID) (domain.ProposalResponseDTO, error) {
		return domain.ProposalResponseDTO{}, utils.ErrProposalAlreadyActiveForHelper
	}

	router := setupProposalRouter(svc, userID.String(), utils.UserTypeBusiness)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/proposal", bytes.NewBuffer(validCreateProposalBody(helperID)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("esperado status 409, obteve: %d", w.Code)
	}
}

func TestProposalController_Create_500InternalError(t *testing.T) {
	userID := uuid.New()
	helperID := uuid.New()
	svc := defaultMockService(userID, helperID)
	svc.CreateFn = func(dto domain.CreateProposalRequestDTO, uid uuid.UUID) (domain.ProposalResponseDTO, error) {
		return domain.ProposalResponseDTO{}, errors.New("unexpected db error")
	}

	router := setupProposalRouter(svc, userID.String(), utils.UserTypeBusiness)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/proposal", bytes.NewBuffer(validCreateProposalBody(helperID)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("esperado status 500, obteve: %d", w.Code)
	}
}

// --- GetByID handler tests ---

func TestProposalController_GetByID_200OK(t *testing.T) {
	userID := uuid.New()
	helperID := uuid.New()
	svc := defaultMockService(userID, helperID)
	proposalID := uuid.New()

	router := setupProposalRouter(svc, userID.String(), utils.UserTypeBusiness)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/proposal/"+proposalID.String(), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("esperado status 200, obteve: %d | body: %s", w.Code, w.Body.String())
	}
}

func TestProposalController_GetByID_400InvalidID(t *testing.T) {
	userID := uuid.New()
	helperID := uuid.New()
	svc := defaultMockService(userID, helperID)

	router := setupProposalRouter(svc, userID.String(), utils.UserTypeBusiness)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/proposal/not-a-uuid", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado status 400, obteve: %d", w.Code)
	}
}

func TestProposalController_GetByID_404NotFound(t *testing.T) {
	userID := uuid.New()
	helperID := uuid.New()
	svc := defaultMockService(userID, helperID)
	svc.GetByIDFn = func(proposalID uuid.UUID, requesterID uuid.UUID) (domain.ProposalResponseDTO, error) {
		return domain.ProposalResponseDTO{}, utils.ErrProposalNotFound
	}

	router := setupProposalRouter(svc, userID.String(), utils.UserTypeBusiness)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/proposal/"+uuid.New().String(), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("esperado status 404, obteve: %d", w.Code)
	}
}

func TestProposalController_GetByID_403NotParticipant(t *testing.T) {
	userID := uuid.New()
	helperID := uuid.New()
	svc := defaultMockService(userID, helperID)
	svc.GetByIDFn = func(proposalID uuid.UUID, requesterID uuid.UUID) (domain.ProposalResponseDTO, error) {
		return domain.ProposalResponseDTO{}, utils.ErrNotProposalParticipant
	}

	router := setupProposalRouter(svc, userID.String(), utils.UserTypeBusiness)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/proposal/"+uuid.New().String(), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("esperado status 403, obteve: %d", w.Code)
	}
}

func TestProposalController_GetByID_500InternalError(t *testing.T) {
	userID := uuid.New()
	helperID := uuid.New()
	svc := defaultMockService(userID, helperID)
	svc.GetByIDFn = func(proposalID uuid.UUID, requesterID uuid.UUID) (domain.ProposalResponseDTO, error) {
		return domain.ProposalResponseDTO{}, errors.New("unexpected error")
	}

	router := setupProposalRouter(svc, userID.String(), utils.UserTypeBusiness)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/proposal/"+uuid.New().String(), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("esperado status 500, obteve: %d", w.Code)
	}
}

// --- UpdateStatus handler tests ---

func TestProposalController_UpdateStatus_200OK(t *testing.T) {
	userID := uuid.New()
	helperID := uuid.New()
	svc := defaultMockService(userID, helperID)

	router := setupProposalRouter(svc, helperID.String(), utils.UserTypeHelper)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPatch, "/proposal/"+uuid.New().String()+"/status", bytes.NewBuffer(validUpdateStatusBody(utils.ProposalStatusAccepted)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("esperado status 200, obteve: %d | body: %s", w.Code, w.Body.String())
	}
}

func TestProposalController_UpdateStatus_400InvalidProposalID(t *testing.T) {
	userID := uuid.New()
	helperID := uuid.New()
	svc := defaultMockService(userID, helperID)

	router := setupProposalRouter(svc, helperID.String(), utils.UserTypeHelper)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPatch, "/proposal/bad-id/status", bytes.NewBuffer(validUpdateStatusBody(utils.ProposalStatusAccepted)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado status 400, obteve: %d", w.Code)
	}
}

func TestProposalController_UpdateStatus_400BadBody(t *testing.T) {
	userID := uuid.New()
	helperID := uuid.New()
	svc := defaultMockService(userID, helperID)

	router := setupProposalRouter(svc, helperID.String(), utils.UserTypeHelper)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPatch, "/proposal/"+uuid.New().String()+"/status", bytes.NewBuffer([]byte(`{}`)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado status 400 para body sem status, obteve: %d", w.Code)
	}
}

func TestProposalController_UpdateStatus_404NotFound(t *testing.T) {
	userID := uuid.New()
	helperID := uuid.New()
	svc := defaultMockService(userID, helperID)
	svc.UpdateStatusFn = func(proposalID uuid.UUID, dto domain.UpdateProposalStatusRequestDTO, requesterID uuid.UUID, requesterType string) (domain.ProposalResponseDTO, error) {
		return domain.ProposalResponseDTO{}, utils.ErrProposalNotFound
	}

	router := setupProposalRouter(svc, helperID.String(), utils.UserTypeHelper)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPatch, "/proposal/"+uuid.New().String()+"/status", bytes.NewBuffer(validUpdateStatusBody(utils.ProposalStatusAccepted)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("esperado status 404, obteve: %d", w.Code)
	}
}

func TestProposalController_UpdateStatus_403Unauthorized(t *testing.T) {
	userID := uuid.New()
	helperID := uuid.New()
	svc := defaultMockService(userID, helperID)
	svc.UpdateStatusFn = func(proposalID uuid.UUID, dto domain.UpdateProposalStatusRequestDTO, requesterID uuid.UUID, requesterType string) (domain.ProposalResponseDTO, error) {
		return domain.ProposalResponseDTO{}, utils.ErrProposalUnauthorized
	}

	router := setupProposalRouter(svc, userID.String(), utils.UserTypeBusiness)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPatch, "/proposal/"+uuid.New().String()+"/status", bytes.NewBuffer(validUpdateStatusBody(utils.ProposalStatusAccepted)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("esperado status 403, obteve: %d", w.Code)
	}
}

func TestProposalController_UpdateStatus_422InvalidStatus(t *testing.T) {
	userID := uuid.New()
	helperID := uuid.New()
	svc := defaultMockService(userID, helperID)
	svc.UpdateStatusFn = func(proposalID uuid.UUID, dto domain.UpdateProposalStatusRequestDTO, requesterID uuid.UUID, requesterType string) (domain.ProposalResponseDTO, error) {
		return domain.ProposalResponseDTO{}, utils.ErrProposalInvalidStatus
	}

	router := setupProposalRouter(svc, helperID.String(), utils.UserTypeHelper)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPatch, "/proposal/"+uuid.New().String()+"/status", bytes.NewBuffer(validUpdateStatusBody(utils.ProposalStatusAccepted)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("esperado status 422, obteve: %d", w.Code)
	}
}

func TestProposalController_UpdateStatus_422TerminalStatus(t *testing.T) {
	userID := uuid.New()
	helperID := uuid.New()
	svc := defaultMockService(userID, helperID)
	svc.UpdateStatusFn = func(proposalID uuid.UUID, dto domain.UpdateProposalStatusRequestDTO, requesterID uuid.UUID, requesterType string) (domain.ProposalResponseDTO, error) {
		return domain.ProposalResponseDTO{}, utils.ErrProposalFinished
	}

	router := setupProposalRouter(svc, helperID.String(), utils.UserTypeHelper)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPatch, "/proposal/"+uuid.New().String()+"/status", bytes.NewBuffer(validUpdateStatusBody(utils.ProposalStatusAccepted)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("esperado status 422, obteve: %d", w.Code)
	}
}

func TestProposalController_UpdateStatus_500InternalError(t *testing.T) {
	userID := uuid.New()
	helperID := uuid.New()
	svc := defaultMockService(userID, helperID)
	svc.UpdateStatusFn = func(proposalID uuid.UUID, dto domain.UpdateProposalStatusRequestDTO, requesterID uuid.UUID, requesterType string) (domain.ProposalResponseDTO, error) {
		return domain.ProposalResponseDTO{}, errors.New("unexpected error")
	}

	router := setupProposalRouter(svc, helperID.String(), utils.UserTypeHelper)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPatch, "/proposal/"+uuid.New().String()+"/status", bytes.NewBuffer(validUpdateStatusBody(utils.ProposalStatusAccepted)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("esperado status 500, obteve: %d", w.Code)
	}
}

// --- List handler tests ---

func TestProposalController_List_200OK(t *testing.T) {
	userID := uuid.New()
	helperID := uuid.New()
	svc := defaultMockService(userID, helperID)

	router := setupProposalRouter(svc, userID.String(), utils.UserTypeBusiness)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/proposal", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("esperado status 200, obteve: %d | body: %s", w.Code, w.Body.String())
	}

	var resp []interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("erro ao decodificar resposta como array: %v", err)
	}
}

func TestProposalController_List_WithStatusFilter(t *testing.T) {
	userID := uuid.New()
	helperID := uuid.New()
	svc := defaultMockService(userID, helperID)

	capturedFilter := ""
	svc.ListFn = func(requesterID uuid.UUID, requesterType string, statusFilter string) ([]domain.ProposalResponseDTO, error) {
		capturedFilter = statusFilter
		return []domain.ProposalResponseDTO{}, nil
	}

	router := setupProposalRouter(svc, userID.String(), utils.UserTypeBusiness)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/proposal?status=accepted", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("esperado status 200, obteve: %d", w.Code)
	}
	if capturedFilter != "accepted" {
		t.Errorf("esperado statusFilter 'accepted', obteve '%s'", capturedFilter)
	}
}

func TestProposalController_List_500InternalError(t *testing.T) {
	userID := uuid.New()
	helperID := uuid.New()
	svc := defaultMockService(userID, helperID)
	svc.ListFn = func(requesterID uuid.UUID, requesterType string, statusFilter string) ([]domain.ProposalResponseDTO, error) {
		return nil, errors.New("db error")
	}

	router := setupProposalRouter(svc, userID.String(), utils.UserTypeBusiness)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/proposal", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("esperado status 500, obteve: %d", w.Code)
	}
}
