package domain

type CreateReviewRequestDTO struct {
	HelperID    string `json:"helper_id" validate:"required,uuid"`
	Rate        int    `json:"rate" validate:"required,min=1,max=5"`
	Review      string `json:"review" validate:"required"`
	ServiceType string `json:"service_type" validate:"required"`
}
