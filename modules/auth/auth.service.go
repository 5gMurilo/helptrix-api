package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/5gMurilo/helptrix-api/core/domain"
	authinterfaces "github.com/5gMurilo/helptrix-api/core/interfaces/auth"
	"github.com/5gMurilo/helptrix-api/core/logger"
	"github.com/5gMurilo/helptrix-api/core/utils"
)

type AuthService struct {
	repo       authinterfaces.IAuthRepository
	tokenMaker authinterfaces.ITokenMaker
}

func NewAuthService(repo authinterfaces.IAuthRepository, tokenMaker authinterfaces.ITokenMaker) authinterfaces.IAuthService {
	return &AuthService{repo: repo, tokenMaker: tokenMaker}
}

func (s *AuthService) Register(dto domain.RegisterRequestDTO) (domain.RegisterResponseDTO, error) {
	log := logger.Get().With(
		slog.String("layer", "service"),
		slog.String("module", "auth"),
		slog.String("operation", "Register"),
	)

	if dto.UserType != utils.UserTypeHelper && dto.UserType != utils.UserTypeBusiness {
		log.Warn("registration rejected: invalid user_type", slog.String("user_type", dto.UserType))
		return domain.RegisterResponseDTO{}, errors.New("invalid user_type: must be 'helper' or 'business'")
	}

	if len(dto.Categories) > utils.MaxCategoriesPerUserRegistration {
		log.Warn("registration rejected: too many categories",
			slog.Int("provided", len(dto.Categories)),
			slog.Int("max", utils.MaxCategoriesPerUserRegistration),
		)
		return domain.RegisterResponseDTO{}, fmt.Errorf("invalid categories: at most %d categories allowed", utils.MaxCategoriesPerUserRegistration)
	}

	if dto.UserType == utils.UserTypeHelper && !utils.CPFRegex.MatchString(dto.Document) {
		log.Warn("registration rejected: invalid CPF format")
		return domain.RegisterResponseDTO{}, errors.New("invalid CPF: must contain exactly 11 numeric digits")
	}

	if dto.UserType == utils.UserTypeBusiness && !utils.CNPJRegex.MatchString(dto.Document) {
		log.Warn("registration rejected: invalid CNPJ format")
		return domain.RegisterResponseDTO{}, errors.New("invalid CNPJ: must contain exactly 14 numeric digits")
	}

	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		log.Error("failed to generate password salt", slog.String("error", err.Error()))
		return domain.RegisterResponseDTO{}, errors.New("error generating password salt")
	}

	saltHex := hex.EncodeToString(salt)
	hash := sha256.Sum256([]byte(saltHex + dto.Password))
	hashHex := hex.EncodeToString(hash[:])
	hashedPassword := saltHex + ":" + hashHex

	user, categoryIDs, err := s.repo.Register(dto, hashedPassword)
	if err != nil {
		log.Error("user registration failed", slog.String("error", err.Error()), slog.String("email", dto.Email))
		return domain.RegisterResponseDTO{}, err
	}

	log.Info("user registered successfully",
		slog.String("user_id", user.ID.String()),
		slog.String("user_type", user.UserType),
	)
	return domain.RegisterResponseDTO{
		ID:         user.ID,
		Name:       user.Name,
		Email:      user.Email,
		UserType:   user.UserType,
		Categories: categoryIDs,
		CreatedAt:  user.CreatedAt,
	}, nil
}

func (s *AuthService) Login(dto domain.LoginRequestDTO) (domain.LoginResponseDTO, error) {
	log := logger.Get().With(
		slog.String("layer", "service"),
		slog.String("module", "auth"),
		slog.String("operation", "Login"),
	)

	user, err := s.repo.FindByEmail(dto.Email)
	if err != nil {
		log.Warn("login failed: email not found", slog.String("email", dto.Email))
		return domain.LoginResponseDTO{}, utils.ErrInvalidCredentials
	}

	parts := strings.SplitN(user.Password, ":", 2)
	if len(parts) != 2 {
		log.Error("corrupted password hash in database", slog.String("user_id", user.ID.String()))
		return domain.LoginResponseDTO{}, utils.ErrInvalidCredentials
	}

	hash := sha256.Sum256([]byte(parts[0] + dto.Password))
	if hex.EncodeToString(hash[:]) != parts[1] {
		log.Warn("login failed: wrong password", slog.String("email", dto.Email))
		return domain.LoginResponseDTO{}, utils.ErrInvalidCredentials
	}

	token, err := s.tokenMaker.CreateToken(user.ID.String(), user.Name, user.Email, user.UserType, 8*time.Hour)
	if err != nil {
		log.Error("token creation failed", slog.String("error", err.Error()), slog.String("user_id", user.ID.String()))
		return domain.LoginResponseDTO{}, utils.ErrInvalidCredentials
	}

	log.Info("user logged in successfully",
		slog.String("user_id", user.ID.String()),
		slog.String("user_type", user.UserType),
	)
	return domain.LoginResponseDTO{ID: user.ID, Token: token}, nil
}
