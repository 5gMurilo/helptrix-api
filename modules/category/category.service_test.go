package category_test

import (
	"errors"
	"testing"
	"time"

	"github.com/5gMurilo/helptrix-api/core/domain"
	categoryinterfaces "github.com/5gMurilo/helptrix-api/core/interfaces/category"
	categorymodule "github.com/5gMurilo/helptrix-api/modules/category"
)

type mockCategoryRepository struct {
	ListFn func() ([]domain.CategoryListItemResponseDTO, error)
	SeedFn func(categories []domain.Category) error
}

func (m *mockCategoryRepository) List() ([]domain.CategoryListItemResponseDTO, error) {
	return m.ListFn()
}

func (m *mockCategoryRepository) Seed(categories []domain.Category) error {
	return m.SeedFn(categories)
}

var _ categoryinterfaces.ICategoryRepository = (*mockCategoryRepository)(nil)

func TestCategoryService_List_success(t *testing.T) {
	t1 := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
	repo := &mockCategoryRepository{
		ListFn: func() ([]domain.CategoryListItemResponseDTO, error) {
			return []domain.CategoryListItemResponseDTO{
				{ID: 1, Name: "A", Description: "d1", CreatedAt: t1, UpdatedAt: t1},
			}, nil
		},
	}
	svc := categorymodule.NewCategoryService(repo)
	out, err := svc.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(out) != 1 || out[0].Name != "A" {
		t.Fatalf("unexpected list: %+v", out)
	}
}

func TestCategoryService_List_empty(t *testing.T) {
	repo := &mockCategoryRepository{
		ListFn: func() ([]domain.CategoryListItemResponseDTO, error) {
			return []domain.CategoryListItemResponseDTO{}, nil
		},
	}
	svc := categorymodule.NewCategoryService(repo)
	out, err := svc.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(out) != 0 {
		t.Fatalf("expected empty slice")
	}
}

func TestCategoryService_List_repoError(t *testing.T) {
	repo := &mockCategoryRepository{
		ListFn: func() ([]domain.CategoryListItemResponseDTO, error) {
			return nil, errors.New("db down")
		},
	}
	svc := categorymodule.NewCategoryService(repo)
	_, err := svc.List()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCategoryService_Seed_success(t *testing.T) {
	var received []domain.Category
	repo := &mockCategoryRepository{
		SeedFn: func(categories []domain.Category) error {
			received = categories
			return nil
		},
	}
	svc := categorymodule.NewCategoryService(repo)
	input := []domain.Category{
		{Name: "IT Support", Description: "Technical assistance"},
		{Name: "Plumbing", Description: "Pipe installation and repair"},
	}
	if err := svc.Seed(input); err != nil {
		t.Fatalf("Seed: %v", err)
	}
	if len(received) != 2 {
		t.Fatalf("expected 2 categories, got %d", len(received))
	}
}

func TestCategoryService_Seed_repoError(t *testing.T) {
	repo := &mockCategoryRepository{
		SeedFn: func(categories []domain.Category) error {
			return errors.New("db down")
		},
	}
	svc := categorymodule.NewCategoryService(repo)
	if err := svc.Seed([]domain.Category{{Name: "X", Description: "Y"}}); err == nil {
		t.Fatal("expected error")
	}
}
