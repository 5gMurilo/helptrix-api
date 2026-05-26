package category

import (
	"log/slog"

	"github.com/5gMurilo/helptrix-api/core/domain"
	categoryinterfaces "github.com/5gMurilo/helptrix-api/core/interfaces/category"
	"github.com/5gMurilo/helptrix-api/core/logger"
)

type CategoryService struct {
	repo categoryinterfaces.ICategoryRepository
}

func NewCategoryService(repo categoryinterfaces.ICategoryRepository) categoryinterfaces.ICategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) List() ([]domain.CategoryListItemResponseDTO, error) {
	log := logger.Get().With(
		slog.String("layer", "service"),
		slog.String("module", "category"),
		slog.String("operation", "List"),
	)

	result, err := s.repo.List()
	if err != nil {
		log.Error("failed to list categories", slog.String("error", err.Error()))
		return nil, err
	}

	return result, nil
}

func (s *CategoryService) Seed(categories []domain.Category) error {
	return s.repo.Seed(categories)
}
