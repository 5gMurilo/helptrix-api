package domain

import "github.com/google/uuid"

// CreateProposalRequestDTO holds the required fields for POST /proposal.
//
//	@name	CreateProposalRequestDTO
type CreateProposalRequestDTO struct {
	HelperID    uuid.UUID `json:"helper_id"    binding:"required"`
	CategoryID  uint      `json:"category_id"  binding:"required"`
	Description string    `json:"description"  binding:"required"`
	Value       float64   `json:"value"        binding:"required,gt=0"`
}

// UpdateProposalStatusRequestDTO holds the required fields for PATCH /proposal/:id/status.
//
//	@name	UpdateProposalStatusRequestDTO
type UpdateProposalStatusRequestDTO struct {
	Status string `json:"status" binding:"required"`
}
