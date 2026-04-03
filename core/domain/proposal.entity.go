package domain

import (
	"time"

	"github.com/google/uuid"
)

type Proposal struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID      uuid.UUID `gorm:"type:uuid;not null;index:idx_proposals_user_id" json:"user_id"`
	HelperID    uuid.UUID `gorm:"type:uuid;not null;index:idx_proposals_helper_id" json:"helper_id"`
	CategoryID  uint      `gorm:"not null" json:"category_id"`
	Description string    `gorm:"type:text;not null" json:"description"`
	Value       float64   `gorm:"type:decimal(10,2);not null" json:"value"`
	Status      string    `gorm:"type:varchar(50);not null;default:pending;index:idx_proposals_status" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Proposal) TableName() string { return "proposals" }
