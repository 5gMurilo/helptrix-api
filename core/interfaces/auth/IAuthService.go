package authinterfaces

import "github.com/5gMurilo/helptrix-api/core/domain"

type IAuthService interface {
	Register(dto domain.RegisterRequestDTO) (domain.RegisterResponseDTO, error)
	Login(dto domain.LoginRequestDTO) (domain.LoginResponseDTO, error)
}
