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
	authmodule "github.com/5gMurilo/helptrix-api/modules/auth"
	"github.com/5gMurilo/helptrix-api/core/utils"
	"github.com/google/uuid"
)

// mockAuthRepository implementa IAuthRepository para testes.
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

// mockTokenMaker implementa ITokenMaker para testes.
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

// buildStoredPassword computa um saltHex:hashHex deterministico usando salt fixo para testes.
func buildStoredPassword(password string) string {
	saltHex := "aabbccddeeff00112233445566778899"
	hash := sha256.Sum256([]byte(saltHex + password))
	hashHex := hex.EncodeToString(hash[:])
	return saltHex + ":" + hashHex
}

func validHelperDTO() domain.RegisterRequestDTO {
	return domain.RegisterRequestDTO{
		Name:       "João da Silva",
		Email:      "joao@example.com",
		Password:   "senha123",
		UserType:   utils.UserTypeHelper,
		Document:   "12345678901",
		Phone:      "11999999999",
		Categories: []uint{1, 2},
		Address: domain.AddressInputDTO{
			Street:       "Rua das Flores",
			Number:       "100",
			Neighborhood: "Centro",
			City:         "São Paulo",
			State:        "SP",
		},
	}
}

func validBusinessDTO() domain.RegisterRequestDTO {
	return domain.RegisterRequestDTO{
		Name:       "Empresa Ltda",
		Email:      "empresa@example.com",
		Password:   "senha123",
		UserType:   utils.UserTypeBusiness,
		Document:   "12345678000195",
		Phone:      "11888888888",
		Categories: []uint{3},
		Address: domain.AddressInputDTO{
			Street:       "Av. Paulista",
			Number:       "1000",
			Neighborhood: "Bela Vista",
			City:         "São Paulo",
			State:        "SP",
		},
	}
}

func TestAuthService_Register(t *testing.T) {
	t.Run("happy path: helper com CPF válido de 11 dígitos", func(t *testing.T) {
		mockUser := domain.User{
			ID:        uuid.New(),
			Name:      "João da Silva",
			Email:     "joao@example.com",
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
		dto := validHelperDTO()

		resp, err := svc.Register(dto)

		if err != nil {
			t.Fatalf("esperado sem erro, obteve: %v", err)
		}
		if resp.ID == (uuid.UUID{}) {
			t.Error("esperado ID não nulo na resposta")
		}
		if len(resp.Categories) != 2 {
			t.Errorf("esperado 2 categorias, obteve: %d", len(resp.Categories))
		}
	})

	t.Run("happy path: business com CNPJ válido de 14 dígitos", func(t *testing.T) {
		mockUser := domain.User{
			ID:        uuid.New(),
			Name:      "Empresa Ltda",
			Email:     "empresa@example.com",
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
		dto := validBusinessDTO()

		resp, err := svc.Register(dto)

		if err != nil {
			t.Fatalf("esperado sem erro, obteve: %v", err)
		}
		if resp.ID == (uuid.UUID{}) {
			t.Error("esperado ID não nulo na resposta")
		}
		if len(resp.Categories) != 1 {
			t.Errorf("esperado 1 categoria, obteve: %d", len(resp.Categories))
		}
	})

	t.Run("user_type inválido", func(t *testing.T) {
		repoCalled := false
		repo := &mockAuthRepository{
			RegisterFn: func(dto domain.RegisterRequestDTO, hashedPassword string) (domain.User, []uint, error) {
				repoCalled = true
				return domain.User{}, nil, nil
			},
		}

		svc := authmodule.NewAuthService(repo, &mockTokenMaker{})
		dto := validHelperDTO()
		dto.UserType = "admin"

		_, err := svc.Register(dto)

		if err == nil {
			t.Fatal("esperado erro para user_type inválido, obteve nil")
		}
		if !strings.Contains(err.Error(), "invalid user_type") {
			t.Errorf("mensagem de erro inesperada: %v", err)
		}
		if repoCalled {
			t.Error("o repositório não deve ser chamado quando user_type é inválido")
		}
	})

	t.Run("helper com documento inválido (CPF deve ter 11 dígitos)", func(t *testing.T) {
		repo := &mockAuthRepository{
			RegisterFn: func(dto domain.RegisterRequestDTO, hashedPassword string) (domain.User, []uint, error) {
				return domain.User{}, nil, nil
			},
		}

		svc := authmodule.NewAuthService(repo, &mockTokenMaker{})
		dto := validHelperDTO()
		dto.Document = "1234"

		_, err := svc.Register(dto)

		if err == nil {
			t.Fatal("esperado erro para CPF inválido, obteve nil")
		}
		if !strings.Contains(err.Error(), "invalid CPF") {
			t.Errorf("mensagem de erro inesperada: %v", err)
		}
	})

	t.Run("business com documento inválido (CNPJ deve ter 14 dígitos, não 11)", func(t *testing.T) {
		repo := &mockAuthRepository{
			RegisterFn: func(dto domain.RegisterRequestDTO, hashedPassword string) (domain.User, []uint, error) {
				return domain.User{}, nil, nil
			},
		}

		svc := authmodule.NewAuthService(repo, &mockTokenMaker{})
		dto := validBusinessDTO()
		dto.Document = "12345678901"

		_, err := svc.Register(dto)

		if err == nil {
			t.Fatal("esperado erro para CNPJ inválido, obteve nil")
		}
		if !strings.Contains(err.Error(), "invalid CNPJ") {
			t.Errorf("mensagem de erro inesperada: %v", err)
		}
	})

	t.Run("usuário duplicado: repositório retorna ErrUserAlreadyRegistered", func(t *testing.T) {
		repo := &mockAuthRepository{
			RegisterFn: func(dto domain.RegisterRequestDTO, hashedPassword string) (domain.User, []uint, error) {
				return domain.User{}, nil, utils.ErrUserAlreadyRegistered
			},
		}

		svc := authmodule.NewAuthService(repo, &mockTokenMaker{})
		dto := validHelperDTO()

		_, err := svc.Register(dto)

		if err == nil {
			t.Fatal("esperado erro de usuário duplicado, obteve nil")
		}
		if !errors.Is(err, utils.ErrUserAlreadyRegistered) {
			t.Errorf("esperado ErrUserAlreadyRegistered, obteve: %v", err)
		}
	})
}

func TestAuthService_Login(t *testing.T) {
	t.Run("happy path: credenciais válidas retornam token sem erro", func(t *testing.T) {
		mockUser := &domain.User{
			ID:       uuid.New(),
			Name:     "João da Silva",
			Email:    "joao@example.com",
			Password: buildStoredPassword("senha123"),
			UserType: utils.UserTypeHelper,
		}

		repo := &mockAuthRepository{
			FindByEmailFn: func(email string) (*domain.User, error) {
				return mockUser, nil
			},
		}
		tokenMaker := &mockTokenMaker{
			CreateTokenFn: func(userID, name, email, userType string, duration time.Duration) (string, error) {
				return "token-valido", nil
			},
		}

		svc := authmodule.NewAuthService(repo, tokenMaker)
		dto := domain.LoginRequestDTO{
			Email:    "joao@example.com",
			Password: "senha123",
		}

		resp, err := svc.Login(dto)

		if err != nil {
			t.Fatalf("esperado sem erro, obteve: %v", err)
		}
		if resp.Token != "token-valido" {
			t.Errorf("esperado token 'token-valido', obteve: %s", resp.Token)
		}
	})

	t.Run("email não encontrado: FindByEmail retorna ErrUserNotFound → ErrInvalidCredentials propagado", func(t *testing.T) {
		repo := &mockAuthRepository{
			FindByEmailFn: func(email string) (*domain.User, error) {
				return nil, utils.ErrUserNotFound
			},
		}

		svc := authmodule.NewAuthService(repo, &mockTokenMaker{})
		dto := domain.LoginRequestDTO{
			Email:    "naoexiste@example.com",
			Password: "qualquersenha",
		}

		_, err := svc.Login(dto)

		if err == nil {
			t.Fatal("esperado erro, obteve nil")
		}
		if !errors.Is(err, utils.ErrInvalidCredentials) {
			t.Errorf("esperado ErrInvalidCredentials, obteve: %v", err)
		}
	})

	t.Run("senha incorreta: usuário encontrado mas senha diverge → ErrInvalidCredentials", func(t *testing.T) {
		mockUser := &domain.User{
			ID:       uuid.New(),
			Name:     "João da Silva",
			Email:    "joao@example.com",
			Password: buildStoredPassword("senha123"),
			UserType: utils.UserTypeHelper,
		}

		repo := &mockAuthRepository{
			FindByEmailFn: func(email string) (*domain.User, error) {
				return mockUser, nil
			},
		}

		svc := authmodule.NewAuthService(repo, &mockTokenMaker{})
		dto := domain.LoginRequestDTO{
			Email:    "joao@example.com",
			Password: "senhaerrada",
		}

		_, err := svc.Login(dto)

		if err == nil {
			t.Fatal("esperado erro para senha incorreta, obteve nil")
		}
		if !errors.Is(err, utils.ErrInvalidCredentials) {
			t.Errorf("esperado ErrInvalidCredentials, obteve: %v", err)
		}
	})
}
