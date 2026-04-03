package category

import (
	"github.com/5gMurilo/helptrix-api/core/domain"
	categoryinterfaces "github.com/5gMurilo/helptrix-api/core/interfaces/category"
)

type CategoryService struct {
	repo categoryinterfaces.ICategoryRepository
}

func NewCategoryService(repo categoryinterfaces.ICategoryRepository) categoryinterfaces.ICategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) List() ([]domain.CategoryListItemResponseDTO, error) {
	return s.repo.List()
}
