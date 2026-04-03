package domain

import "time"

// CreateServiceRequestDTO holds the required and optional fields for POST /service.
//
//	@name	CreateServiceRequestDTO
type CreateServiceRequestDTO struct {
	Name          string    `json:"name"           binding:"required"`
	Description   string    `json:"description"    binding:"required"`
	ActuationDays []string  `json:"actuation_days" binding:"required,min=1"`
	Value         string    `json:"value"          binding:"required"`
	StartTime     string    `json:"start_time"     binding:"required"`
	EndTime       string    `json:"end_time"       binding:"required"`
	OfferSince    time.Time `json:"offer_since"    binding:"required"`
	CategoryID    uint      `json:"category_id"    binding:"required"`
	Photos        []string  `json:"photos,omitempty"`
}

// UpdateServiceRequestDTO holds optional fields for PUT /service/:id.
// All fields are pointers so the service layer can detect absent vs. zero-value.
//
//	@name	UpdateServiceRequestDTO
type UpdateServiceRequestDTO struct {
	Name          *string    `json:"name,omitempty"`
	Description   *string    `json:"description,omitempty"`
	ActuationDays []string   `json:"actuation_days,omitempty"`
	Value         *string    `json:"value,omitempty"`
	StartTime     *string    `json:"start_time,omitempty"`
	EndTime       *string    `json:"end_time,omitempty"`
	OfferSince    *time.Time `json:"offer_since,omitempty"`
	CategoryID    *uint      `json:"category_id,omitempty"`
	Photos        []string   `json:"photos,omitempty"`
}
