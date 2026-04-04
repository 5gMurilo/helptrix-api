package otpmodule

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/5gMurilo/helptrix-api/core/domain"
	emailinterfaces "github.com/5gMurilo/helptrix-api/core/interfaces/email"
	otpinterfaces "github.com/5gMurilo/helptrix-api/core/interfaces/otp"
	"github.com/5gMurilo/helptrix-api/core/utils"
	"github.com/google/uuid"
)

type OtpService struct {
	repo        otpinterfaces.IOtpRepository
	emailSender emailinterfaces.IEmailSender
}

func NewOtpService(repo otpinterfaces.IOtpRepository, emailSender emailinterfaces.IEmailSender) otpinterfaces.IOtpService {
	return &OtpService{repo: repo, emailSender: emailSender}
}

func (s *OtpService) Send(dto domain.SendOTPRequestDTO) (domain.SendOTPResponseDTO, error) {
	existing, err := s.repo.FindActiveByEmail(dto.Email)

	log.Println("existing otp", existing, err)

	if err == nil && existing != nil {
		if err = s.repo.UpdateStatus(existing.ID, utils.OTPStatusExpired); err != nil {
			return domain.SendOTPResponseDTO{}, errors.New("error expiring previous otp")
		}
	}

	buf := make([]byte, 2)
	if _, err = rand.Read(buf); err != nil {
		return domain.SendOTPResponseDTO{}, errors.New("error generating otp code")
	}
	n := binary.BigEndian.Uint16(buf) % 10000
	code := fmt.Sprintf("%04d", n)

	otp := domain.OTP{
		Email:     dto.Email,
		Code:      code,
		Status:    utils.OTPStatusWaiting,
		ExpiresAt: time.Now().Add(utils.OTPExpirationDuration),
	}

	created, err := s.repo.Create(otp)
	if err != nil {
		return domain.SendOTPResponseDTO{}, errors.New("error creating otp")
	}

	if err = s.emailSender.Send(
		dto.Email,
		"Your OTP code for Helptrix",
		fmt.Sprintf("<html><body><h1>HELPTRIX</h1><p>Your OTP code for Helptrix registration is: <strong>%s</strong></p></body></html>", code),
	); err != nil {
		return domain.SendOTPResponseDTO{}, errors.New("error sending otp email")
	}

	return domain.SendOTPResponseDTO{ID: created.ID, Message: "OTP sent successfully"}, nil
}

func (s *OtpService) Confirm(dto domain.ConfirmOTPRequestDTO) (domain.ConfirmOTPResponseDTO, error) {
	id, err := uuid.Parse(dto.ID)
	if err != nil {
		return domain.ConfirmOTPResponseDTO{}, utils.ErrOTPNotFound
	}

	otp, err := s.repo.FindByID(id)
	if err != nil {
		return domain.ConfirmOTPResponseDTO{}, err
	}

	if otp.Status != utils.OTPStatusWaiting {
		return domain.ConfirmOTPResponseDTO{}, utils.ErrOTPNotWaiting
	}

	if !time.Now().Before(otp.ExpiresAt) {
		_ = s.repo.UpdateStatus(id, utils.OTPStatusExpired)
		return domain.ConfirmOTPResponseDTO{}, utils.ErrOTPExpired
	}

	if otp.Code != dto.Code {
		return domain.ConfirmOTPResponseDTO{}, utils.ErrOTPInvalid
	}

	if err = s.repo.UpdateStatus(id, utils.OTPStatusConfirmed); err != nil {
		return domain.ConfirmOTPResponseDTO{}, errors.New("error confirming otp")
	}

	return domain.ConfirmOTPResponseDTO{Message: "OTP confirmed successfully"}, nil
}
