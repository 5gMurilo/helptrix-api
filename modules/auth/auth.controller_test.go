package auth_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/5gMurilo/helptrix-api/core/domain"
	"github.com/5gMurilo/helptrix-api/core/utils"
	authmodule "github.com/5gMurilo/helptrix-api/modules/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// mockAuthService implementa IAuthService para testes do controller.
type mockAuthService struct {
	RegisterFn func(dto domain.RegisterRequestDTO) (domain.RegisterResponseDTO, error)
	LoginFn    func(dto domain.LoginRequestDTO) (domain.LoginResponseDTO, error)
}

func (m *mockAuthService) Register(dto domain.RegisterRequestDTO) (domain.RegisterResponseDTO, error) {
	return m.RegisterFn(dto)
}

func (m *mockAuthService) Login(dto domain.LoginRequestDTO) (domain.LoginResponseDTO, error) {
	return m.LoginFn(dto)
}

func validRegisterBody() []byte {
	body := map[string]interface{}{
		"name":      "João da Silva",
		"email":     "joao@example.com",
		"password":  "senha123",
		"user_type": utils.UserTypeHelper,
		"document":  "12345678901",
		"phone":     "11999999999",
		"categories": []uint{1, 2},
		"address": map[string]interface{}{
			"street":       "Rua das Flores",
			"number":       "100",
			"neighborhood": "Centro",
			"city":         "São Paulo",
			"state":        "SP",
		},
	}
	b, _ := json.Marshal(body)
	return b
}

func setupRouter(svc *mockAuthService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	ctrl := authmodule.NewAuthController(svc)
	r.POST("/auth/register", ctrl.Register)
	r.POST("/auth/login", ctrl.Login)
	return r
}

func validLoginBody() []byte {
	body := map[string]interface{}{
		"email":    "joao@example.com",
		"password": "senha123",
	}
	b, _ := json.Marshal(body)
	return b
}

func TestAuthController_Register_201Created(t *testing.T) {
	mockResp := domain.RegisterResponseDTO{
		ID:         uuid.New(),
		Name:       "João da Silva",
		Email:      "joao@example.com",
		UserType:   utils.UserTypeHelper,
		Categories: []uint{1, 2},
		CreatedAt:  time.Now(),
	}

	svc := &mockAuthService{
		RegisterFn: func(dto domain.RegisterRequestDTO) (domain.RegisterResponseDTO, error) {
			return mockResp, nil
		},
	}

	router := setupRouter(svc)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(validRegisterBody()))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("esperado status 201, obteve: %d", w.Code)
	}

	var respBody map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &respBody); err != nil {
		t.Fatalf("erro ao decodificar resposta: %v", err)
	}
	if _, ok := respBody["id"]; !ok {
		t.Error("resposta deve conter a chave 'id'")
	}
}

func TestAuthController_Register_409Conflict(t *testing.T) {
	svc := &mockAuthService{
		RegisterFn: func(dto domain.RegisterRequestDTO) (domain.RegisterResponseDTO, error) {
			return domain.RegisterResponseDTO{}, utils.ErrUserAlreadyRegistered
		},
	}

	router := setupRouter(svc)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(validRegisterBody()))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("esperado status 409, obteve: %d", w.Code)
	}

	var respBody map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &respBody); err != nil {
		t.Fatalf("erro ao decodificar resposta: %v", err)
	}
	if respBody["error"] != utils.ErrUserAlreadyRegistered.Error() {
		t.Errorf("mensagem de erro inesperada: %s", respBody["error"])
	}
}

func TestAuthController_Register_500InternalServerError(t *testing.T) {
	svc := &mockAuthService{
		RegisterFn: func(dto domain.RegisterRequestDTO) (domain.RegisterResponseDTO, error) {
			return domain.RegisterResponseDTO{}, errors.New("unexpected db error")
		},
	}

	router := setupRouter(svc)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(validRegisterBody()))
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

func TestAuthController_Register_400BadRequest(t *testing.T) {
	svc := &mockAuthService{
		RegisterFn: func(dto domain.RegisterRequestDTO) (domain.RegisterResponseDTO, error) {
			return domain.RegisterResponseDTO{}, nil
		},
	}

	router := setupRouter(svc)
	w := httptest.NewRecorder()

	// Corpo inválido: campos obrigatórios ausentes
	invalidBody := []byte(`{"email": "invalido"}`)
	req, _ := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(invalidBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado status 400, obteve: %d", w.Code)
	}
}

func TestAuthController_Login_201Created(t *testing.T) {
	svc := &mockAuthService{
		LoginFn: func(dto domain.LoginRequestDTO) (domain.LoginResponseDTO, error) {
			return domain.LoginResponseDTO{Token: "tok"}, nil
		},
	}

	router := setupRouter(svc)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(validLoginBody()))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("esperado status 201, obteve: %d", w.Code)
	}

	var respBody map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &respBody); err != nil {
		t.Fatalf("erro ao decodificar resposta: %v", err)
	}
	if token, ok := respBody["token"]; !ok || token != "tok" {
		t.Errorf("esperado token 'tok', obteve: %v", respBody["token"])
	}
}

func TestAuthController_Login_401Unauthorized(t *testing.T) {
	svc := &mockAuthService{
		LoginFn: func(dto domain.LoginRequestDTO) (domain.LoginResponseDTO, error) {
			return domain.LoginResponseDTO{}, utils.ErrInvalidCredentials
		},
	}

	router := setupRouter(svc)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(validLoginBody()))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("esperado status 401, obteve: %d", w.Code)
	}

	var respBody map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &respBody); err != nil {
		t.Fatalf("erro ao decodificar resposta: %v", err)
	}
	if _, ok := respBody["error"]; !ok {
		t.Error("resposta deve conter a chave 'error'")
	}
}

func TestAuthController_Login_400BadRequest(t *testing.T) {
	svc := &mockAuthService{}

	router := setupRouter(svc)
	w := httptest.NewRecorder()

	// Corpo inválido: campo password ausente
	invalidBody := []byte(`{"email": "joao@example.com"}`)
	req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(invalidBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado status 400, obteve: %d", w.Code)
	}
}
