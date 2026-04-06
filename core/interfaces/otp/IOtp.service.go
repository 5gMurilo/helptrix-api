package otpinterfaces

import "github.com/5gMurilo/helptrix-api/core/domain"

type IOtpService interface {
	Send(dto domain.SendOTPRequestDTO) (domain.SendOTPResponseDTO, error)
	Confirm(dto domain.ConfirmOTPRequestDTO) (domain.ConfirmOTPResponseDTO, error)
}
