package domain

// AddressUpdateDTO holds optional address fields for profile update.
//
//	@name	AddressUpdateDTO
type AddressUpdateDTO struct {
	Street       string `json:"street"`
	Number       string `json:"number"`
	Complement   string `json:"complement,omitempty"`
	Neighborhood string `json:"neighborhood"`
	City         string `json:"city"`
	State        string `json:"state"`
}

// UpdateProfileRequestDTO holds the optional fields a user may update on their profile.
// All fields are optional; only non-zero values are applied.
//
//	@name	UpdateProfileRequestDTO
type UpdateProfileRequestDTO struct {
	Email          string            `json:"email,omitempty"`
	Phone          string            `json:"phone,omitempty"`
	Biography      string            `json:"biography,omitempty"`
	ProfilePicture string            `json:"profile_picture,omitempty"`
	Categories     []uint            `json:"categories,omitempty"`
	Address        *AddressUpdateDTO `json:"address,omitempty"`
}

// ProfileFilters carries optional query-param filters for GET /user/profile/:id.
// Only applied when the requester has user_type == "business".
type ProfileFilters struct {
	CategoryID    *uint
	ActuationDays []string
}
