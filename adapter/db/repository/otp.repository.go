package repository

import (
	"errors"

	"github.com/5gMurilo/helptrix-api/core/domain"
	otpinterfaces "github.com/5gMurilo/helptrix-api/core/interfaces/otp"
	"github.com/5gMurilo/helptrix-api/core/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type otpRepository struct {
	db *gorm.DB
}

func NewOtpRepository(db *gorm.DB) otpinterfaces.IOtpRepository {
	return &otpRepository{db: db}
}

func (r *otpRepository) Create(otp domain.OTP) (domain.OTP, error) {
	tx := r.db.Begin()
	if tx.Error != nil {
		return domain.OTP{}, tx.Error
	}
	if err := tx.Create(&otp).Error; err != nil {
		tx.Rollback()
		return domain.OTP{}, errors.New("error creating otp record")
	}
	if err := tx.Commit().Error; err != nil {
		return domain.OTP{}, errors.New("error committing otp creation")
	}
	return otp, nil
}

func (r *otpRepository) FindByID(id uuid.UUID) (*domain.OTP, error) {
	var otp domain.OTP
	result := r.db.Where("id = ?", id).First(&otp)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, utils.ErrOTPNotFound
		}
		return nil, result.Error
	}
	return &otp, nil
}

func (r *otpRepository) FindActiveByEmail(email string) (*domain.OTP, error) {
	var otp domain.OTP
	result := r.db.Where("email = ? AND status = ?", email, utils.OTPStatusWaiting).First(&otp)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &otp, nil
}

func (r *otpRepository) UpdateStatus(id uuid.UUID, status string) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	result := tx.Model(&domain.OTP{}).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		tx.Rollback()
		return errors.New("error updating otp status")
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return utils.ErrOTPNotFound
	}
	return tx.Commit().Error
}
