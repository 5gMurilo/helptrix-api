package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Address struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID       uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	Street       string         `gorm:"not null" json:"street"`
	Number       string         `gorm:"not null" json:"number"`
	Complement   string         `json:"complement,omitempty"`
	Neighborhood string         `gorm:"not null" json:"neighborhood"`
	ZipCode      string         `json:"zip_code"`
	City         string         `gorm:"not null" json:"city"`
	State        string         `gorm:"not null" json:"state"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Address) TableName() string {
	return "addresses"
}
