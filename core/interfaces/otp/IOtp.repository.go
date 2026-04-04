package otpinterfaces

import (
	"github.com/5gMurilo/helptrix-api/core/domain"
	"github.com/google/uuid"
)

type IOtpRepository interface {
	Create(otp domain.OTP) (domain.OTP, error)
	FindByID(id uuid.UUID) (*domain.OTP, error)
	FindActiveByEmail(email string) (*domain.OTP, error)
	UpdateStatus(id uuid.UUID, status string) error
}
