package authinterfaces

import "github.com/5gMurilo/helptrix-api/core/domain"

type IAuthRepository interface {
	Register(dto domain.RegisterRequestDTO, hashedPassword string) (domain.User, []uint, error)
	FindByEmail(email string) (*domain.User, error)
}
