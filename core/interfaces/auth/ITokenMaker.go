package authinterfaces

import (
	"time"

	"github.com/5gMurilo/helptrix-api/adapter/auth"
)

// ITokenMaker define o contrato para criacao e verificacao de tokens de autenticacao.
type ITokenMaker interface {
	CreateToken(userID, name, email, userType string, duration time.Duration) (string, error)
	VerifyToken(token string) (*auth.Payload, error)
}
