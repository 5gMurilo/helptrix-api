package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Review struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	BusinessID  uuid.UUID      `gorm:"type:uuid;not null;index:idx_reviews_business_id" json:"business_id"`
	HelperID    uuid.UUID      `gorm:"type:uuid;not null;index:idx_reviews_helper_id" json:"helper_id"`
	CategoryID  uint           `gorm:"not null" json:"category_id"`
	Rate        int            `gorm:"not null" json:"rate"`
	Review      string         `gorm:"type:text;not null" json:"review"`
	ServiceType string         `gorm:"type:varchar(255);not null" json:"service_type"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index:idx_reviews_deleted_at" json:"-"`
}

func (Review) TableName() string {
	return "reviews"
}
