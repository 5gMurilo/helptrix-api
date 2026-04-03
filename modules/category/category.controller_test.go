package category_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/5gMurilo/helptrix-api/core/domain"
	categoryinterfaces "github.com/5gMurilo/helptrix-api/core/interfaces/category"
	categorymodule "github.com/5gMurilo/helptrix-api/modules/category"
	"github.com/gin-gonic/gin"
)

type mockCategoryService struct {
	ListFn func() ([]domain.CategoryListItemResponseDTO, error)
}

func (m *mockCategoryService) List() ([]domain.CategoryListItemResponseDTO, error) {
	return m.ListFn()
}

var _ categoryinterfaces.ICategoryService = (*mockCategoryService)(nil)

func setupCategoryRouter(svc *mockCategoryService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	ctrl := categorymodule.NewCategoryController(svc)
	r.GET("/category", ctrl.List)
	return r
}

func TestCategoryController_List_success(t *testing.T) {
	t1 := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
	svc := &mockCategoryService{
		ListFn: func() ([]domain.CategoryListItemResponseDTO, error) {
			return []domain.CategoryListItemResponseDTO{
				{ID: 2, Name: "B", Description: "x", CreatedAt: t1, UpdatedAt: t1},
			}, nil
		},
	}
	r := setupCategoryRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/category", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status %d", w.Code)
	}
	var body []domain.CategoryListItemResponseDTO
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if len(body) != 1 || body[0].ID != 2 {
		t.Fatalf("body %+v", body)
	}
}

func TestCategoryController_List_serviceError(t *testing.T) {
	svc := &mockCategoryService{
		ListFn: func() ([]domain.CategoryListItemResponseDTO, error) {
			return nil, errors.New("fail")
		},
	}
	r := setupCategoryRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/category", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status %d", w.Code)
	}
}
