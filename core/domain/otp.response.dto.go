package domain

import "github.com/google/uuid"

// SendOTPResponseDTO represents the response after requesting an OTP.
//
//	@name SendOTPResponseDTO
type SendOTPResponseDTO struct {
	ID      uuid.UUID `json:"id"`
	Message string    `json:"message"`
}

// ConfirmOTPResponseDTO represents the response after confirming an OTP.
//
//	@name ConfirmOTPResponseDTO
type ConfirmOTPResponseDTO struct {
	Message string `json:"message"`
}
