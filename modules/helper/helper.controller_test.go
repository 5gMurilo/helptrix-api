package helper_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/5gMurilo/helptrix-api/adapter/auth"
	"github.com/5gMurilo/helptrix-api/core/domain"
	"github.com/5gMurilo/helptrix-api/core/utils"
	helpermodule "github.com/5gMurilo/helptrix-api/modules/helper"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type mockHelperService struct {
	searchFn func(string, domain.HelperSearchParams) (domain.HelperListResponseDTO, error)
	lastReq  string
	lastP    domain.HelperSearchParams
}

func (m *mockHelperService) Search(requesterType string, params domain.HelperSearchParams) (domain.HelperListResponseDTO, error) {
	m.lastReq = requesterType
	m.lastP = params
	if m.searchFn != nil {
		return m.searchFn(requesterType, params)
	}
	return domain.HelperListResponseDTO{}, nil
}

func TestList_BusinessOnlyError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &mockHelperService{
		searchFn: func(string, domain.HelperSearchParams) (domain.HelperListResponseDTO, error) {
			return domain.HelperListResponseDTO{}, utils.ErrBusinessOnly
		},
	}
	ctrl := helpermodule.NewHelperController(svc)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/helper", nil)
	c.Set("authorization_payload", &auth.Payload{UserType: utils.UserTypeBusiness})
	ctrl.List(c)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}
	var body map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if body["error"] != utils.ErrBusinessOnly.Error() {
		t.Fatalf("unexpected body: %v", body)
	}
}

func TestList_InternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &mockHelperService{
		searchFn: func(string, domain.HelperSearchParams) (domain.HelperListResponseDTO, error) {
			return domain.HelperListResponseDTO{}, errors.New("db failure")
		},
	}
	ctrl := helpermodule.NewHelperController(svc)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/helper", nil)
	c.Set("authorization_payload", &auth.Payload{UserType: utils.UserTypeBusiness})
	ctrl.List(c)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
	var body map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if body["error"] != "internal server error" {
		t.Fatalf("unexpected body: %v", body)
	}
}

func TestList_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	id := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	svc := &mockHelperService{
		searchFn: func(string, domain.HelperSearchParams) (domain.HelperListResponseDTO, error) {
			return domain.HelperListResponseDTO{
				Data: []domain.HelperCardDTO{
					{ID: id, Name: "João", Categories: []domain.HelperCategoryDTO{{ID: 1, Name: "X"}}},
				},
				Total:    1,
				Page:     1,
				PageSize: 20,
			}, nil
		},
	}
	ctrl := helpermodule.NewHelperController(svc)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/helper", nil)
	c.Set("authorization_payload", &auth.Payload{UserType: utils.UserTypeBusiness})
	ctrl.List(c)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp domain.HelperListResponseDTO
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if len(resp.Data) != 1 || resp.Data[0].Name != "João" {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestList_DefaultPagination(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &mockHelperService{
		searchFn: func(string, domain.HelperSearchParams) (domain.HelperListResponseDTO, error) {
			return domain.HelperListResponseDTO{Data: []domain.HelperCardDTO{}}, nil
		},
	}
	ctrl := helpermodule.NewHelperController(svc)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/helper", nil)
	c.Set("authorization_payload", &auth.Payload{UserType: utils.UserTypeBusiness})
	ctrl.List(c)
	if svc.lastP.Page != 1 || svc.lastP.PageSize != 20 {
		t.Fatalf("expected Page=1 PageSize=20, got %+v", svc.lastP)
	}
	_ = w
}

func TestList_InvalidPageFallsToDefault(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &mockHelperService{
		searchFn: func(string, domain.HelperSearchParams) (domain.HelperListResponseDTO, error) {
			return domain.HelperListResponseDTO{Data: []domain.HelperCardDTO{}}, nil
		},
	}
	ctrl := helpermodule.NewHelperController(svc)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/helper?page=abc", nil)
	c.Set("authorization_payload", &auth.Payload{UserType: utils.UserTypeBusiness})
	ctrl.List(c)
	if svc.lastP.Page != 1 {
		t.Fatalf("expected Page=1, got %d", svc.lastP.Page)
	}
	_ = w
}

func TestList_InvalidCategoryIDIgnored(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &mockHelperService{
		searchFn: func(string, domain.HelperSearchParams) (domain.HelperListResponseDTO, error) {
			return domain.HelperListResponseDTO{Data: []domain.HelperCardDTO{}}, nil
		},
	}
	ctrl := helpermodule.NewHelperController(svc)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/helper?category_id=not-a-number", nil)
	c.Set("authorization_payload", &auth.Payload{UserType: utils.UserTypeBusiness})
	ctrl.List(c)
	if svc.lastP.CategoryID != nil {
		t.Fatalf("expected CategoryID nil, got %v", svc.lastP.CategoryID)
	}
	_ = w
}

func TestList_ValidCategoryID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &mockHelperService{
		searchFn: func(string, domain.HelperSearchParams) (domain.HelperListResponseDTO, error) {
			return domain.HelperListResponseDTO{Data: []domain.HelperCardDTO{}}, nil
		},
	}
	ctrl := helpermodule.NewHelperController(svc)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/helper?category_id=3", nil)
	c.Set("authorization_payload", &auth.Payload{UserType: utils.UserTypeBusiness})
	ctrl.List(c)
	if svc.lastP.CategoryID == nil || *svc.lastP.CategoryID != 3 {
		t.Fatalf("expected CategoryID 3, got %v", svc.lastP.CategoryID)
	}
	_ = w
}

func TestList_EmptyResult(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &mockHelperService{
		searchFn: func(string, domain.HelperSearchParams) (domain.HelperListResponseDTO, error) {
			return domain.HelperListResponseDTO{
				Data:     []domain.HelperCardDTO{},
				Total:    0,
				Page:     1,
				PageSize: 20,
			}, nil
		},
	}
	ctrl := helpermodule.NewHelperController(svc)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/helper", nil)
	c.Set("authorization_payload", &auth.Payload{UserType: utils.UserTypeBusiness})
	ctrl.List(c)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	raw := w.Body.String()
	if !json.Valid([]byte(raw)) {
		t.Fatalf("invalid json: %s", raw)
	}
	var m map[string]json.RawMessage
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	var data []json.RawMessage
	if err := json.Unmarshal(m["data"], &data); err != nil {
		t.Fatalf("data must be array, got %s", string(m["data"]))
	}
	if len(data) != 0 {
		t.Fatalf("expected empty data array")
	}
}
