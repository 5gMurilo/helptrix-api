package categoryinterfaces

import "github.com/5gMurilo/helptrix-api/core/domain"

type ICategoryRepository interface {
	List() ([]domain.CategoryListItemResponseDTO, error)
}
