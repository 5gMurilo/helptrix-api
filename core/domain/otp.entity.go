package domain

import (
	"time"

	"github.com/google/uuid"
)

// OTP represents a one-time password record.
//
//	@name OTP
type OTP struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Email     string    `gorm:"type:varchar(255);not null;index"               json:"email"`
	Code      string    `gorm:"type:char(4);not null"                          json:"-"`
	Status    string    `gorm:"type:varchar(20);not null;default:'waiting'"    json:"status"`
	CreatedAt time.Time `                                                       json:"created_at"`
	ExpiresAt time.Time `gorm:"not null"                                        json:"expires_at"`
}

func (OTP) TableName() string { return "otps" }
