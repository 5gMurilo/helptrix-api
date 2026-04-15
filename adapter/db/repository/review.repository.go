package repository

import (
	"errors"
	"fmt"

	"github.com/5gMurilo/helptrix-api/core/domain"
	reviewinterfaces "github.com/5gMurilo/helptrix-api/core/interfaces/review"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type reviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) reviewinterfaces.IReviewRepository {
	return &reviewRepository{db: db}
}

func (r *reviewRepository) Create(review *domain.Review) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Verify proposal exists and is finished
	var proposal domain.Proposal
	err := tx.Where("helper_id = ? AND user_id = ? AND status = ?", review.HelperID, review.BusinessID, "finished").
		First(&proposal).Error
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("no finished proposal found for this business and helper")
		}
		return fmt.Errorf("error verifying proposal: %w", err)
	}

	// Set category_id from the proposal
	review.CategoryID = proposal.CategoryID

	// Check for existing review (active, not soft-deleted)
	var existing domain.Review
	existingErr := tx.Where("business_id = ? AND helper_id = ? AND deleted_at IS NULL",
		review.BusinessID, review.HelperID).
		First(&existing).Error
	if existingErr == nil {
		tx.Rollback()
		return fmt.Errorf("business has already reviewed this helper")
	}
	if !errors.Is(existingErr, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return fmt.Errorf("error checking existing review: %w", existingErr)
	}

	if err := tx.Create(review).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("error creating review: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}

func (r *reviewRepository) ListByBusiness(businessID uuid.UUID) ([]domain.Review, error) {
	var reviews []domain.Review
	err := r.db.Where("business_id = ? AND deleted_at IS NULL", businessID).
		Order("created_at DESC").
		Find(&reviews).Error
	if err != nil {
		return nil, fmt.Errorf("error listing reviews by business: %w", err)
	}
	return reviews, nil
}

func (r *reviewRepository) ListByHelper(helperID uuid.UUID) ([]domain.Review, error) {
	var reviews []domain.Review
	err := r.db.Where("helper_id = ? AND deleted_at IS NULL", helperID).
		Order("created_at DESC").
		Find(&reviews).Error
	if err != nil {
		return nil, fmt.Errorf("error listing reviews by helper: %w", err)
	}
	return reviews, nil
}

func (r *reviewRepository) GetByBusinessAndHelper(businessID, helperID uuid.UUID) (*domain.Review, error) {
	var review domain.Review
	err := r.db.Where("business_id = ? AND helper_id = ? AND deleted_at IS NULL",
		businessID, helperID).
		First(&review).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("error fetching review: %w", err)
	}
	return &review, nil
}
