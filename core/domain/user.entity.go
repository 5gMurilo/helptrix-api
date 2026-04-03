package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name           string         `gorm:"not null" json:"name"`
	Email          string         `gorm:"uniqueIndex;not null" json:"email"`
	Document       string         `gorm:"uniqueIndex;not null" json:"document"`
	Password       string         `gorm:"not null" json:"-"`
	Phone          string         `gorm:"not null" json:"phone"`
	UserType       string         `gorm:"not null" json:"user_type"`
	Biography      string         `gorm:"type:text" json:",omitempty"`
	ProfilePicture string         `gorm:"type:varchar(500)" json:",omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	Address        Address        `gorm:"foreignKey:UserID" json:"-"`
	Categories     []Category     `gorm:"many2many:user_categories;" json:"-"`
	Services       []Service      `gorm:"foreignKey:UserID" json:"-"`
}

func (User) TableName() string {
	return "users"
}
