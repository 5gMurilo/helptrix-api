package auth_test

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"testing"
	"time"

	pasetoauth "github.com/5gMurilo/helptrix-api/adapter/auth"
	"github.com/5gMurilo/helptrix-api/core/domain"
	"github.com/5gMurilo/helptrix-api/core/utils"
	authmodule "github.com/5gMurilo/helptrix-api/modules/auth"
	"github.com/google/uuid"
)

type mockAuthRepository struct {
	RegisterFn    func(dto domain.RegisterRequestDTO, hashedPassword string) (domain.User, []uint, error)
	FindByEmailFn func(email string) (*domain.User, error)
}

func (m *mockAuthRepository) Register(dto domain.RegisterRequestDTO, hashedPassword string) (domain.User, []uint, error) {
	return m.RegisterFn(dto, hashedPassword)
}

func (m *mockAuthRepository) FindByEmail(email string) (*domain.User, error) {
	return m.FindByEmailFn(email)
}

type mockTokenMaker struct {
	CreateTokenFn func(userID, name, email, userType string, duration time.Duration) (string, error)
	VerifyTokenFn func(token string) (*pasetoauth.Payload, error)
}

func (m *mockTokenMaker) CreateToken(userID, name, email, userType string, duration time.Duration) (string, error) {
	return m.CreateTokenFn(userID, name, email, userType, duration)
}

func (m *mockTokenMaker) VerifyToken(token string) (*pasetoauth.Payload, error) {
	return m.VerifyTokenFn(token)
}

func buildStoredPassword(t *testing.T, password string) string {
	t.Helper()
	saltHex := "aabbccddeeff00112233445566778899"
	hash := sha256.Sum256([]byte(saltHex + password))
	hashHex := hex.EncodeToString(hash[:])
	return saltHex + ":" + hashHex
}

func validHelperDTO(t *testing.T) domain.RegisterRequestDTO {
	t.Helper()
	return domain.RegisterRequestDTO{
		Name:       "John Silva",
		Email:      "john@example.com",
		Password:   "password123",
		UserType:   utils.UserTypeHelper,
		Document:   "12345678901",
		Phone:      "11999999999",
		Categories: []uint{1, 2},
		Address: domain.AddressInputDTO{
			Street:       "Flower Street",
			Number:       "100",
			Neighborhood: "Downtown",
			City:         "Sao Paulo",
			State:        "SP",
		},
	}
}

func validBusinessDTO(t *testing.T) domain.RegisterRequestDTO {
	t.Helper()
	return domain.RegisterRequestDTO{
		Name:       "Company Ltd",
		Email:      "company@example.com",
		Password:   "password123",
		UserType:   utils.UserTypeBusiness,
		Document:   "12345678000195",
		Phone:      "11888888888",
		Categories: []uint{3},
		Address: domain.AddressInputDTO{
			Street:       "Paulista Ave",
			Number:       "1000",
			Neighborhood: "Bela Vista",
			City:         "Sao Paulo",
			State:        "SP",
		},
	}
}

func TestAuthService_Register(t *testing.T) {
	t.Run("happy path: helper with valid 11-digit CPF", func(t *testing.T) {
		// Arrange
		mockUser := domain.User{
			ID:        uuid.New(),
			Name:      "John Silva",
			Email:     "john@example.com",
			UserType:  utils.UserTypeHelper,
			CreatedAt: time.Now(),
		}
		mockCategories := []uint{1, 2}
		repo := &mockAuthRepository{
			RegisterFn: func(dto domain.RegisterRequestDTO, hashedPassword string) (domain.User, []uint, error) {
				return mockUser, mockCategories, nil
			},
		}
		svc := authmodule.NewAuthService(repo, &mockTokenMaker{})
		dto := validHelperDTO(t)

		// Act
		resp, err := svc.Register(dto)

		// Assert
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if resp.ID == (uuid.UUID{}) {
			t.Error("expected non-nil ID in response")
		}
		if len(resp.Categories) != 2 {
			t.Errorf("expected 2 categories, got: %d", len(resp.Categories))
		}
	})

	t.Run("happy path: business with valid 14-digit CNPJ", func(t *testing.T) {
		// Arrange
		mockUser := domain.User{
			ID:        uuid.New(),
			Name:      "Company Ltd",
			Email:     "company@example.com",
			UserType:  utils.UserTypeBusiness,
			CreatedAt: time.Now(),
		}
		mockCategories := []uint{3}
		repo := &mockAuthRepository{
			RegisterFn: func(dto domain.RegisterRequestDTO, hashedPassword string) (domain.User, []uint, error) {
				return mockUser, mockCategories, nil
			},
		}
		svc := authmodule.NewAuthService(repo, &mockTokenMaker{})
		dto := validBusinessDTO(t)

		// Act
		resp, err := svc.Register(dto)

		// Assert
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if resp.ID == (uuid.UUID{}) {
			t.Error("expected non-nil ID in response")
		}
		if len(resp.Categories) != 1 {
			t.Errorf("expected 1 category, got: %d", len(resp.Categories))
		}
	})

	t.Run("invalid user_type", func(t *testing.T) {
		// Arrange
		repoCalled := false
		repo := &mockAuthRepository{
			RegisterFn: func(dto domain.RegisterRequestDTO, hashedPassword string) (domain.User, []uint, error) {
				repoCalled = true
				return domain.User{}, nil, nil
			},
		}
		svc := authmodule.NewAuthService(repo, &mockTokenMaker{})
		dto := validHelperDTO(t)
		dto.UserType = "admin"

		// Act
		_, err := svc.Register(dto)

		// Assert
		if err == nil {
			t.Fatal("expected error for invalid user_type, got nil")
		}
		if !strings.Contains(err.Error(), "invalid user_type") {
			t.Errorf("unexpected error message: %v", err)
		}
		if repoCalled {
			t.Error("repository must not be called when user_type is invalid")
		}
	})

	t.Run("more than 3 categories", func(t *testing.T) {
		// Arrange
		repoCalled := false
		repo := &mockAuthRepository{
			RegisterFn: func(dto domain.RegisterRequestDTO, hashedPassword string) (domain.User, []uint, error) {
				repoCalled = true
				return domain.User{}, nil, nil
			},
		}
		svc := authmodule.NewAuthService(repo, &mockTokenMaker{})
		dto := validHelperDTO(t)
		dto.Categories = []uint{1, 2, 3, 4}

		// Act
		_, err := svc.Register(dto)

		// Assert
		if err == nil {
			t.Fatal("expected error for more than 3 categories, got nil")
		}
		if !strings.Contains(err.Error(), "invalid categories") {
			t.Errorf("unexpected error message: %v", err)
		}
		if repoCalled {
			t.Error("repository must not be called when categories exceed limit")
		}
	})

	t.Run("helper with invalid CPF", func(t *testing.T) {
		// Arrange
		repo := &mockAuthRepository{
			RegisterFn: func(dto domain.RegisterRequestDTO, hashedPassword string) (domain.User, []uint, error) {
				return domain.User{}, nil, nil
			},
		}
		svc := authmodule.NewAuthService(repo, &mockTokenMaker{})
		dto := validHelperDTO(t)
		dto.Document = "1234"

		// Act
		_, err := svc.Register(dto)

		// Assert
		if err == nil {
			t.Fatal("expected error for invalid CPF, got nil")
		}
		if !strings.Contains(err.Error(), "invalid CPF") {
			t.Errorf("unexpected error message: %v", err)
		}
	})

	t.Run("business with invalid CNPJ", func(t *testing.T) {
		// Arrange
		repo := &mockAuthRepository{
			RegisterFn: func(dto domain.RegisterRequestDTO, hashedPassword string) (domain.User, []uint, error) {
				return domain.User{}, nil, nil
			},
		}
		svc := authmodule.NewAuthService(repo, &mockTokenMaker{})
		dto := validBusinessDTO(t)
		dto.Document = "12345678901"

		// Act
		_, err := svc.Register(dto)

		// Assert
		if err == nil {
			t.Fatal("expected error for invalid CNPJ, got nil")
		}
		if !strings.Contains(err.Error(), "invalid CNPJ") {
			t.Errorf("unexpected error message: %v", err)
		}
	})

	t.Run("duplicate user: repository returns ErrUserAlreadyRegistered", func(t *testing.T) {
		// Arrange
		repo := &mockAuthRepository{
			RegisterFn: func(dto domain.RegisterRequestDTO, hashedPassword string) (domain.User, []uint, error) {
				return domain.User{}, nil, utils.ErrUserAlreadyRegistered
			},
		}
		svc := authmodule.NewAuthService(repo, &mockTokenMaker{})
		dto := validHelperDTO(t)

		// Act
		_, err := svc.Register(dto)

		// Assert
		if err == nil {
			t.Fatal("expected duplicate user error, got nil")
		}
		if !errors.Is(err, utils.ErrUserAlreadyRegistered) {
			t.Errorf("expected ErrUserAlreadyRegistered, got: %v", err)
		}
	})
}

func TestAuthService_Login(t *testing.T) {
	t.Run("happy path: valid credentials return token", func(t *testing.T) {
		// Arrange
		mockUser := &domain.User{
			ID:       uuid.New(),
			Name:     "John Silva",
			Email:    "john@example.com",
			Password: buildStoredPassword(t, "password123"),
			UserType: utils.UserTypeHelper,
		}
		repo := &mockAuthRepository{
			FindByEmailFn: func(email string) (*domain.User, error) {
				return mockUser, nil
			},
		}
		tokenMaker := &mockTokenMaker{
			CreateTokenFn: func(userID, name, email, userType string, duration time.Duration) (string, error) {
				return "valid-token", nil
			},
		}
		svc := authmodule.NewAuthService(repo, tokenMaker)
		dto := domain.LoginRequestDTO{Email: "john@example.com", Password: "password123"}

		// Act
		resp, err := svc.Login(dto)

		// Assert
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if resp.Token != "valid-token" {
			t.Errorf("expected token 'valid-token', got: %s", resp.Token)
		}
	})

	t.Run("email not found: FindByEmail returns ErrUserNotFound", func(t *testing.T) {
		// Arrange
		repo := &mockAuthRepository{
			FindByEmailFn: func(email string) (*domain.User, error) {
				return nil, utils.ErrUserNotFound
			},
		}
		svc := authmodule.NewAuthService(repo, &mockTokenMaker{})
		dto := domain.LoginRequestDTO{Email: "ghost@example.com", Password: "anypassword"}

		// Act
		_, err := svc.Login(dto)

		// Assert
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, utils.ErrInvalidCredentials) {
			t.Errorf("expected ErrInvalidCredentials, got: %v", err)
		}
	})

	t.Run("wrong password: hash mismatch returns ErrInvalidCredentials", func(t *testing.T) {
		// Arrange
		mockUser := &domain.User{
			ID:       uuid.New(),
			Name:     "John Silva",
			Email:    "john@example.com",
			Password: buildStoredPassword(t, "password123"),
			UserType: utils.UserTypeHelper,
		}
		repo := &mockAuthRepository{
			FindByEmailFn: func(email string) (*domain.User, error) {
				return mockUser, nil
			},
		}
		svc := authmodule.NewAuthService(repo, &mockTokenMaker{})
		dto := domain.LoginRequestDTO{Email: "john@example.com", Password: "wrongpassword"}

		// Act
		_, err := svc.Login(dto)

		// Assert
		if err == nil {
			t.Fatal("expected error for wrong password, got nil")
		}
		if !errors.Is(err, utils.ErrInvalidCredentials) {
			t.Errorf("expected ErrInvalidCredentials, got: %v", err)
		}
	})

	t.Run("malformed stored password: missing colon separator", func(t *testing.T) {
		// Arrange
		mockUser := &domain.User{
			ID:       uuid.New(),
			Name:     "John Silva",
			Email:    "john@example.com",
			Password: "nocolonseparator",
			UserType: utils.UserTypeHelper,
		}
		repo := &mockAuthRepository{
			FindByEmailFn: func(email string) (*domain.User, error) {
				return mockUser, nil
			},
		}
		svc := authmodule.NewAuthService(repo, &mockTokenMaker{})
		dto := domain.LoginRequestDTO{Email: "john@example.com", Password: "password123"}

		// Act
		_, err := svc.Login(dto)

		// Assert
		if err == nil {
			t.Fatal("expected error for malformed stored password, got nil")
		}
		if !errors.Is(err, utils.ErrInvalidCredentials) {
			t.Errorf("expected ErrInvalidCredentials, got: %v", err)
		}
	})

	t.Run("CreateToken error propagated as ErrInvalidCredentials", func(t *testing.T) {
		// Arrange
		mockUser := &domain.User{
			ID:       uuid.New(),
			Name:     "John Silva",
			Email:    "john@example.com",
			Password: buildStoredPassword(t, "password123"),
			UserType: utils.UserTypeHelper,
		}
		repo := &mockAuthRepository{
			FindByEmailFn: func(email string) (*domain.User, error) {
				return mockUser, nil
			},
		}
		tokenMaker := &mockTokenMaker{
			CreateTokenFn: func(userID, name, email, userType string, duration time.Duration) (string, error) {
				return "", errors.New("token generation failed")
			},
		}
		svc := authmodule.NewAuthService(repo, tokenMaker)
		dto := domain.LoginRequestDTO{Email: "john@example.com", Password: "password123"}

		// Act
		_, err := svc.Login(dto)

		// Assert
		if err == nil {
			t.Fatal("expected error when CreateToken fails, got nil")
		}
		if !errors.Is(err, utils.ErrInvalidCredentials) {
			t.Errorf("expected ErrInvalidCredentials, got: %v", err)
		}
	})
}
