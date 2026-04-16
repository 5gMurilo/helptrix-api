package repository

import (
	"errors"
	"fmt"

	"github.com/5gMurilo/helptrix-api/core/domain"
	proposalinterfaces "github.com/5gMurilo/helptrix-api/core/interfaces/proposal"
	"github.com/5gMurilo/helptrix-api/core/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type proposalRepository struct {
	db *gorm.DB
}

func NewProposalRepository(db *gorm.DB) proposalinterfaces.IProposalRepository {
	return &proposalRepository{db: db}
}

func (r *proposalRepository) Create(dto domain.CreateProposalRequestDTO, userID uuid.UUID) (domain.Proposal, error) {
	proposal := domain.Proposal{
		UserID:      userID,
		HelperID:    dto.HelperID,
		CategoryID:  dto.CategoryID,
		Description: dto.Description,
		Value:       dto.Value,
		Status:      utils.ProposalStatusPending,
	}

	if err := r.db.Create(&proposal).Error; err != nil {
		return domain.Proposal{}, fmt.Errorf("error creating proposal: %w", err)
	}

	return proposal, nil
}

func (r *proposalRepository) FindByID(id uuid.UUID) (*domain.Proposal, error) {
	var proposal domain.Proposal

	if err := r.db.Where("id = ?", id).First(&proposal).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrProposalNotFound
		}
		return nil, fmt.Errorf("error finding proposal: %w", err)
	}

	return &proposal, nil
}

func (r *proposalRepository) UpdateStatus(id uuid.UUID, status string) (*domain.Proposal, error) {
	var proposal domain.Proposal

	if err := r.db.Where("id = ?", id).First(&proposal).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrProposalNotFound
		}
		return nil, fmt.Errorf("error finding proposal: %w", err)
	}

	if err := r.db.Model(&proposal).Update("status", status).Error; err != nil {
		return nil, fmt.Errorf("error updating proposal status: %w", err)
	}

	return &proposal, nil
}

func (r *proposalRepository) HasBlockingProposalForHelper(userID, helperID uuid.UUID) (bool, error) {
	var count int64

	if err := r.db.Model(&domain.Proposal{}).
		Where("user_id = ? AND helper_id = ? AND status != ?", userID, helperID, utils.ProposalStatusAccepted).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("error checking proposal for helper: %w", err)
	}

	return count > 0, nil
}

func (r *proposalRepository) ListByUserID(userID uuid.UUID, statusFilter string) ([]domain.ProposalResponseDTO, error) {
	var proposals []domain.Proposal

	query := r.db.Where("user_id = ?", userID)
	if statusFilter != "" {
		query = query.Where("status = ?", statusFilter)
	}

	if err := query.Find(&proposals).Error; err != nil {
		return nil, fmt.Errorf("error listing proposals by user: %w", err)
	}

	return mapProposalsToDTO(proposals), nil
}

func (r *proposalRepository) ListByHelperID(helperID uuid.UUID, statusFilter string) ([]domain.ProposalResponseDTO, error) {
	var proposals []domain.Proposal

	query := r.db.Where("helper_id = ?", helperID)
	if statusFilter != "" {
		query = query.Where("status = ?", statusFilter)
	}

	if err := query.Find(&proposals).Error; err != nil {
		return nil, fmt.Errorf("error listing proposals by helper: %w", err)
	}

	return mapProposalsToDTO(proposals), nil
}

func mapProposalsToDTO(proposals []domain.Proposal) []domain.ProposalResponseDTO {
	result := make([]domain.ProposalResponseDTO, 0, len(proposals))
	for _, p := range proposals {
		result = append(result, domain.ProposalResponseDTO{
			ID:          p.ID,
			UserID:      p.UserID,
			HelperID:    p.HelperID,
			CategoryID:  p.CategoryID,
			Description: p.Description,
			Value:       p.Value,
			Status:      p.Status,
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
		})
	}
	return result
}
