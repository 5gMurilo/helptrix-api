package domain

import (
	"time"

	"github.com/google/uuid"
)

// LoginResponseDTO representa a resposta apos autenticacao bem-sucedida.
//
//	@name	LoginResponseDTO
type LoginResponseDTO struct {
	Token string `json:"token"`
}

// RegisterResponseDTO represents the response payload after successful user registration.
//
//	@name	RegisterResponseDTO
type RegisterResponseDTO struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	UserType   string    `json:"user_type"`
	Categories []uint    `json:"categories"`
	CreatedAt  time.Time `json:"created_at"`
}
