package domain

import (
	"time"

	"github.com/google/uuid"
)

// ProposalResponseDTO is the full proposal shape returned by all proposal endpoints.
//
//	@name	ProposalResponseDTO
type ProposalResponseDTO struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	HelperID    uuid.UUID `json:"helper_id"`
	CategoryID  uint      `json:"category_id"`
	Description string    `json:"description"`
	Value       float64   `json:"value"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
