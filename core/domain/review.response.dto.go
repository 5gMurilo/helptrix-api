package domain

import (
	"time"
)

// ReviewListResponseDTO is the response shape for review list endpoints.
//
//	@name	ReviewListResponseDTO
type ReviewListResponseDTO struct {
	Rate            int       `json:"rate"`
	Review          string    `json:"review,omitempty"`
	ServiceType     string    `json:"service_type"`
	BusinessID      string    `json:"business_id"`
	BusinessName    string    `json:"business_name" gorm:"column:business_name"`
	BusinessPicture string    `json:"business_picture" gorm:"column:business_picture"`
	CreatedAt       time.Time `json:"created_at"`
}
