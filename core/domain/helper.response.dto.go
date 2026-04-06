package domain

import "github.com/google/uuid"

type HelperCategoryDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type HelperCardDTO struct {
	ID             uuid.UUID           `json:"id"`
	Name           string              `json:"name"`
	Biography      string              `json:"biography"`
	ProfilePicture string              `json:"profile_picture"`
	City           string              `json:"city"`
	State          string              `json:"state"`
	Categories     []HelperCategoryDTO `json:"categories"`
}

type HelperListResponseDTO struct {
	Data     []HelperCardDTO `json:"data"`
	Total    int64           `json:"total"`
	Page     int             `json:"page"`
	PageSize int             `json:"page_size"`
}
