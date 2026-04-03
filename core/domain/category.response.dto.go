package domain

import "time"

// CategoryListItemResponseDTO is one category in the public catalog list.
//
//	@name	CategoryListItemResponseDTO
type CategoryListItemResponseDTO struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
