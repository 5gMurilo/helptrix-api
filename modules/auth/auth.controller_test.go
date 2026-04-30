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

func validRegisterBody(t *testing.T) []byte {
	t.Helper()
	body := map[string]interface{}{
		"name":       "John Silva",
		"email":      "john@example.com",
		"password":   "password123",
		"user_type":  utils.UserTypeHelper,
		"document":   "12345678901",
		"phone":      "11999999999",
		"categories": []uint{1, 2},
		"address": map[string]interface{}{
			"street":       "Flower Street",
			"number":       "100",
			"neighborhood": "Downtown",
			"city":         "Sao Paulo",
			"state":        "SP",
			"zip_code":     "01001000",
		},
	}
	b, _ := json.Marshal(body)
	return b
}

func validLoginBody(t *testing.T) []byte {
	t.Helper()
	body := map[string]interface{}{
		"email":    "john@example.com",
		"password": "password123",
	}
	b, _ := json.Marshal(body)
	return b
}

func setupRouter(t *testing.T, svc *mockAuthService) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	ctrl := authmodule.NewAuthController(svc)
	r.POST("/auth/register", ctrl.Register)
	r.POST("/auth/login", ctrl.Login)
	return r
}

func TestAuthController_Register(t *testing.T) {
	t.Run("201: valid registration", func(t *testing.T) {
		// Arrange
		mockResp := domain.RegisterResponseDTO{
			ID:         uuid.New(),
			Name:       "John Silva",
			Email:      "john@example.com",
			UserType:   utils.UserTypeHelper,
			Categories: []uint{1, 2},
			CreatedAt:  time.Now(),
		}
		svc := &mockAuthService{
			RegisterFn: func(dto domain.RegisterRequestDTO) (domain.RegisterResponseDTO, error) {
				return mockResp, nil
			},
		}
		router := setupRouter(t, svc)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(validRegisterBody(t)))
		req.Header.Set("Content-Type", "application/json")

		// Act
		router.ServeHTTP(w, req)

		// Assert
		if w.Code != http.StatusCreated {
			t.Fatalf("expected status 201, got: %d", w.Code)
		}
		var respBody map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &respBody); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if _, ok := respBody["id"]; !ok {
			t.Error("response must contain key 'id'")
		}
	})

	t.Run("409: duplicate user", func(t *testing.T) {
		// Arrange
		svc := &mockAuthService{
			RegisterFn: func(dto domain.RegisterRequestDTO) (domain.RegisterResponseDTO, error) {
				return domain.RegisterResponseDTO{}, utils.ErrUserAlreadyRegistered
			},
		}
		router := setupRouter(t, svc)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(validRegisterBody(t)))
		req.Header.Set("Content-Type", "application/json")

		// Act
		router.ServeHTTP(w, req)

		// Assert
		if w.Code != http.StatusConflict {
			t.Fatalf("expected status 409, got: %d", w.Code)
		}
		var respBody map[string]string
		if err := json.Unmarshal(w.Body.Bytes(), &respBody); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if respBody["error"] != utils.ErrUserAlreadyRegistered.Error() {
			t.Errorf("unexpected error message: %s", respBody["error"])
		}
	})

	t.Run("500: internal server error", func(t *testing.T) {
		// Arrange
		svc := &mockAuthService{
			RegisterFn: func(dto domain.RegisterRequestDTO) (domain.RegisterResponseDTO, error) {
				return domain.RegisterResponseDTO{}, errors.New("unexpected db error")
			},
		}
		router := setupRouter(t, svc)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(validRegisterBody(t)))
		req.Header.Set("Content-Type", "application/json")

		// Act
		router.ServeHTTP(w, req)

		// Assert
		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected status 500, got: %d", w.Code)
		}
		var respBody map[string]string
		if err := json.Unmarshal(w.Body.Bytes(), &respBody); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if respBody["error"] != "internal server error" {
			t.Errorf("unexpected error message: %s", respBody["error"])
		}
	})

	t.Run("400: invalid request body", func(t *testing.T) {
		// Arrange
		svc := &mockAuthService{
			RegisterFn: func(dto domain.RegisterRequestDTO) (domain.RegisterResponseDTO, error) {
				return domain.RegisterResponseDTO{}, nil
			},
		}
		router := setupRouter(t, svc)
		w := httptest.NewRecorder()
		invalidBody := []byte(`{"email": "invalid"}`)
		req, _ := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(invalidBody))
		req.Header.Set("Content-Type", "application/json")

		// Act
		router.ServeHTTP(w, req)

		// Assert
		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400, got: %d", w.Code)
		}
	})
}

func TestAuthController_Login(t *testing.T) {
	t.Run("201: valid credentials", func(t *testing.T) {
		// Arrange
		svc := &mockAuthService{
			LoginFn: func(dto domain.LoginRequestDTO) (domain.LoginResponseDTO, error) {
				return domain.LoginResponseDTO{Token: "tok"}, nil
			},
		}
		router := setupRouter(t, svc)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(validLoginBody(t)))
		req.Header.Set("Content-Type", "application/json")

		// Act
		router.ServeHTTP(w, req)

		// Assert
		if w.Code != http.StatusCreated {
			t.Fatalf("expected status 201, got: %d", w.Code)
		}
		var respBody map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &respBody); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if token, ok := respBody["token"]; !ok || token != "tok" {
			t.Errorf("expected token 'tok', got: %v", respBody["token"])
		}
	})

	t.Run("401: invalid credentials", func(t *testing.T) {
		// Arrange
		svc := &mockAuthService{
			LoginFn: func(dto domain.LoginRequestDTO) (domain.LoginResponseDTO, error) {
				return domain.LoginResponseDTO{}, utils.ErrInvalidCredentials
			},
		}
		router := setupRouter(t, svc)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(validLoginBody(t)))
		req.Header.Set("Content-Type", "application/json")

		// Act
		router.ServeHTTP(w, req)

		// Assert
		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected status 401, got: %d", w.Code)
		}
		var respBody map[string]string
		if err := json.Unmarshal(w.Body.Bytes(), &respBody); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if _, ok := respBody["error"]; !ok {
			t.Error("response must contain key 'error'")
		}
	})

	t.Run("500: internal server error", func(t *testing.T) {
		// Arrange
		svc := &mockAuthService{
			LoginFn: func(dto domain.LoginRequestDTO) (domain.LoginResponseDTO, error) {
				return domain.LoginResponseDTO{}, errors.New("db crashed")
			},
		}
		router := setupRouter(t, svc)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(validLoginBody(t)))
		req.Header.Set("Content-Type", "application/json")

		// Act
		router.ServeHTTP(w, req)

		// Assert
		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected status 500, got: %d", w.Code)
		}
		var respBody map[string]string
		if err := json.Unmarshal(w.Body.Bytes(), &respBody); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if respBody["error"] != "internal server error" {
			t.Errorf("unexpected error message: %s", respBody["error"])
		}
	})

	t.Run("400: invalid request body", func(t *testing.T) {
		// Arrange
		svc := &mockAuthService{}
		router := setupRouter(t, svc)
		w := httptest.NewRecorder()
		invalidBody := []byte(`{"email": "john@example.com"}`)
		req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(invalidBody))
		req.Header.Set("Content-Type", "application/json")

		// Act
		router.ServeHTTP(w, req)

		// Assert
		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400, got: %d", w.Code)
		}
	})
}
