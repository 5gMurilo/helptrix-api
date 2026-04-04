package domain

// SendOTPRequestDTO represents the payload for requesting an OTP.
//
//	@name SendOTPRequestDTO
type SendOTPRequestDTO struct {
	Email string `json:"email" binding:"required,email"`
}

// ConfirmOTPRequestDTO represents the payload for confirming an OTP.
//
//	@name ConfirmOTPRequestDTO
type ConfirmOTPRequestDTO struct {
	ID   string `json:"id"   binding:"required,uuid"`
	Code string `json:"code" binding:"required,len=4"`
}
