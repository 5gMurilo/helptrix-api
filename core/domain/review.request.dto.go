package domain

type CreateReviewRequestDTO struct {
	ProposalID  string `json:"proposal_id" binding:"required,uuid"`
	HelperID    string `json:"helper_id" binding:"required,uuid"`
	Rate        int    `json:"rate" binding:"required,min=1,max=5"`
	Review      string `json:"review" binding:"required"`
	ServiceType string `json:"service_type" binding:"required"`
}
