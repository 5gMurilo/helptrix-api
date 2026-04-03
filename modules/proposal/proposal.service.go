package proposal

import (
	"github.com/5gMurilo/helptrix-api/core/domain"
	proposalinterfaces "github.com/5gMurilo/helptrix-api/core/interfaces/proposal"
	"github.com/5gMurilo/helptrix-api/core/utils"
	"github.com/google/uuid"
)

type ProposalService struct {
	repo proposalinterfaces.IProposalRepository
}

func NewProposalService(repo proposalinterfaces.IProposalRepository) proposalinterfaces.IProposalService {
	return &ProposalService{repo: repo}
}

func (s *ProposalService) Create(dto domain.CreateProposalRequestDTO, userID uuid.UUID) (domain.ProposalResponseDTO, error) {
	hasActive, err := s.repo.HasActiveProposal(userID)
	if err != nil {
		return domain.ProposalResponseDTO{}, err
	}
	if hasActive {
		return domain.ProposalResponseDTO{}, utils.ErrProposalAlreadyActive
	}

	proposal, err := s.repo.Create(dto, userID)
	if err != nil {
		return domain.ProposalResponseDTO{}, err
	}

	return toResponseDTO(proposal), nil
}

func (s *ProposalService) GetByID(proposalID uuid.UUID, requesterID uuid.UUID) (domain.ProposalResponseDTO, error) {
	proposal, err := s.repo.FindByID(proposalID)
	if err != nil {
		return domain.ProposalResponseDTO{}, err
	}

	if proposal.UserID != requesterID && proposal.HelperID != requesterID {
		return domain.ProposalResponseDTO{}, utils.ErrNotProposalParticipant
	}

	return toResponseDTO(*proposal), nil
}

func (s *ProposalService) UpdateStatus(
	proposalID uuid.UUID,
	dto domain.UpdateProposalStatusRequestDTO,
	requesterID uuid.UUID,
	requesterType string,
) (domain.ProposalResponseDTO, error) {
	proposal, err := s.repo.FindByID(proposalID)
	if err != nil {
		return domain.ProposalResponseDTO{}, err
	}

	terminalStatuses := map[string]bool{
		utils.ProposalStatusRefused:   true,
		utils.ProposalStatusCancelled: true,
		utils.ProposalStatusFinished:  true,
	}
	if terminalStatuses[proposal.Status] {
		return domain.ProposalResponseDTO{}, utils.ErrProposalFinished
	}

	validStatuses := map[string]bool{
		utils.ProposalStatusPending:    true,
		utils.ProposalStatusAccepted:   true,
		utils.ProposalStatusRefused:    true,
		utils.ProposalStatusInProgress: true,
		utils.ProposalStatusCancelled:  true,
		utils.ProposalStatusFinished:   true,
	}
	if !validStatuses[dto.Status] {
		return domain.ProposalResponseDTO{}, utils.ErrProposalInvalidStatus
	}

	if dto.Status == utils.ProposalStatusCancelled {
		if requesterID != proposal.UserID && requesterID != proposal.HelperID {
			return domain.ProposalResponseDTO{}, utils.ErrProposalUnauthorized
		}
	} else {
		if requesterType != utils.UserTypeHelper || requesterID != proposal.HelperID {
			return domain.ProposalResponseDTO{}, utils.ErrProposalUnauthorized
		}
	}

	allowedTransitions := map[string]map[string]bool{
		utils.ProposalStatusPending: {
			utils.ProposalStatusAccepted:  true,
			utils.ProposalStatusRefused:   true,
			utils.ProposalStatusCancelled: true,
		},
		utils.ProposalStatusAccepted: {
			utils.ProposalStatusInProgress: true,
			utils.ProposalStatusCancelled:  true,
		},
		utils.ProposalStatusInProgress: {
			utils.ProposalStatusFinished:  true,
			utils.ProposalStatusCancelled: true,
		},
	}

	if allowed, ok := allowedTransitions[proposal.Status]; !ok || !allowed[dto.Status] {
		return domain.ProposalResponseDTO{}, utils.ErrProposalInvalidStatus
	}

	updated, err := s.repo.UpdateStatus(proposalID, dto.Status)
	if err != nil {
		return domain.ProposalResponseDTO{}, err
	}

	return toResponseDTO(*updated), nil
}

func (s *ProposalService) List(requesterID uuid.UUID, requesterType string, statusFilter string) ([]domain.ProposalResponseDTO, error) {
	if requesterType == utils.UserTypeBusiness {
		return s.repo.ListByUserID(requesterID, statusFilter)
	}
	return s.repo.ListByHelperID(requesterID, statusFilter)
}

func toResponseDTO(p domain.Proposal) domain.ProposalResponseDTO {
	return domain.ProposalResponseDTO{
		ID:          p.ID,
		UserID:      p.UserID,
		HelperID:    p.HelperID,
		CategoryID:  p.CategoryID,
		Description: p.Description,
		Value:       p.Value,
		Status:      p.Status,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}
