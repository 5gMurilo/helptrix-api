package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserCategory struct {
	UserID     uuid.UUID      `gorm:"type:uuid;not null;primaryKey" json:"user_id"`
	CategoryID uint           `gorm:"not null;primaryKey" json:"category_id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (UserCategory) TableName() string {
	return "user_categories"
}
