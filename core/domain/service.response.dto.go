package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// ServiceCategoryDTO is the minimal category shape embedded inside ServiceResponseDTO.
//
//	@name	ServiceCategoryDTO
type ServiceCategoryDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// ServiceResponseDTO is the full service shape returned by POST /service and
// embedded inside GetProfileResponseDTO.Services.
//
//	@name	ServiceResponseDTO
type ServiceResponseDTO struct {
	ID            uuid.UUID          `json:"id"`
	Name          string             `json:"name"`
	Description   string             `json:"description"`
	ActuationDays []string           `json:"actuation_days"`
	Value         decimal.Decimal    `json:"value"`
	StartTime     string             `json:"start_time"`
	EndTime       string             `json:"end_time"`
	OfferSince    time.Time          `json:"offer_since"`
	Category      ServiceCategoryDTO `json:"category"`
	Photos        []string           `json:"photos,omitempty"`
}
