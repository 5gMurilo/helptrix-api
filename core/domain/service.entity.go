package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Service struct {
	ID            uuid.UUID       `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID        uuid.UUID       `gorm:"type:uuid;not null" json:"user_id"`
	CategoryID    uint            `gorm:"not null" json:"category_id"`
	Name          string          `gorm:"type:varchar(255);not null" json:"name"`
	Description   string          `gorm:"type:text;not null" json:"description"`
	ActuationDays datatypes.JSON  `gorm:"type:jsonb;not null" json:"actuation_days"`
	Value         decimal.Decimal `gorm:"type:decimal(10,2);not null" json:"value"`
	StartTime     string          `gorm:"type:varchar(5);not null" json:"start_time"` // HH:MM
	EndTime       string          `gorm:"type:varchar(5);not null" json:"end_time"`   // HH:MM
	OfferSince    time.Time       `gorm:"not null" json:"offer_since"`
	Photos        datatypes.JSON  `gorm:"type:jsonb" json:"photos,omitempty"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
	DeletedAt     gorm.DeletedAt  `gorm:"index" json:"-"`
	User          User            `gorm:"foreignKey:UserID" json:"-"`
	Category      Category        `gorm:"foreignKey:CategoryID" json:"-"`
}

func (Service) TableName() string { return "services" }
