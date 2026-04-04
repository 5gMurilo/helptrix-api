package otpmodule_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/5gMurilo/helptrix-api/core/domain"
	"github.com/5gMurilo/helptrix-api/core/utils"
	otpmodule "github.com/5gMurilo/helptrix-api/modules/otp"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// mockOtpService implementa IOtpService para testes do controller.
type mockOtpService struct {
	SendFn    func(dto domain.SendOTPRequestDTO) (domain.SendOTPResponseDTO, error)
	ConfirmFn func(dto domain.ConfirmOTPRequestDTO) (domain.ConfirmOTPResponseDTO, error)
}

func (m *mockOtpService) Send(dto domain.SendOTPRequestDTO) (domain.SendOTPResponseDTO, error) {
	return m.SendFn(dto)
}

func (m *mockOtpService) Confirm(dto domain.ConfirmOTPRequestDTO) (domain.ConfirmOTPResponseDTO, error) {
	return m.ConfirmFn(dto)
}

func setupOtpRouter(svc *mockOtpService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	ctrl := otpmodule.NewOtpController(svc)
	r.POST("/otp/send", ctrl.Send)
	r.POST("/otp/confirm", ctrl.Confirm)
	return r
}

func validSendOTPBody() []byte {
	body := map[string]interface{}{
		"email": "user@example.com",
	}
	b, _ := json.Marshal(body)
	return b
}

func validConfirmOTPBody(id string) []byte {
	body := map[string]interface{}{
		"id":   id,
		"code": "1234",
	}
	b, _ := json.Marshal(body)
	return b
}

func TestOtpController_Send_200OK(t *testing.T) {
	mockResp := domain.SendOTPResponseDTO{
		ID:      uuid.New(),
		Message: "OTP sent successfully",
	}

	svc := &mockOtpService{
		SendFn: func(dto domain.SendOTPRequestDTO) (domain.SendOTPResponseDTO, error) {
			return mockResp, nil
		},
	}

	router := setupOtpRouter(svc)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/otp/send", bytes.NewBuffer(validSendOTPBody()))
	req.Header.Set("Content-Type", "application/json")
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
	if respBody["message"] != "OTP sent successfully" {
		t.Errorf("mensagem inesperada: %v", respBody["message"])
	}
}

func TestOtpController_Send_400BadRequest(t *testing.T) {
	svc := &mockOtpService{}

	router := setupOtpRouter(svc)
	w := httptest.NewRecorder()

	invalidBody := []byte(`{"not_email": "invalid"}`)
	req, _ := http.NewRequest(http.MethodPost, "/otp/send", bytes.NewBuffer(invalidBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado status 400, obteve: %d", w.Code)
	}

	var respBody map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &respBody); err != nil {
		t.Fatalf("erro ao decodificar resposta: %v", err)
	}
	if respBody["error"] != "invalid request" {
		t.Errorf("mensagem de erro inesperada: %s", respBody["error"])
	}
}

func TestOtpController_Send_500InternalServerError(t *testing.T) {
	svc := &mockOtpService{
		SendFn: func(dto domain.SendOTPRequestDTO) (domain.SendOTPResponseDTO, error) {
			return domain.SendOTPResponseDTO{}, errors.New("error creating otp")
		},
	}

	router := setupOtpRouter(svc)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/otp/send", bytes.NewBuffer(validSendOTPBody()))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("esperado status 500, obteve: %d", w.Code)
	}

	var respBody map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &respBody); err != nil {
		t.Fatalf("erro ao decodificar resposta: %v", err)
	}
	if _, ok := respBody["error"]; !ok {
		t.Error("resposta deve conter a chave 'error'")
	}
}

func TestOtpController_Confirm_200OK(t *testing.T) {
	otpID := uuid.New()

	svc := &mockOtpService{
		ConfirmFn: func(dto domain.ConfirmOTPRequestDTO) (domain.ConfirmOTPResponseDTO, error) {
			return domain.ConfirmOTPResponseDTO{Message: "OTP confirmed successfully"}, nil
		},
	}

	router := setupOtpRouter(svc)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/otp/confirm", bytes.NewBuffer(validConfirmOTPBody(otpID.String())))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("esperado status 200, obteve: %d", w.Code)
	}

	var respBody map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &respBody); err != nil {
		t.Fatalf("erro ao decodificar resposta: %v", err)
	}
	if respBody["message"] != "OTP confirmed successfully" {
		t.Errorf("mensagem inesperada: %s", respBody["message"])
	}
}

func TestOtpController_Confirm_400BadRequest(t *testing.T) {
	svc := &mockOtpService{}

	router := setupOtpRouter(svc)
	w := httptest.NewRecorder()

	// id ausente e code ausente
	invalidBody := []byte(`{"not_valid": true}`)
	req, _ := http.NewRequest(http.MethodPost, "/otp/confirm", bytes.NewBuffer(invalidBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado status 400, obteve: %d", w.Code)
	}

	var respBody map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &respBody); err != nil {
		t.Fatalf("erro ao decodificar resposta: %v", err)
	}
	if respBody["error"] != "invalid request" {
		t.Errorf("mensagem de erro inesperada: %s", respBody["error"])
	}
}

func TestOtpController_Confirm_404NotFound(t *testing.T) {
	otpID := uuid.New()

	svc := &mockOtpService{
		ConfirmFn: func(dto domain.ConfirmOTPRequestDTO) (domain.ConfirmOTPResponseDTO, error) {
			return domain.ConfirmOTPResponseDTO{}, utils.ErrOTPNotFound
		},
	}

	router := setupOtpRouter(svc)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/otp/confirm", bytes.NewBuffer(validConfirmOTPBody(otpID.String())))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("esperado status 404, obteve: %d", w.Code)
	}

	var respBody map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &respBody); err != nil {
		t.Fatalf("erro ao decodificar resposta: %v", err)
	}
	if _, ok := respBody["error"]; !ok {
		t.Error("resposta deve conter a chave 'error'")
	}
}

func TestOtpController_Confirm_422ErrOTPNotWaiting(t *testing.T) {
	otpID := uuid.New()

	svc := &mockOtpService{
		ConfirmFn: func(dto domain.ConfirmOTPRequestDTO) (domain.ConfirmOTPResponseDTO, error) {
			return domain.ConfirmOTPResponseDTO{}, utils.ErrOTPNotWaiting
		},
	}

	router := setupOtpRouter(svc)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/otp/confirm", bytes.NewBuffer(validConfirmOTPBody(otpID.String())))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("esperado status 422, obteve: %d", w.Code)
	}

	var respBody map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &respBody); err != nil {
		t.Fatalf("erro ao decodificar resposta: %v", err)
	}
	if respBody["error"] != utils.ErrOTPNotWaiting.Error() {
		t.Errorf("mensagem de erro inesperada: %s", respBody["error"])
	}
}

func TestOtpController_Confirm_422ErrOTPExpired(t *testing.T) {
	otpID := uuid.New()

	svc := &mockOtpService{
		ConfirmFn: func(dto domain.ConfirmOTPRequestDTO) (domain.ConfirmOTPResponseDTO, error) {
			return domain.ConfirmOTPResponseDTO{}, utils.ErrOTPExpired
		},
	}

	router := setupOtpRouter(svc)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/otp/confirm", bytes.NewBuffer(validConfirmOTPBody(otpID.String())))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("esperado status 422, obteve: %d", w.Code)
	}

	var respBody map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &respBody); err != nil {
		t.Fatalf("erro ao decodificar resposta: %v", err)
	}
	if respBody["error"] != utils.ErrOTPExpired.Error() {
		t.Errorf("mensagem de erro inesperada: %s", respBody["error"])
	}
}

func TestOtpController_Confirm_422ErrOTPInvalid(t *testing.T) {
	otpID := uuid.New()

	svc := &mockOtpService{
		ConfirmFn: func(dto domain.ConfirmOTPRequestDTO) (domain.ConfirmOTPResponseDTO, error) {
			return domain.ConfirmOTPResponseDTO{}, utils.ErrOTPInvalid
		},
	}

	router := setupOtpRouter(svc)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/otp/confirm", bytes.NewBuffer(validConfirmOTPBody(otpID.String())))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("esperado status 422, obteve: %d", w.Code)
	}

	var respBody map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &respBody); err != nil {
		t.Fatalf("erro ao decodificar resposta: %v", err)
	}
	if respBody["error"] != utils.ErrOTPInvalid.Error() {
		t.Errorf("mensagem de erro inesperada: %s", respBody["error"])
	}
}

func TestOtpController_Confirm_500InternalServerError(t *testing.T) {
	otpID := uuid.New()

	svc := &mockOtpService{
		ConfirmFn: func(dto domain.ConfirmOTPRequestDTO) (domain.ConfirmOTPResponseDTO, error) {
			return domain.ConfirmOTPResponseDTO{}, errors.New("error confirming otp")
		},
	}

	router := setupOtpRouter(svc)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/otp/confirm", bytes.NewBuffer(validConfirmOTPBody(otpID.String())))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("esperado status 500, obteve: %d", w.Code)
	}

	var respBody map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &respBody); err != nil {
		t.Fatalf("erro ao decodificar resposta: %v", err)
	}
	if _, ok := respBody["error"]; !ok {
		t.Error("resposta deve conter a chave 'error'")
	}
}
