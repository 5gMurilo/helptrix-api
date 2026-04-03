package domain

import (
	"time"

	"github.com/google/uuid"
)

// ProfileCategoryDTO is the category shape embedded inside GetProfileResponseDTO.
//
//	@name	ProfileCategoryDTO
type ProfileCategoryDTO struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ProfileAddressDTO is the address shape embedded inside GetProfileResponseDTO.
//
//	@name	ProfileAddressDTO
type ProfileAddressDTO struct {
	Street       string `json:"street"`
	Number       string `json:"number"`
	Complement   string `json:"complement,omitempty"`
	Neighborhood string `json:"neighborhood"`
	City         string `json:"city"`
	State        string `json:"state"`
}

// ProfileServiceRateDTO holds aggregated review data for a service inside GetProfileResponseDTO.
//
//	@name	ProfileServiceRateDTO
type ProfileServiceRateDTO struct {
	RateQuantity  int     `json:"rate_quantity"`
	AverageRating float64 `json:"average_rating"`
}

// ProfileServiceDTO is the service shape embedded inside GetProfileResponseDTO.
//
//	@name	ProfileServiceDTO
type ProfileServiceDTO struct {
	Title         string                `json:"title"`
	Description   string                `json:"description"`
	ActuationDays []string              `json:"actuation_days"`
	StartTime     string                `json:"start_time"`
	EndTime       string                `json:"end_time"`
	OfferSince    time.Time             `json:"offer_since"`
	Value         interface{}           `json:"value"`
	Category      ProfileCategoryDTO    `json:"category"`
	Rate          ProfileServiceRateDTO `json:"rate"`
	Photos        []string              `json:"photos,omitempty"`
}

// GetProfileResponseDTO is the full profile shape returned by GET /user/profile/:id.
//
//	@name	GetProfileResponseDTO
type GetProfileResponseDTO struct {
	ID             uuid.UUID            `json:"id"`
	Name           string               `json:"name"`
	Email          string               `json:"email"`
	Phone          string               `json:"phone"`
	Biography      string               `json:"biography,omitempty"`
	ProfilePicture string               `json:"profile_picture,omitempty"`
	UserType       string               `json:"user_type"`
	Categories     []ProfileCategoryDTO `json:"categories"`
	Address        *ProfileAddressDTO   `json:"address,omitempty"`
	Reviews        []interface{}        `json:"reviews"`
	Services       []ServiceResponseDTO `json:"services,omitempty"`
	CreatedAt      time.Time            `json:"created_at"`
}

// UpdateProfileResponseDTO is an empty struct used only for Swagger documentation of PUT /user/profile/:id.
// The actual HTTP response is 204 No Content.
//
//	@name	UpdateProfileResponseDTO
type UpdateProfileResponseDTO struct{}
