package service_test

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
	"github.com/5gMurilo/helptrix-api/core/utils"
	servicemodule "github.com/5gMurilo/helptrix-api/modules/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// mockServiceService implementa IServiceService para testes do controller.
type mockServiceService struct {
	CreateFn   func(userID uuid.UUID, userType string, dto domain.CreateServiceRequestDTO) (domain.ServiceResponseDTO, error)
	ListFn     func(userID uuid.UUID, userType string) ([]domain.ServiceResponseDTO, error)
	GetByIDFn  func(serviceID uuid.UUID, userID uuid.UUID, userType string) (domain.ServiceResponseDTO, error)
	UpdateFn   func(serviceID uuid.UUID, userID uuid.UUID, userType string, dto domain.UpdateServiceRequestDTO) (domain.ServiceResponseDTO, error)
	DeleteFn   func(serviceID uuid.UUID, userID uuid.UUID, userType string) error
}

func (m *mockServiceService) Create(userID uuid.UUID, userType string, dto domain.CreateServiceRequestDTO) (domain.ServiceResponseDTO, error) {
	return m.CreateFn(userID, userType, dto)
}

func (m *mockServiceService) List(userID uuid.UUID, userType string) ([]domain.ServiceResponseDTO, error) {
	return m.ListFn(userID, userType)
}

func (m *mockServiceService) GetByID(serviceID uuid.UUID, userID uuid.UUID, userType string) (domain.ServiceResponseDTO, error) {
	return m.GetByIDFn(serviceID, userID, userType)
}

func (m *mockServiceService) Update(serviceID uuid.UUID, userID uuid.UUID, userType string, dto domain.UpdateServiceRequestDTO) (domain.ServiceResponseDTO, error) {
	return m.UpdateFn(serviceID, userID, userType, dto)
}

func (m *mockServiceService) Delete(serviceID uuid.UUID, userID uuid.UUID, userType string) error {
	return m.DeleteFn(serviceID, userID, userType)
}

func setupServiceRouter(svc *mockServiceService, userID string, userType string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	ctrl := servicemodule.NewServiceController(svc)

	injectPayload := func(c *gin.Context) {
		c.Set("authorization_payload", &auth.Payload{
			UserID:   userID,
			UserType: userType,
		})
	}

	r.POST("/service", func(c *gin.Context) {
		injectPayload(c)
		ctrl.Create(c)
	})
	r.GET("/service", func(c *gin.Context) {
		injectPayload(c)
		ctrl.List(c)
	})
	r.GET("/service/:id", func(c *gin.Context) {
		injectPayload(c)
		ctrl.GetByID(c)
	})
	r.PUT("/service/:id", func(c *gin.Context) {
		injectPayload(c)
		ctrl.Update(c)
	})
	r.DELETE("/service/:id", func(c *gin.Context) {
		injectPayload(c)
		ctrl.Delete(c)
	})
	return r
}

func defaultMockSvc() *mockServiceService {
	return &mockServiceService{
		CreateFn: func(userID uuid.UUID, userType string, dto domain.CreateServiceRequestDTO) (domain.ServiceResponseDTO, error) {
			return validServiceResponse(), nil
		},
		ListFn: func(userID uuid.UUID, userType string) ([]domain.ServiceResponseDTO, error) {
			return []domain.ServiceResponseDTO{validServiceResponse()}, nil
		},
		GetByIDFn: func(serviceID uuid.UUID, userID uuid.UUID, userType string) (domain.ServiceResponseDTO, error) {
			return validServiceResponse(), nil
		},
		UpdateFn: func(serviceID uuid.UUID, userID uuid.UUID, userType string, dto domain.UpdateServiceRequestDTO) (domain.ServiceResponseDTO, error) {
			return validServiceResponse(), nil
		},
		DeleteFn: func(serviceID uuid.UUID, userID uuid.UUID, userType string) error {
			return nil
		},
	}
}

func validCreateServiceBody() []byte {
	body := map[string]interface{}{
		"name":           "Serviço de Encanamento",
		"description":    "Conserto de canos e torneiras",
		"actuation_days": []string{"monday", "tuesday"},
		"value":          "150.00",
		"start_time":     "08:00",
		"end_time":       "18:00",
		"offer_since":    time.Now().Format(time.RFC3339),
		"category_id":    1,
	}
	b, _ := json.Marshal(body)
	return b
}

func validServiceResponse() domain.ServiceResponseDTO {
	return domain.ServiceResponseDTO{
		ID:            uuid.New(),
		Name:          "Serviço de Encanamento",
		Description:   "Conserto de canos e torneiras",
		ActuationDays: []string{"monday", "tuesday"},
		Value:         decimal.NewFromFloat(150.00),
		StartTime:     "08:00",
		EndTime:       "18:00",
		OfferSince:    time.Now(),
		Category: domain.ServiceCategoryDTO{
			ID:   1,
			Name: "Encanamento",
		},
	}
}

func TestServiceController_Create_201Created(t *testing.T) {
	mockResp := validServiceResponse()

	svc := &mockServiceService{
		CreateFn: func(userID uuid.UUID, userType string, dto domain.CreateServiceRequestDTO) (domain.ServiceResponseDTO, error) {
			return mockResp, nil
		},
	}

	router := setupServiceRouter(svc, uuid.New().String(), utils.UserTypeHelper)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/service", bytes.NewBuffer(validCreateServiceBody()))
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

func TestServiceController_Create_400BindError(t *testing.T) {
	svc := &mockServiceService{
		CreateFn: func(userID uuid.UUID, userType string, dto domain.CreateServiceRequestDTO) (domain.ServiceResponseDTO, error) {
			return domain.ServiceResponseDTO{}, nil
		},
	}

	router := setupServiceRouter(svc, uuid.New().String(), utils.UserTypeHelper)
	w := httptest.NewRecorder()

	// Corpo inválido: campos obrigatórios ausentes
	invalidBody := []byte(`{"name": "incompleto"}`)
	req, _ := http.NewRequest(http.MethodPost, "/service", bytes.NewBuffer(invalidBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado status 400, obteve: %d", w.Code)
	}
}

func TestServiceController_Create_400InvalidUUID(t *testing.T) {
	svc := &mockServiceService{
		CreateFn: func(userID uuid.UUID, userType string, dto domain.CreateServiceRequestDTO) (domain.ServiceResponseDTO, error) {
			return domain.ServiceResponseDTO{}, nil
		},
	}

	router := setupServiceRouter(svc, "invalid-uuid", utils.UserTypeHelper)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/service", bytes.NewBuffer(validCreateServiceBody()))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado status 400 para UUID inválido, obteve: %d", w.Code)
	}

	var respBody map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &respBody); err != nil {
		t.Fatalf("erro ao decodificar resposta: %v", err)
	}
	if respBody["error"] != "invalid user id" {
		t.Errorf("mensagem de erro inesperada: %s", respBody["error"])
	}
}

func TestServiceController_Create_403HelperOnly(t *testing.T) {
	svc := &mockServiceService{
		CreateFn: func(userID uuid.UUID, userType string, dto domain.CreateServiceRequestDTO) (domain.ServiceResponseDTO, error) {
			return domain.ServiceResponseDTO{}, utils.ErrHelperOnly
		},
	}

	router := setupServiceRouter(svc, uuid.New().String(), utils.UserTypeBusiness)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/service", bytes.NewBuffer(validCreateServiceBody()))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("esperado status 403, obteve: %d", w.Code)
	}

	var respBody map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &respBody); err != nil {
		t.Fatalf("erro ao decodificar resposta: %v", err)
	}
	if respBody["error"] != utils.ErrHelperOnly.Error() {
		t.Errorf("mensagem de erro inesperada: %s", respBody["error"])
	}
}

func TestServiceController_Create_409NameConflict(t *testing.T) {
	svc := &mockServiceService{
		CreateFn: func(userID uuid.UUID, userType string, dto domain.CreateServiceRequestDTO) (domain.ServiceResponseDTO, error) {
			return domain.ServiceResponseDTO{}, utils.ErrServiceNameNotUnique
		},
	}

	router := setupServiceRouter(svc, uuid.New().String(), utils.UserTypeHelper)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/service", bytes.NewBuffer(validCreateServiceBody()))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("esperado status 409, obteve: %d", w.Code)
	}

	var respBody map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &respBody); err != nil {
		t.Fatalf("erro ao decodificar resposta: %v", err)
	}
	if respBody["error"] != utils.ErrServiceNameNotUnique.Error() {
		t.Errorf("mensagem de erro inesperada: %s", respBody["error"])
	}
}

func TestServiceController_Create_422CategoryNotAssigned(t *testing.T) {
	svc := &mockServiceService{
		CreateFn: func(userID uuid.UUID, userType string, dto domain.CreateServiceRequestDTO) (domain.ServiceResponseDTO, error) {
			return domain.ServiceResponseDTO{}, utils.ErrCategoryNotAssignedToUser
		},
	}

	router := setupServiceRouter(svc, uuid.New().String(), utils.UserTypeHelper)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/service", bytes.NewBuffer(validCreateServiceBody()))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("esperado status 422, obteve: %d", w.Code)
	}

	var respBody map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &respBody); err != nil {
		t.Fatalf("erro ao decodificar resposta: %v", err)
	}
	if respBody["error"] != utils.ErrCategoryNotAssignedToUser.Error() {
		t.Errorf("mensagem de erro inesperada: %s", respBody["error"])
	}
}

func TestServiceController_Create_500InternalServerError(t *testing.T) {
	svc := &mockServiceService{
		CreateFn: func(userID uuid.UUID, userType string, dto domain.CreateServiceRequestDTO) (domain.ServiceResponseDTO, error) {
			return domain.ServiceResponseDTO{}, errors.New("unexpected db error")
		},
	}

	router := setupServiceRouter(svc, uuid.New().String(), utils.UserTypeHelper)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/service", bytes.NewBuffer(validCreateServiceBody()))
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

// ---- List tests ----

func TestServiceController_List_200OK(t *testing.T) {
	svc := defaultMockSvc()
	router := setupServiceRouter(svc, uuid.New().String(), utils.UserTypeHelper)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/service", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("esperado status 200, obteve: %d", w.Code)
	}

	var respBody []interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &respBody); err != nil {
		t.Fatalf("erro ao decodificar resposta: %v", err)
	}
	if len(respBody) == 0 {
		t.Error("esperado pelo menos um serviço na resposta")
	}
}

func TestServiceController_List_403HelperOnly(t *testing.T) {
	svc := defaultMockSvc()
	svc.ListFn = func(userID uuid.UUID, userType string) ([]domain.ServiceResponseDTO, error) {
		return nil, utils.ErrHelperOnly
	}

	router := setupServiceRouter(svc, uuid.New().String(), utils.UserTypeBusiness)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/service", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("esperado status 403, obteve: %d", w.Code)
	}
}

func TestServiceController_List_500InternalServerError(t *testing.T) {
	svc := defaultMockSvc()
	svc.ListFn = func(userID uuid.UUID, userType string) ([]domain.ServiceResponseDTO, error) {
		return nil, errors.New("db error")
	}

	router := setupServiceRouter(svc, uuid.New().String(), utils.UserTypeHelper)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/service", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("esperado status 500, obteve: %d", w.Code)
	}
}

// ---- GetByID tests ----

func TestServiceController_GetByID_200OK(t *testing.T) {
	svc := defaultMockSvc()
	router := setupServiceRouter(svc, uuid.New().String(), utils.UserTypeHelper)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/service/"+uuid.New().String(), nil)
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

func TestServiceController_GetByID_400InvalidUUID(t *testing.T) {
	svc := defaultMockSvc()
	router := setupServiceRouter(svc, uuid.New().String(), utils.UserTypeHelper)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/service/not-a-uuid", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado status 400, obteve: %d", w.Code)
	}
}

func TestServiceController_GetByID_403HelperOnly(t *testing.T) {
	svc := defaultMockSvc()
	svc.GetByIDFn = func(serviceID uuid.UUID, userID uuid.UUID, userType string) (domain.ServiceResponseDTO, error) {
		return domain.ServiceResponseDTO{}, utils.ErrHelperOnly
	}

	router := setupServiceRouter(svc, uuid.New().String(), utils.UserTypeBusiness)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/service/"+uuid.New().String(), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("esperado status 403, obteve: %d", w.Code)
	}
}

func TestServiceController_GetByID_404NotFound(t *testing.T) {
	svc := defaultMockSvc()
	svc.GetByIDFn = func(serviceID uuid.UUID, userID uuid.UUID, userType string) (domain.ServiceResponseDTO, error) {
		return domain.ServiceResponseDTO{}, utils.ErrServiceNotFound
	}

	router := setupServiceRouter(svc, uuid.New().String(), utils.UserTypeHelper)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/service/"+uuid.New().String(), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("esperado status 404, obteve: %d", w.Code)
	}
}

func TestServiceController_GetByID_500InternalServerError(t *testing.T) {
	svc := defaultMockSvc()
	svc.GetByIDFn = func(serviceID uuid.UUID, userID uuid.UUID, userType string) (domain.ServiceResponseDTO, error) {
		return domain.ServiceResponseDTO{}, errors.New("db error")
	}

	router := setupServiceRouter(svc, uuid.New().String(), utils.UserTypeHelper)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/service/"+uuid.New().String(), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("esperado status 500, obteve: %d", w.Code)
	}
}

// ---- Update tests ----

func validUpdateServiceBody() []byte {
	name := "Serviço Atualizado"
	body := map[string]interface{}{
		"name": name,
	}
	b, _ := json.Marshal(body)
	return b
}

func TestServiceController_Update_200OK(t *testing.T) {
	svc := defaultMockSvc()
	router := setupServiceRouter(svc, uuid.New().String(), utils.UserTypeHelper)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/service/"+uuid.New().String(), bytes.NewBuffer(validUpdateServiceBody()))
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
}

func TestServiceController_Update_400InvalidUUID(t *testing.T) {
	svc := defaultMockSvc()
	router := setupServiceRouter(svc, uuid.New().String(), utils.UserTypeHelper)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/service/not-a-uuid", bytes.NewBuffer(validUpdateServiceBody()))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado status 400, obteve: %d", w.Code)
	}
}

func TestServiceController_Update_403HelperOnly(t *testing.T) {
	svc := defaultMockSvc()
	svc.UpdateFn = func(serviceID uuid.UUID, userID uuid.UUID, userType string, dto domain.UpdateServiceRequestDTO) (domain.ServiceResponseDTO, error) {
		return domain.ServiceResponseDTO{}, utils.ErrHelperOnly
	}

	router := setupServiceRouter(svc, uuid.New().String(), utils.UserTypeBusiness)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/service/"+uuid.New().String(), bytes.NewBuffer(validUpdateServiceBody()))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("esperado status 403, obteve: %d", w.Code)
	}
}

func TestServiceController_Update_404NotFound(t *testing.T) {
	svc := defaultMockSvc()
	svc.UpdateFn = func(serviceID uuid.UUID, userID uuid.UUID, userType string, dto domain.UpdateServiceRequestDTO) (domain.ServiceResponseDTO, error) {
		return domain.ServiceResponseDTO{}, utils.ErrServiceNotFound
	}

	router := setupServiceRouter(svc, uuid.New().String(), utils.UserTypeHelper)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/service/"+uuid.New().String(), bytes.NewBuffer(validUpdateServiceBody()))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("esperado status 404, obteve: %d", w.Code)
	}
}

func TestServiceController_Update_409NameConflict(t *testing.T) {
	svc := defaultMockSvc()
	svc.UpdateFn = func(serviceID uuid.UUID, userID uuid.UUID, userType string, dto domain.UpdateServiceRequestDTO) (domain.ServiceResponseDTO, error) {
		return domain.ServiceResponseDTO{}, utils.ErrServiceNameNotUnique
	}

	router := setupServiceRouter(svc, uuid.New().String(), utils.UserTypeHelper)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/service/"+uuid.New().String(), bytes.NewBuffer(validUpdateServiceBody()))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("esperado status 409, obteve: %d", w.Code)
	}
}

func TestServiceController_Update_422CategoryNotAssigned(t *testing.T) {
	svc := defaultMockSvc()
	svc.UpdateFn = func(serviceID uuid.UUID, userID uuid.UUID, userType string, dto domain.UpdateServiceRequestDTO) (domain.ServiceResponseDTO, error) {
		return domain.ServiceResponseDTO{}, utils.ErrCategoryNotAssignedToUser
	}

	router := setupServiceRouter(svc, uuid.New().String(), utils.UserTypeHelper)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/service/"+uuid.New().String(), bytes.NewBuffer(validUpdateServiceBody()))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("esperado status 422, obteve: %d", w.Code)
	}
}

func TestServiceController_Update_500InternalServerError(t *testing.T) {
	svc := defaultMockSvc()
	svc.UpdateFn = func(serviceID uuid.UUID, userID uuid.UUID, userType string, dto domain.UpdateServiceRequestDTO) (domain.ServiceResponseDTO, error) {
		return domain.ServiceResponseDTO{}, errors.New("db error")
	}

	router := setupServiceRouter(svc, uuid.New().String(), utils.UserTypeHelper)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/service/"+uuid.New().String(), bytes.NewBuffer(validUpdateServiceBody()))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("esperado status 500, obteve: %d", w.Code)
	}
}

// ---- Delete tests ----

func TestServiceController_Delete_204NoContent(t *testing.T) {
	svc := defaultMockSvc()
	router := setupServiceRouter(svc, uuid.New().String(), utils.UserTypeHelper)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/service/"+uuid.New().String(), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("esperado status 204, obteve: %d", w.Code)
	}
}

func TestServiceController_Delete_400InvalidUUID(t *testing.T) {
	svc := defaultMockSvc()
	router := setupServiceRouter(svc, uuid.New().String(), utils.UserTypeHelper)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/service/not-a-uuid", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado status 400, obteve: %d", w.Code)
	}
}

func TestServiceController_Delete_403HelperOnly(t *testing.T) {
	svc := defaultMockSvc()
	svc.DeleteFn = func(serviceID uuid.UUID, userID uuid.UUID, userType string) error {
		return utils.ErrHelperOnly
	}

	router := setupServiceRouter(svc, uuid.New().String(), utils.UserTypeBusiness)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/service/"+uuid.New().String(), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("esperado status 403, obteve: %d", w.Code)
	}
}

func TestServiceController_Delete_404NotFound(t *testing.T) {
	svc := defaultMockSvc()
	svc.DeleteFn = func(serviceID uuid.UUID, userID uuid.UUID, userType string) error {
		return utils.ErrServiceNotFound
	}

	router := setupServiceRouter(svc, uuid.New().String(), utils.UserTypeHelper)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/service/"+uuid.New().String(), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("esperado status 404, obteve: %d", w.Code)
	}
}

func TestServiceController_Delete_500InternalServerError(t *testing.T) {
	svc := defaultMockSvc()
	svc.DeleteFn = func(serviceID uuid.UUID, userID uuid.UUID, userType string) error {
		return errors.New("db error")
	}

	router := setupServiceRouter(svc, uuid.New().String(), utils.UserTypeHelper)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/service/"+uuid.New().String(), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("esperado status 500, obteve: %d", w.Code)
	}
}
