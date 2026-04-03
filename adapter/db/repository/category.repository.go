package repository

import (
	"fmt"

	"github.com/5gMurilo/helptrix-api/core/domain"
	categoryinterfaces "github.com/5gMurilo/helptrix-api/core/interfaces/category"
	"gorm.io/gorm"
)

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) categoryinterfaces.ICategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) List() ([]domain.CategoryListItemResponseDTO, error) {
	var rows []domain.Category
	if err := r.db.Order("id ASC").Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("error listing categories: %w", err)
	}

	out := make([]domain.CategoryListItemResponseDTO, 0, len(rows))
	for _, c := range rows {
		out = append(out, domain.CategoryListItemResponseDTO{
			ID:          c.ID,
			Name:        c.Name,
			Description: c.Description,
			CreatedAt:   c.CreatedAt,
			UpdatedAt:   c.UpdatedAt,
		})
	}
	return out, nil
}
