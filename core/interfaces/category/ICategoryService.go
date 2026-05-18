package categoryinterfaces

import "github.com/5gMurilo/helptrix-api/core/domain"

type ICategoryService interface {
	List() ([]domain.CategoryListItemResponseDTO, error)
	Seed(categories []domain.Category) error
}
