package proposalinterfaces

import (
	"github.com/5gMurilo/helptrix-api/core/domain"
	"github.com/google/uuid"
)

type IProposalService interface {
	Create(dto domain.CreateProposalRequestDTO, userID uuid.UUID) (domain.ProposalResponseDTO, error)
	GetByID(proposalID uuid.UUID, requesterID uuid.UUID) (domain.ProposalResponseDTO, error)
	UpdateStatus(proposalID uuid.UUID, dto domain.UpdateProposalStatusRequestDTO, requesterID uuid.UUID, requesterType string) (domain.ProposalResponseDTO, error)
	List(requesterID uuid.UUID, requesterType string, statusFilter string) ([]domain.ProposalResponseDTO, error)
}
