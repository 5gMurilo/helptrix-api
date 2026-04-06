package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/5gMurilo/helptrix-api/core/domain"
	authinterfaces "github.com/5gMurilo/helptrix-api/core/interfaces/auth"
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
	if dto.UserType != utils.UserTypeHelper && dto.UserType != utils.UserTypeBusiness {
		return domain.RegisterResponseDTO{}, errors.New("invalid user_type: must be 'helper' or 'business'")
	}

	if len(dto.Categories) > utils.MaxCategoriesPerUserRegistration {
		return domain.RegisterResponseDTO{}, fmt.Errorf("invalid categories: at most %d categories allowed", utils.MaxCategoriesPerUserRegistration)
	}

	if dto.UserType == utils.UserTypeHelper && !utils.CPFRegex.MatchString(dto.Document) {
		return domain.RegisterResponseDTO{}, errors.New("invalid CPF: must contain exactly 11 numeric digits")
	}

	if dto.UserType == utils.UserTypeBusiness && !utils.CNPJRegex.MatchString(dto.Document) {
		return domain.RegisterResponseDTO{}, errors.New("invalid CNPJ: must contain exactly 14 numeric digits")
	}

	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return domain.RegisterResponseDTO{}, errors.New("error generating password salt")
	}

	saltHex := hex.EncodeToString(salt)
	hash := sha256.Sum256([]byte(saltHex + dto.Password))
	hashHex := hex.EncodeToString(hash[:])
	hashedPassword := saltHex + ":" + hashHex

	user, categoryIDs, err := s.repo.Register(dto, hashedPassword)
	if err != nil {
		return domain.RegisterResponseDTO{}, err
	}

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
	user, err := s.repo.FindByEmail(dto.Email)
	if err != nil {
		return domain.LoginResponseDTO{}, utils.ErrInvalidCredentials
	}

	parts := strings.SplitN(user.Password, ":", 2)
	if len(parts) != 2 {
		return domain.LoginResponseDTO{}, utils.ErrInvalidCredentials
	}

	hash := sha256.Sum256([]byte(parts[0] + dto.Password))
	if hex.EncodeToString(hash[:]) != parts[1] {
		return domain.LoginResponseDTO{}, utils.ErrInvalidCredentials
	}

	token, err := s.tokenMaker.CreateToken(user.ID.String(), user.Name, user.Email, user.UserType, 8*time.Hour)
	if err != nil {
		return domain.LoginResponseDTO{}, utils.ErrInvalidCredentials
	}

	return domain.LoginResponseDTO{ID: user.ID, Token: token}, nil
}
