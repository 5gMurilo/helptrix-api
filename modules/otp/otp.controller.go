package otpmodule

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/5gMurilo/helptrix-api/core/domain"
	otpinterfaces "github.com/5gMurilo/helptrix-api/core/interfaces/otp"
	"github.com/5gMurilo/helptrix-api/core/logger"
	"github.com/5gMurilo/helptrix-api/core/utils"
	"github.com/gin-gonic/gin"
)

type OtpController struct {
	svc otpinterfaces.IOtpService
}

func NewOtpController(svc otpinterfaces.IOtpService) otpinterfaces.IOtpController {
	return &OtpController{svc: svc}
}

// Send godoc
//
//	@Summary		Send OTP
//	@Description	Generates a 4-digit OTP, invalidates any existing waiting OTP for the email, and delivers it by email
//	@Tags			otp
//	@Accept			json
//	@Produce		json
//	@Param			body	body		domain.SendOTPRequestDTO	true	"Send OTP request"
//	@Success		200		{object}	domain.SendOTPResponseDTO
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/otp/send [post]
func (ctrl *OtpController) Send(c *gin.Context) {
	log := logger.Get().With(
		slog.String("layer", "controller"),
		slog.String("module", "otp"),
		slog.String("operation", "Send"),
	)

	var dto domain.SendOTPRequestDTO

	if err := c.ShouldBindJSON(&dto); err != nil {
		log.Warn("invalid request body", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	response, err := ctrl.svc.Send(dto)
	if err != nil {
		log.Error("OTP send failed", slog.String("error", err.Error()), slog.String("email", dto.Email))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Confirm godoc
//
//	@Summary		Confirm OTP
//	@Description	Validates the submitted 4-digit OTP code and marks it as confirmed
//	@Tags			otp
//	@Accept			json
//	@Produce		json
//	@Param			body	body		domain.ConfirmOTPRequestDTO	true	"Confirm OTP request"
//	@Success		200		{object}	domain.ConfirmOTPResponseDTO
//	@Failure		400		{object}	map[string]string
//	@Failure		404		{object}	map[string]string
//	@Failure		422		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/otp/confirm [post]
func (ctrl *OtpController) Confirm(c *gin.Context) {
	log := logger.Get().With(
		slog.String("layer", "controller"),
		slog.String("module", "otp"),
		slog.String("operation", "Confirm"),
	)

	var dto domain.ConfirmOTPRequestDTO

	if err := c.ShouldBindJSON(&dto); err != nil {
		log.Warn("invalid request body", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	response, err := ctrl.svc.Confirm(dto)
	if err != nil {
		if errors.Is(err, utils.ErrOTPNotFound) {
			log.Warn("OTP confirmation failed: OTP not found", slog.String("otp_id", dto.ID))
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, utils.ErrOTPNotWaiting) || errors.Is(err, utils.ErrOTPExpired) || errors.Is(err, utils.ErrOTPInvalid) {
			log.Warn("OTP confirmation failed: OTP validation error", slog.String("error", err.Error()), slog.String("otp_id", dto.ID))
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}
		log.Error("OTP confirmation failed with unexpected error", slog.String("error", err.Error()), slog.String("otp_id", dto.ID))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
