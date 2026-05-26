package otpmodule

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/5gMurilo/helptrix-api/core/domain"
	emailinterfaces "github.com/5gMurilo/helptrix-api/core/interfaces/email"
	otpinterfaces "github.com/5gMurilo/helptrix-api/core/interfaces/otp"
	"github.com/5gMurilo/helptrix-api/core/logger"
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
	log := logger.Get().With(
		slog.String("layer", "service"),
		slog.String("module", "otp"),
		slog.String("operation", "Send"),
		slog.String("email", dto.Email),
	)

	existing, err := s.repo.FindActiveByEmail(dto.Email)
	if err == nil && existing != nil {
		if err = s.repo.UpdateStatus(existing.ID, utils.OTPStatusExpired); err != nil {
			log.Error("failed to expire previous OTP", slog.String("error", err.Error()), slog.String("otp_id", existing.ID.String()))
			return domain.SendOTPResponseDTO{}, errors.New("error expiring previous otp")
		}
		log.Info("previous OTP expired", slog.String("otp_id", existing.ID.String()))
	}

	buf := make([]byte, 2)
	if _, err = rand.Read(buf); err != nil {
		log.Error("failed to generate OTP code", slog.String("error", err.Error()))
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
		log.Error("failed to persist OTP", slog.String("error", err.Error()))
		return domain.SendOTPResponseDTO{}, errors.New("error creating otp")
	}

	if err = s.emailSender.Send(
		dto.Email,
		"Your OTP code for Helptrix",
		fmt.Sprintf("<html><body><h1>HELPTRIX</h1><p>Your OTP code for Helptrix registration is: <strong>%s</strong></p></body></html>", code),
	); err != nil {
		log.Error("failed to send OTP email", slog.String("error", err.Error()))
		return domain.SendOTPResponseDTO{}, errors.New("error sending otp email")
	}

	log.Info("OTP sent successfully", slog.String("otp_id", created.ID.String()))
	return domain.SendOTPResponseDTO{ID: created.ID, Message: "OTP sent successfully"}, nil
}

func (s *OtpService) Confirm(dto domain.ConfirmOTPRequestDTO) (domain.ConfirmOTPResponseDTO, error) {
	log := logger.Get().With(
		slog.String("layer", "service"),
		slog.String("module", "otp"),
		slog.String("operation", "Confirm"),
		slog.String("otp_id", dto.ID),
	)

	id, err := uuid.Parse(dto.ID)
	if err != nil {
		log.Warn("OTP confirmation rejected: invalid OTP ID format")
		return domain.ConfirmOTPResponseDTO{}, utils.ErrOTPNotFound
	}

	otp, err := s.repo.FindByID(id)
	if err != nil {
		log.Warn("OTP not found", slog.String("error", err.Error()))
		return domain.ConfirmOTPResponseDTO{}, err
	}

	if otp.Status != utils.OTPStatusWaiting {
		log.Warn("OTP confirmation rejected: OTP is not in waiting state", slog.String("current_status", otp.Status))
		return domain.ConfirmOTPResponseDTO{}, utils.ErrOTPNotWaiting
	}

	if !time.Now().Before(otp.ExpiresAt) {
		_ = s.repo.UpdateStatus(id, utils.OTPStatusExpired)
		log.Warn("OTP confirmation rejected: OTP has expired")
		return domain.ConfirmOTPResponseDTO{}, utils.ErrOTPExpired
	}

	if otp.Code != dto.Code {
		log.Warn("OTP confirmation rejected: invalid code provided")
		return domain.ConfirmOTPResponseDTO{}, utils.ErrOTPInvalid
	}

	if err = s.repo.UpdateStatus(id, utils.OTPStatusConfirmed); err != nil {
		log.Error("failed to confirm OTP status", slog.String("error", err.Error()))
		return domain.ConfirmOTPResponseDTO{}, errors.New("error confirming otp")
	}

	log.Info("OTP confirmed successfully")
	return domain.ConfirmOTPResponseDTO{Message: "OTP confirmed successfully"}, nil
}
