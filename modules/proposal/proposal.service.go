package proposal

import (
	"log/slog"

	"github.com/5gMurilo/helptrix-api/core/domain"
	proposalinterfaces "github.com/5gMurilo/helptrix-api/core/interfaces/proposal"
	"github.com/5gMurilo/helptrix-api/core/logger"
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
	log := logger.Get().With(
		slog.String("layer", "service"),
		slog.String("module", "proposal"),
		slog.String("operation", "Create"),
		slog.String("user_id", userID.String()),
	)

	hasBlocking, err := s.repo.HasBlockingProposalForHelper(userID, dto.HelperID)
	if err != nil {
		log.Error("failed to check blocking proposals", slog.String("error", err.Error()), slog.String("helper_id", dto.HelperID.String()))
		return domain.ProposalResponseDTO{}, err
	}
	if hasBlocking {
		log.Warn("proposal creation rejected: pending proposal already exists for this helper", slog.String("helper_id", dto.HelperID.String()))
		return domain.ProposalResponseDTO{}, utils.ErrProposalAlreadyPendingForHelper
	}

	proposal, err := s.repo.Create(dto, userID)
	if err != nil {
		log.Error("failed to create proposal", slog.String("error", err.Error()))
		return domain.ProposalResponseDTO{}, err
	}

	log.Info("proposal created successfully", slog.String("proposal_id", proposal.ID.String()))
	return toResponseDTO(proposal), nil
}

func (s *ProposalService) GetByID(proposalID uuid.UUID, requesterID uuid.UUID) (domain.ProposalResponseDTO, error) {
	log := logger.Get().With(
		slog.String("layer", "service"),
		slog.String("module", "proposal"),
		slog.String("operation", "GetByID"),
		slog.String("proposal_id", proposalID.String()),
		slog.String("requester_id", requesterID.String()),
	)

	proposal, err := s.repo.FindByID(proposalID)
	if err != nil {
		log.Error("proposal not found", slog.String("error", err.Error()))
		return domain.ProposalResponseDTO{}, err
	}

	if proposal.UserID != requesterID && proposal.HelperID != requesterID {
		log.Warn("proposal access denied: requester is not a participant")
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
	log := logger.Get().With(
		slog.String("layer", "service"),
		slog.String("module", "proposal"),
		slog.String("operation", "UpdateStatus"),
		slog.String("proposal_id", proposalID.String()),
		slog.String("requester_id", requesterID.String()),
	)

	proposal, err := s.repo.FindByID(proposalID)
	if err != nil {
		log.Error("proposal not found for status update", slog.String("error", err.Error()))
		return domain.ProposalResponseDTO{}, err
	}

	terminalStatuses := map[string]bool{
		utils.ProposalStatusRefused:   true,
		utils.ProposalStatusCancelled: true,
		utils.ProposalStatusFinished:  true,
	}
	if terminalStatuses[proposal.Status] {
		log.Warn("status update rejected: proposal is already in a terminal state",
			slog.String("current_status", proposal.Status),
			slog.String("requested_status", dto.Status),
		)
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
		log.Warn("status update rejected: invalid target status", slog.String("requested_status", dto.Status))
		return domain.ProposalResponseDTO{}, utils.ErrProposalInvalidStatus
	}

	if dto.Status == utils.ProposalStatusCancelled {
		if requesterID != proposal.UserID && requesterID != proposal.HelperID {
			log.Warn("status update rejected: requester is not a participant")
			return domain.ProposalResponseDTO{}, utils.ErrProposalUnauthorized
		}
	} else {
		if requesterType != utils.UserTypeHelper || requesterID != proposal.HelperID {
			log.Warn("status update rejected: only the assigned helper can perform this transition",
				slog.String("requester_type", requesterType),
			)
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
		log.Warn("status update rejected: transition not allowed",
			slog.String("from_status", proposal.Status),
			slog.String("to_status", dto.Status),
		)
		return domain.ProposalResponseDTO{}, utils.ErrProposalInvalidStatus
	}

	updated, err := s.repo.UpdateStatus(proposalID, dto.Status)
	if err != nil {
		log.Error("failed to update proposal status", slog.String("error", err.Error()))
		return domain.ProposalResponseDTO{}, err
	}

	log.Info("proposal status updated",
		slog.String("from_status", proposal.Status),
		slog.String("to_status", dto.Status),
	)
	return toResponseDTO(*updated), nil
}

func (s *ProposalService) List(requesterID uuid.UUID, requesterType string, statusFilter string, p domain.PaginationParams) ([]domain.ProposalResponseDTO, error) {
	log := logger.Get().With(
		slog.String("layer", "service"),
		slog.String("module", "proposal"),
		slog.String("operation", "List"),
		slog.String("requester_id", requesterID.String()),
	)

	var (
		result []domain.ProposalResponseDTO
		err    error
	)

	if requesterType == utils.UserTypeBusiness {
		result, err = s.repo.ListByUserID(requesterID, statusFilter, p)
	} else {
		result, err = s.repo.ListByHelperID(requesterID, statusFilter, p)
	}

	if err != nil {
		log.Error("failed to list proposals", slog.String("error", err.Error()))
		return nil, err
	}

	return result, nil
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
