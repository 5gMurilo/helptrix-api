package proposalinterfaces

import (
	"github.com/5gMurilo/helptrix-api/core/domain"
	"github.com/google/uuid"
)

type IProposalRepository interface {
	Create(dto domain.CreateProposalRequestDTO, userID uuid.UUID) (domain.Proposal, error)
	FindByID(id uuid.UUID) (*domain.Proposal, error)
	UpdateStatus(id uuid.UUID, status string) (*domain.Proposal, error)
	ListByUserID(userID uuid.UUID, statusFilter string) ([]domain.ProposalResponseDTO, error)
	ListByHelperID(helperID uuid.UUID, statusFilter string) ([]domain.ProposalResponseDTO, error)
	HasBlockingProposalForHelper(userID uuid.UUID, helperID uuid.UUID) (bool, error)
}
