package repository

import (
	"errors"
	"testing"

	"github.com/5gMurilo/helptrix-api/core/domain"
	"github.com/5gMurilo/helptrix-api/core/utils"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupReviewRepositoryTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file:"+uuid.New().String()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite database: %v", err)
	}

	if err := db.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
		t.Fatalf("failed to enable foreign keys: %v", err)
	}

	statements := []string{
		`CREATE TABLE users (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			document TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			phone TEXT NOT NULL,
			user_type TEXT NOT NULL,
			biography TEXT,
			profile_picture TEXT,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		)`,
		`CREATE TABLE categories (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			description TEXT NOT NULL,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		)`,
		`CREATE TABLE proposals (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			helper_id TEXT NOT NULL,
			category_id INTEGER NOT NULL,
			description TEXT NOT NULL,
			value DECIMAL NOT NULL,
			status TEXT NOT NULL,
			created_at DATETIME,
			updated_at DATETIME
		)`,
		`CREATE TABLE reviews (
			id TEXT PRIMARY KEY,
			proposal_id TEXT NOT NULL,
			business_id TEXT NOT NULL,
			helper_id TEXT NOT NULL,
			category_id INTEGER NOT NULL,
			rate INTEGER NOT NULL,
			review TEXT NOT NULL,
			service_type TEXT NOT NULL,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		)`,
		`CREATE UNIQUE INDEX idx_reviews_unique_active_proposal
			ON reviews(proposal_id)
			WHERE deleted_at IS NULL`,
	}

	for _, statement := range statements {
		if err := db.Exec(statement).Error; err != nil {
			t.Fatalf("failed to migrate test database: %v", err)
		}
	}

	return db
}

func seedReviewRepositoryUsers(t *testing.T, db *gorm.DB) (uuid.UUID, uuid.UUID, uint) {
	t.Helper()

	businessID := uuid.New()
	helperID := uuid.New()
	category := domain.Category{Name: "Plumbing", Description: "Plumbing services"}

	if err := db.Create(&category).Error; err != nil {
		t.Fatalf("failed to create category: %v", err)
	}

	users := []domain.User{
		{
			ID:       businessID,
			Name:     "Business",
			Email:    businessID.String() + "@example.com",
			Document: "12345678901234",
			Password: "secret",
			Phone:    "11999999999",
			UserType: utils.UserTypeBusiness,
		},
		{
			ID:       helperID,
			Name:     "Helper",
			Email:    helperID.String() + "@example.com",
			Document: "12345678901",
			Password: "secret",
			Phone:    "11888888888",
			UserType: utils.UserTypeHelper,
		},
	}

	if err := db.Create(&users).Error; err != nil {
		t.Fatalf("failed to create users: %v", err)
	}

	return businessID, helperID, category.ID
}

func seedReviewRepositoryProposal(
	t *testing.T,
	db *gorm.DB,
	businessID uuid.UUID,
	helperID uuid.UUID,
	categoryID uint,
	status string,
) uuid.UUID {
	t.Helper()

	proposalID := uuid.New()
	proposal := domain.Proposal{
		ID:          proposalID,
		UserID:      businessID,
		HelperID:    helperID,
		CategoryID:  categoryID,
		Description: "Fix kitchen sink",
		Value:       100,
		Status:      status,
	}

	if err := db.Create(&proposal).Error; err != nil {
		t.Fatalf("failed to create proposal: %v", err)
	}

	return proposalID
}

func newRepositoryReview(proposalID, businessID, helperID uuid.UUID) *domain.Review {
	return &domain.Review{
		ID:          uuid.New(),
		ProposalID:  proposalID,
		BusinessID:  businessID,
		HelperID:    helperID,
		Rate:        5,
		Review:      "Great service",
		ServiceType: "Plumbing",
	}
}

func TestReviewRepository_Create_SuccessForFinishedProposal(t *testing.T) {
	db := setupReviewRepositoryTestDB(t)
	businessID, helperID, categoryID := seedReviewRepositoryUsers(t, db)
	proposalID := seedReviewRepositoryProposal(t, db, businessID, helperID, categoryID, utils.ProposalStatusFinished)
	repo := NewReviewRepository(db)

	review := newRepositoryReview(proposalID, businessID, helperID)
	err := repo.Create(review)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if review.CategoryID != categoryID {
		t.Fatalf("expected category %d, got %d", categoryID, review.CategoryID)
	}
}

func TestReviewRepository_Create_RejectsUnfinishedProposal(t *testing.T) {
	statuses := []string{
		utils.ProposalStatusPending,
		utils.ProposalStatusAccepted,
		utils.ProposalStatusInProgress,
		utils.ProposalStatusCancelled,
		utils.ProposalStatusRefused,
	}

	for _, status := range statuses {
		t.Run(status, func(t *testing.T) {
			db := setupReviewRepositoryTestDB(t)
			businessID, helperID, categoryID := seedReviewRepositoryUsers(t, db)
			proposalID := seedReviewRepositoryProposal(t, db, businessID, helperID, categoryID, status)
			repo := NewReviewRepository(db)

			err := repo.Create(newRepositoryReview(proposalID, businessID, helperID))

			if !errors.Is(err, utils.ErrProposalNotFinished) {
				t.Fatalf("expected ErrProposalNotFinished, got %v", err)
			}
		})
	}
}

func TestReviewRepository_Create_RejectsMismatchedHelper(t *testing.T) {
	db := setupReviewRepositoryTestDB(t)
	businessID, helperID, categoryID := seedReviewRepositoryUsers(t, db)
	proposalID := seedReviewRepositoryProposal(t, db, businessID, helperID, categoryID, utils.ProposalStatusFinished)
	repo := NewReviewRepository(db)

	err := repo.Create(newRepositoryReview(proposalID, businessID, uuid.New()))

	if !errors.Is(err, utils.ErrReviewProposalMismatch) {
		t.Fatalf("expected ErrReviewProposalMismatch, got %v", err)
	}
}

func TestReviewRepository_Create_RejectsDuplicateProposalReview(t *testing.T) {
	db := setupReviewRepositoryTestDB(t)
	businessID, helperID, categoryID := seedReviewRepositoryUsers(t, db)
	proposalID := seedReviewRepositoryProposal(t, db, businessID, helperID, categoryID, utils.ProposalStatusFinished)
	repo := NewReviewRepository(db)

	if err := repo.Create(newRepositoryReview(proposalID, businessID, helperID)); err != nil {
		t.Fatalf("expected first review to succeed, got %v", err)
	}

	err := repo.Create(newRepositoryReview(proposalID, businessID, helperID))

	if !errors.Is(err, utils.ErrReviewAlreadyExists) {
		t.Fatalf("expected ErrReviewAlreadyExists, got %v", err)
	}
}

func TestReviewRepository_Create_AllowsSamePairForDifferentFinishedProposals(t *testing.T) {
	db := setupReviewRepositoryTestDB(t)
	businessID, helperID, categoryID := seedReviewRepositoryUsers(t, db)
	firstProposalID := seedReviewRepositoryProposal(t, db, businessID, helperID, categoryID, utils.ProposalStatusFinished)
	secondProposalID := seedReviewRepositoryProposal(t, db, businessID, helperID, categoryID, utils.ProposalStatusFinished)
	repo := NewReviewRepository(db)

	if err := repo.Create(newRepositoryReview(firstProposalID, businessID, helperID)); err != nil {
		t.Fatalf("expected first review to succeed, got %v", err)
	}

	if err := repo.Create(newRepositoryReview(secondProposalID, businessID, helperID)); err != nil {
		t.Fatalf("expected second review to succeed, got %v", err)
	}
}
