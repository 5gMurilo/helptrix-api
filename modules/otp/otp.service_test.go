package otpmodule_test

import (
	"errors"
	"testing"
	"time"

	"github.com/5gMurilo/helptrix-api/core/domain"
	otpmodule "github.com/5gMurilo/helptrix-api/modules/otp"
	"github.com/5gMurilo/helptrix-api/core/utils"
	"github.com/google/uuid"
)

// mockOtpRepository implementa IOtpRepository para testes.
type mockOtpRepository struct {
	CreateFn            func(otp domain.OTP) (domain.OTP, error)
	FindByIDFn          func(id uuid.UUID) (*domain.OTP, error)
	FindActiveByEmailFn func(email string) (*domain.OTP, error)
	UpdateStatusFn      func(id uuid.UUID, status string) error
}

func (m *mockOtpRepository) Create(otp domain.OTP) (domain.OTP, error) {
	return m.CreateFn(otp)
}

func (m *mockOtpRepository) FindByID(id uuid.UUID) (*domain.OTP, error) {
	return m.FindByIDFn(id)
}

func (m *mockOtpRepository) FindActiveByEmail(email string) (*domain.OTP, error) {
	return m.FindActiveByEmailFn(email)
}

func (m *mockOtpRepository) UpdateStatus(id uuid.UUID, status string) error {
	return m.UpdateStatusFn(id, status)
}

// mockEmailSender implementa IEmailSender para testes.
type mockEmailSender struct {
	SendFn func(to, subject, htmlBody string) error
}

func (m *mockEmailSender) Send(to, subject, htmlBody string) error {
	return m.SendFn(to, subject, htmlBody)
}

func TestOtpService_Send(t *testing.T) {
	t.Run("happy path: e-mail sem OTP ativo", func(t *testing.T) {
		existingOTP := domain.OTP{
			ID:        uuid.New(),
			Email:     "user@example.com",
			Status:    utils.OTPStatusWaiting,
			ExpiresAt: time.Now().Add(30 * time.Minute),
		}

		repo := &mockOtpRepository{
			FindActiveByEmailFn: func(email string) (*domain.OTP, error) {
				return nil, nil
			},
			CreateFn: func(otp domain.OTP) (domain.OTP, error) {
				existingOTP.Code = otp.Code
				existingOTP.Email = otp.Email
				return existingOTP, nil
			},
		}
		emailSender := &mockEmailSender{
			SendFn: func(to, subject, htmlBody string) error {
				return nil
			},
		}

		svc := otpmodule.NewOtpService(repo, emailSender)
		dto := domain.SendOTPRequestDTO{Email: "user@example.com"}

		resp, err := svc.Send(dto)

		if err != nil {
			t.Fatalf("esperado sem erro, obteve: %v", err)
		}
		if resp.ID == (uuid.UUID{}) {
			t.Error("esperado ID não nulo na resposta")
		}
		if resp.Message != "OTP sent successfully" {
			t.Errorf("mensagem inesperada: %s", resp.Message)
		}
	})

	t.Run("happy path: e-mail com OTP em waiting é expirado antes de criar novo", func(t *testing.T) {
		existingID := uuid.New()
		newID := uuid.New()
		updateStatusCalled := false

		existingOTP := &domain.OTP{
			ID:        existingID,
			Email:     "user@example.com",
			Status:    utils.OTPStatusWaiting,
			ExpiresAt: time.Now().Add(10 * time.Minute),
		}

		repo := &mockOtpRepository{
			FindActiveByEmailFn: func(email string) (*domain.OTP, error) {
				return existingOTP, nil
			},
			UpdateStatusFn: func(id uuid.UUID, status string) error {
				if id == existingID && status == utils.OTPStatusExpired {
					updateStatusCalled = true
				}
				return nil
			},
			CreateFn: func(otp domain.OTP) (domain.OTP, error) {
				otp.ID = newID
				return otp, nil
			},
		}
		emailSender := &mockEmailSender{
			SendFn: func(to, subject, htmlBody string) error {
				return nil
			},
		}

		svc := otpmodule.NewOtpService(repo, emailSender)
		dto := domain.SendOTPRequestDTO{Email: "user@example.com"}

		resp, err := svc.Send(dto)

		if err != nil {
			t.Fatalf("esperado sem erro, obteve: %v", err)
		}
		if !updateStatusCalled {
			t.Error("UpdateStatus deveria ter sido chamado para expirar OTP anterior")
		}
		if resp.ID != newID {
			t.Error("esperado ID do novo OTP na resposta")
		}
	})

	t.Run("erro ao expirar OTP anterior", func(t *testing.T) {
		createCalled := false
		existingOTP := &domain.OTP{
			ID:        uuid.New(),
			Email:     "user@example.com",
			Status:    utils.OTPStatusWaiting,
			ExpiresAt: time.Now().Add(10 * time.Minute),
		}

		repo := &mockOtpRepository{
			FindActiveByEmailFn: func(email string) (*domain.OTP, error) {
				return existingOTP, nil
			},
			UpdateStatusFn: func(id uuid.UUID, status string) error {
				return errors.New("db error")
			},
			CreateFn: func(otp domain.OTP) (domain.OTP, error) {
				createCalled = true
				return otp, nil
			},
		}
		emailSender := &mockEmailSender{}

		svc := otpmodule.NewOtpService(repo, emailSender)
		dto := domain.SendOTPRequestDTO{Email: "user@example.com"}

		_, err := svc.Send(dto)

		if err == nil {
			t.Fatal("esperado erro, obteve nil")
		}
		if err.Error() != "error expiring previous otp" {
			t.Errorf("mensagem de erro inesperada: %v", err)
		}
		if createCalled {
			t.Error("Create não deve ser chamado quando expiracao falha")
		}
	})

	t.Run("erro ao criar OTP", func(t *testing.T) {
		emailSendCalled := false

		repo := &mockOtpRepository{
			FindActiveByEmailFn: func(email string) (*domain.OTP, error) {
				return nil, nil
			},
			CreateFn: func(otp domain.OTP) (domain.OTP, error) {
				return domain.OTP{}, errors.New("db error")
			},
		}
		emailSender := &mockEmailSender{
			SendFn: func(to, subject, htmlBody string) error {
				emailSendCalled = true
				return nil
			},
		}

		svc := otpmodule.NewOtpService(repo, emailSender)
		dto := domain.SendOTPRequestDTO{Email: "user@example.com"}

		_, err := svc.Send(dto)

		if err == nil {
			t.Fatal("esperado erro, obteve nil")
		}
		if err.Error() != "error creating otp" {
			t.Errorf("mensagem de erro inesperada: %v", err)
		}
		if emailSendCalled {
			t.Error("emailSender.Send não deve ser chamado quando criação falha")
		}
	})

	t.Run("erro ao enviar e-mail", func(t *testing.T) {
		newID := uuid.New()

		repo := &mockOtpRepository{
			FindActiveByEmailFn: func(email string) (*domain.OTP, error) {
				return nil, nil
			},
			CreateFn: func(otp domain.OTP) (domain.OTP, error) {
				otp.ID = newID
				return otp, nil
			},
		}
		emailSender := &mockEmailSender{
			SendFn: func(to, subject, htmlBody string) error {
				return errors.New("smtp error")
			},
		}

		svc := otpmodule.NewOtpService(repo, emailSender)
		dto := domain.SendOTPRequestDTO{Email: "user@example.com"}

		_, err := svc.Send(dto)

		if err == nil {
			t.Fatal("esperado erro, obteve nil")
		}
		if err.Error() != "error sending otp email" {
			t.Errorf("mensagem de erro inesperada: %v", err)
		}
	})
}

func TestOtpService_Confirm(t *testing.T) {
	t.Run("happy path: código correto, OTP waiting e não expirado", func(t *testing.T) {
		otpID := uuid.New()
		validOTP := &domain.OTP{
			ID:        otpID,
			Email:     "user@example.com",
			Code:      "1234",
			Status:    utils.OTPStatusWaiting,
			ExpiresAt: time.Now().Add(30 * time.Minute),
		}

		repo := &mockOtpRepository{
			FindByIDFn: func(id uuid.UUID) (*domain.OTP, error) {
				return validOTP, nil
			},
			UpdateStatusFn: func(id uuid.UUID, status string) error {
				return nil
			},
		}

		svc := otpmodule.NewOtpService(repo, &mockEmailSender{})
		dto := domain.ConfirmOTPRequestDTO{
			ID:   otpID.String(),
			Code: "1234",
		}

		resp, err := svc.Confirm(dto)

		if err != nil {
			t.Fatalf("esperado sem erro, obteve: %v", err)
		}
		if resp.Message != "OTP confirmed successfully" {
			t.Errorf("mensagem inesperada: %s", resp.Message)
		}
	})

	t.Run("OTP não encontrado", func(t *testing.T) {
		repo := &mockOtpRepository{
			FindByIDFn: func(id uuid.UUID) (*domain.OTP, error) {
				return nil, utils.ErrOTPNotFound
			},
		}

		svc := otpmodule.NewOtpService(repo, &mockEmailSender{})
		dto := domain.ConfirmOTPRequestDTO{
			ID:   uuid.New().String(),
			Code: "1234",
		}

		_, err := svc.Confirm(dto)

		if err == nil {
			t.Fatal("esperado erro, obteve nil")
		}
		if !errors.Is(err, utils.ErrOTPNotFound) {
			t.Errorf("esperado ErrOTPNotFound, obteve: %v", err)
		}
	})

	t.Run("ID UUID inválido retorna ErrOTPNotFound", func(t *testing.T) {
		findByIDCalled := false
		repo := &mockOtpRepository{
			FindByIDFn: func(id uuid.UUID) (*domain.OTP, error) {
				findByIDCalled = true
				return nil, nil
			},
		}

		svc := otpmodule.NewOtpService(repo, &mockEmailSender{})
		dto := domain.ConfirmOTPRequestDTO{
			ID:   "not-a-uuid",
			Code: "1234",
		}

		_, err := svc.Confirm(dto)

		if err == nil {
			t.Fatal("esperado erro, obteve nil")
		}
		if !errors.Is(err, utils.ErrOTPNotFound) {
			t.Errorf("esperado ErrOTPNotFound, obteve: %v", err)
		}
		if findByIDCalled {
			t.Error("FindByID não deve ser chamado com UUID inválido")
		}
	})

	t.Run("OTP com status confirmed retorna ErrOTPNotWaiting", func(t *testing.T) {
		otpID := uuid.New()
		confirmedOTP := &domain.OTP{
			ID:        otpID,
			Email:     "user@example.com",
			Code:      "1234",
			Status:    utils.OTPStatusConfirmed,
			ExpiresAt: time.Now().Add(30 * time.Minute),
		}

		repo := &mockOtpRepository{
			FindByIDFn: func(id uuid.UUID) (*domain.OTP, error) {
				return confirmedOTP, nil
			},
		}

		svc := otpmodule.NewOtpService(repo, &mockEmailSender{})
		dto := domain.ConfirmOTPRequestDTO{
			ID:   otpID.String(),
			Code: "1234",
		}

		_, err := svc.Confirm(dto)

		if err == nil {
			t.Fatal("esperado erro, obteve nil")
		}
		if !errors.Is(err, utils.ErrOTPNotWaiting) {
			t.Errorf("esperado ErrOTPNotWaiting, obteve: %v", err)
		}
	})

	t.Run("OTP com status expired retorna ErrOTPNotWaiting", func(t *testing.T) {
		otpID := uuid.New()
		expiredStatusOTP := &domain.OTP{
			ID:        otpID,
			Email:     "user@example.com",
			Code:      "1234",
			Status:    utils.OTPStatusExpired,
			ExpiresAt: time.Now().Add(30 * time.Minute),
		}

		repo := &mockOtpRepository{
			FindByIDFn: func(id uuid.UUID) (*domain.OTP, error) {
				return expiredStatusOTP, nil
			},
		}

		svc := otpmodule.NewOtpService(repo, &mockEmailSender{})
		dto := domain.ConfirmOTPRequestDTO{
			ID:   otpID.String(),
			Code: "1234",
		}

		_, err := svc.Confirm(dto)

		if err == nil {
			t.Fatal("esperado erro, obteve nil")
		}
		if !errors.Is(err, utils.ErrOTPNotWaiting) {
			t.Errorf("esperado ErrOTPNotWaiting, obteve: %v", err)
		}
	})

	t.Run("OTP waiting mas expirado por tempo retorna ErrOTPExpired e chama UpdateStatus", func(t *testing.T) {
		otpID := uuid.New()
		updateStatusCalled := false
		timeExpiredOTP := &domain.OTP{
			ID:        otpID,
			Email:     "user@example.com",
			Code:      "1234",
			Status:    utils.OTPStatusWaiting,
			ExpiresAt: time.Now().Add(-1 * time.Minute), // expirado
		}

		repo := &mockOtpRepository{
			FindByIDFn: func(id uuid.UUID) (*domain.OTP, error) {
				return timeExpiredOTP, nil
			},
			UpdateStatusFn: func(id uuid.UUID, status string) error {
				if id == otpID && status == utils.OTPStatusExpired {
					updateStatusCalled = true
				}
				return nil
			},
		}

		svc := otpmodule.NewOtpService(repo, &mockEmailSender{})
		dto := domain.ConfirmOTPRequestDTO{
			ID:   otpID.String(),
			Code: "1234",
		}

		_, err := svc.Confirm(dto)

		if err == nil {
			t.Fatal("esperado erro, obteve nil")
		}
		if !errors.Is(err, utils.ErrOTPExpired) {
			t.Errorf("esperado ErrOTPExpired, obteve: %v", err)
		}
		if !updateStatusCalled {
			t.Error("UpdateStatus deveria ter sido chamado para marcar OTP como expirado")
		}
	})

	t.Run("código incorreto retorna ErrOTPInvalid", func(t *testing.T) {
		otpID := uuid.New()
		validOTP := &domain.OTP{
			ID:        otpID,
			Email:     "user@example.com",
			Code:      "1234",
			Status:    utils.OTPStatusWaiting,
			ExpiresAt: time.Now().Add(30 * time.Minute),
		}

		repo := &mockOtpRepository{
			FindByIDFn: func(id uuid.UUID) (*domain.OTP, error) {
				return validOTP, nil
			},
		}

		svc := otpmodule.NewOtpService(repo, &mockEmailSender{})
		dto := domain.ConfirmOTPRequestDTO{
			ID:   otpID.String(),
			Code: "9999",
		}

		_, err := svc.Confirm(dto)

		if err == nil {
			t.Fatal("esperado erro, obteve nil")
		}
		if !errors.Is(err, utils.ErrOTPInvalid) {
			t.Errorf("esperado ErrOTPInvalid, obteve: %v", err)
		}
	})

	t.Run("erro ao atualizar status para confirmed propaga erro", func(t *testing.T) {
		otpID := uuid.New()
		validOTP := &domain.OTP{
			ID:        otpID,
			Email:     "user@example.com",
			Code:      "1234",
			Status:    utils.OTPStatusWaiting,
			ExpiresAt: time.Now().Add(30 * time.Minute),
		}

		repo := &mockOtpRepository{
			FindByIDFn: func(id uuid.UUID) (*domain.OTP, error) {
				return validOTP, nil
			},
			UpdateStatusFn: func(id uuid.UUID, status string) error {
				return errors.New("db error")
			},
		}

		svc := otpmodule.NewOtpService(repo, &mockEmailSender{})
		dto := domain.ConfirmOTPRequestDTO{
			ID:   otpID.String(),
			Code: "1234",
		}

		_, err := svc.Confirm(dto)

		if err == nil {
			t.Fatal("esperado erro, obteve nil")
		}
		if err.Error() != "error confirming otp" {
			t.Errorf("mensagem de erro inesperada: %v", err)
		}
	})
}
