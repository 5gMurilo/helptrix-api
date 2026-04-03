package user_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/5gMurilo/helptrix-api/adapter/auth"
	"github.com/5gMurilo/helptrix-api/core/domain"
	"github.com/5gMurilo/helptrix-api/core/utils"
	usermodule "github.com/5gMurilo/helptrix-api/modules/user"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// mockUserService implementa IUserService para testes do controller.
type mockUserService struct {
	GetProfileFn    func(requesterID uuid.UUID, requesterType string, targetID uuid.UUID, filters domain.ProfileFilters) (domain.GetProfileResponseDTO, error)
	UpdateProfileFn func(requesterID uuid.UUID, targetID uuid.UUID, dto domain.UpdateProfileRequestDTO) error
	DeleteProfileFn func(requesterID uuid.UUID, targetID uuid.UUID) error
}

func (m *mockUserService) GetProfile(requesterID uuid.UUID, requesterType string, targetID uuid.UUID, filters domain.ProfileFilters) (domain.GetProfileResponseDTO, error) {
	return m.GetProfileFn(requesterID, requesterType, targetID, filters)
}

func (m *mockUserService) UpdateProfile(requesterID uuid.UUID, targetID uuid.UUID, dto domain.UpdateProfileRequestDTO) error {
	return m.UpdateProfileFn(requesterID, targetID, dto)
}

func (m *mockUserService) DeleteProfile(requesterID uuid.UUID, targetID uuid.UUID) error {
	return m.DeleteProfileFn(requesterID, targetID)
}

func setupUserRouter(svc *mockUserService) (*gin.Engine, uuid.UUID) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	ctrl := usermodule.NewUserController(svc)

	userID := uuid.New()
	payload := &auth.Payload{
		UserID:   userID.String(),
		UserType: utils.UserTypeHelper,
	}

	r.GET("/user/profile/:id", func(c *gin.Context) {
		c.Set("authorization_payload", payload)
		ctrl.GetProfile(c)
	})
	r.PUT("/user/profile/:id", func(c *gin.Context) {
		c.Set("authorization_payload", payload)
		ctrl.UpdateProfile(c)
	})
	r.DELETE("/user/profile/:id", func(c *gin.Context) {
		c.Set("authorization_payload", payload)
		ctrl.DeleteProfile(c)
	})

	return r, userID
}

// ---------- GetProfile ----------

func TestUserController_GetProfile_200OK(t *testing.T) {
	targetID := uuid.New()

	svc := &mockUserService{
		GetProfileFn: func(requesterID uuid.UUID, requesterType string, targetID uuid.UUID, filters domain.ProfileFilters) (domain.GetProfileResponseDTO, error) {
			return domain.GetProfileResponseDTO{
				ID:       targetID,
				Name:     "Maria Silva",
				Email:    "maria@example.com",
				UserType: utils.UserTypeHelper,
			}, nil
		},
	}

	router, _ := setupUserRouter(svc)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/user/profile/"+targetID.String(), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("esperado status 200, obteve: %d", w.Code)
	}

	var respBody map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &respBody); err != nil {
		t.Fatalf("erro ao decodificar resposta: %v", err)
	}
	if _, ok := respBody["id"]; !ok {
		t.Error("resposta deve conter a chave 'id'")
	}
}

func TestUserController_GetProfile_400BadRequest_InvalidID(t *testing.T) {
	svc := &mockUserService{
		GetProfileFn: func(requesterID uuid.UUID, requesterType string, targetID uuid.UUID, filters domain.ProfileFilters) (domain.GetProfileResponseDTO, error) {
			return domain.GetProfileResponseDTO{}, nil
		},
	}

	router, _ := setupUserRouter(svc)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/user/profile/not-a-uuid", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado status 400, obteve: %d", w.Code)
	}
}

func TestUserController_GetProfile_404NotFound(t *testing.T) {
	targetID := uuid.New()

	svc := &mockUserService{
		GetProfileFn: func(requesterID uuid.UUID, requesterType string, targetID uuid.UUID, filters domain.ProfileFilters) (domain.GetProfileResponseDTO, error) {
			return domain.GetProfileResponseDTO{}, utils.ErrUserNotFound
		},
	}

	router, _ := setupUserRouter(svc)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/user/profile/"+targetID.String(), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("esperado status 404, obteve: %d", w.Code)
	}

	var respBody map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &respBody); err != nil {
		t.Fatalf("erro ao decodificar resposta: %v", err)
	}
	if respBody["error"] != utils.ErrUserNotFound.Error() {
		t.Errorf("mensagem de erro inesperada: %s", respBody["error"])
	}
}

func TestUserController_GetProfile_500InternalServerError(t *testing.T) {
	targetID := uuid.New()

	svc := &mockUserService{
		GetProfileFn: func(requesterID uuid.UUID, requesterType string, targetID uuid.UUID, filters domain.ProfileFilters) (domain.GetProfileResponseDTO, error) {
			return domain.GetProfileResponseDTO{}, errors.New("unexpected db error")
		},
	}

	router, _ := setupUserRouter(svc)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/user/profile/"+targetID.String(), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("esperado status 500, obteve: %d", w.Code)
	}

	var respBody map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &respBody); err != nil {
		t.Fatalf("erro ao decodificar resposta: %v", err)
	}
	if respBody["error"] != "internal server error" {
		t.Errorf("mensagem de erro inesperada: %s", respBody["error"])
	}
}

// ---------- UpdateProfile ----------

func TestUserController_UpdateProfile_204NoContent(t *testing.T) {
	targetID := uuid.New()

	svc := &mockUserService{
		UpdateProfileFn: func(requesterID uuid.UUID, targetID uuid.UUID, dto domain.UpdateProfileRequestDTO) error {
			return nil
		},
	}

	router, userID := setupUserRouter(svc)
	_ = userID
	body, _ := json.Marshal(domain.UpdateProfileRequestDTO{Email: "novo@example.com"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/user/profile/"+targetID.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("esperado status 204, obteve: %d", w.Code)
	}
}

func TestUserController_UpdateProfile_400BadRequest_InvalidID(t *testing.T) {
	svc := &mockUserService{
		UpdateProfileFn: func(requesterID uuid.UUID, targetID uuid.UUID, dto domain.UpdateProfileRequestDTO) error {
			return nil
		},
	}

	router, _ := setupUserRouter(svc)
	body, _ := json.Marshal(domain.UpdateProfileRequestDTO{Email: "novo@example.com"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/user/profile/not-a-uuid", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado status 400, obteve: %d", w.Code)
	}
}

func TestUserController_UpdateProfile_403Forbidden(t *testing.T) {
	targetID := uuid.New()

	svc := &mockUserService{
		UpdateProfileFn: func(requesterID uuid.UUID, targetID uuid.UUID, dto domain.UpdateProfileRequestDTO) error {
			return utils.ErrNotOwner
		},
	}

	router, _ := setupUserRouter(svc)
	body, _ := json.Marshal(domain.UpdateProfileRequestDTO{Email: "novo@example.com"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/user/profile/"+targetID.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("esperado status 403, obteve: %d", w.Code)
	}

	var respBody map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &respBody); err != nil {
		t.Fatalf("erro ao decodificar resposta: %v", err)
	}
	if respBody["error"] != utils.ErrNotOwner.Error() {
		t.Errorf("mensagem de erro inesperada: %s", respBody["error"])
	}
}

func TestUserController_UpdateProfile_404NotFound(t *testing.T) {
	targetID := uuid.New()

	svc := &mockUserService{
		UpdateProfileFn: func(requesterID uuid.UUID, targetID uuid.UUID, dto domain.UpdateProfileRequestDTO) error {
			return utils.ErrUserNotFound
		},
	}

	router, _ := setupUserRouter(svc)
	body, _ := json.Marshal(domain.UpdateProfileRequestDTO{Email: "novo@example.com"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/user/profile/"+targetID.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("esperado status 404, obteve: %d", w.Code)
	}
}

func TestUserController_UpdateProfile_409Conflict(t *testing.T) {
	targetID := uuid.New()

	svc := &mockUserService{
		UpdateProfileFn: func(requesterID uuid.UUID, targetID uuid.UUID, dto domain.UpdateProfileRequestDTO) error {
			return utils.ErrCategoryHasLinkedServices
		},
	}

	router, _ := setupUserRouter(svc)
	body, _ := json.Marshal(domain.UpdateProfileRequestDTO{Categories: []uint{1}})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/user/profile/"+targetID.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("esperado status 409, obteve: %d", w.Code)
	}

	var respBody map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &respBody); err != nil {
		t.Fatalf("erro ao decodificar resposta: %v", err)
	}
	if respBody["error"] != utils.ErrCategoryHasLinkedServices.Error() {
		t.Errorf("mensagem de erro inesperada: %s", respBody["error"])
	}
}

func TestUserController_UpdateProfile_500InternalServerError(t *testing.T) {
	targetID := uuid.New()

	svc := &mockUserService{
		UpdateProfileFn: func(requesterID uuid.UUID, targetID uuid.UUID, dto domain.UpdateProfileRequestDTO) error {
			return errors.New("unexpected db error")
		},
	}

	router, _ := setupUserRouter(svc)
	body, _ := json.Marshal(domain.UpdateProfileRequestDTO{Email: "novo@example.com"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/user/profile/"+targetID.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("esperado status 500, obteve: %d", w.Code)
	}

	var respBody map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &respBody); err != nil {
		t.Fatalf("erro ao decodificar resposta: %v", err)
	}
	if respBody["error"] != "internal server error" {
		t.Errorf("mensagem de erro inesperada: %s", respBody["error"])
	}
}

// ---------- DeleteProfile ----------

func TestUserController_DeleteProfile_204NoContent(t *testing.T) {
	targetID := uuid.New()

	svc := &mockUserService{
		DeleteProfileFn: func(requesterID uuid.UUID, targetID uuid.UUID) error {
			return nil
		},
	}

	router, _ := setupUserRouter(svc)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/user/profile/"+targetID.String(), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("esperado status 204, obteve: %d", w.Code)
	}
}

func TestUserController_DeleteProfile_400BadRequest_InvalidID(t *testing.T) {
	svc := &mockUserService{
		DeleteProfileFn: func(requesterID uuid.UUID, targetID uuid.UUID) error {
			return nil
		},
	}

	router, _ := setupUserRouter(svc)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/user/profile/not-a-uuid", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado status 400, obteve: %d", w.Code)
	}
}

func TestUserController_DeleteProfile_403Forbidden(t *testing.T) {
	targetID := uuid.New()

	svc := &mockUserService{
		DeleteProfileFn: func(requesterID uuid.UUID, targetID uuid.UUID) error {
			return utils.ErrNotOwner
		},
	}

	router, _ := setupUserRouter(svc)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/user/profile/"+targetID.String(), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("esperado status 403, obteve: %d", w.Code)
	}

	var respBody map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &respBody); err != nil {
		t.Fatalf("erro ao decodificar resposta: %v", err)
	}
	if respBody["error"] != utils.ErrNotOwner.Error() {
		t.Errorf("mensagem de erro inesperada: %s", respBody["error"])
	}
}

func TestUserController_DeleteProfile_404NotFound(t *testing.T) {
	targetID := uuid.New()

	svc := &mockUserService{
		DeleteProfileFn: func(requesterID uuid.UUID, targetID uuid.UUID) error {
			return utils.ErrUserNotFound
		},
	}

	router, _ := setupUserRouter(svc)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/user/profile/"+targetID.String(), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("esperado status 404, obteve: %d", w.Code)
	}

	var respBody map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &respBody); err != nil {
		t.Fatalf("erro ao decodificar resposta: %v", err)
	}
	if respBody["error"] != utils.ErrUserNotFound.Error() {
		t.Errorf("mensagem de erro inesperada: %s", respBody["error"])
	}
}

func TestUserController_DeleteProfile_500InternalServerError(t *testing.T) {
	targetID := uuid.New()

	svc := &mockUserService{
		DeleteProfileFn: func(requesterID uuid.UUID, targetID uuid.UUID) error {
			return errors.New("unexpected db error")
		},
	}

	router, _ := setupUserRouter(svc)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/user/profile/"+targetID.String(), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("esperado status 500, obteve: %d", w.Code)
	}

	var respBody map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &respBody); err != nil {
		t.Fatalf("erro ao decodificar resposta: %v", err)
	}
	if respBody["error"] != "internal server error" {
		t.Errorf("mensagem de erro inesperada: %s", respBody["error"])
	}
}
