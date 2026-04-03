package auth

import (
	"errors"
	"time"

	gopasseto "aidanwoods.dev/go-paseto"
	"github.com/google/uuid"
)

const (
	minSymmetricKeySize = 32
	TokenDuration       = 8 * time.Hour
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

type Payload struct {
	UserID    string
	Name      string
	Email     string
	UserType  string
	IssuedAt  time.Time
	ExpiredAt time.Time
}

type PasetoMaker struct {
	symmetricKey gopasseto.V2SymmetricKey
}

func NewPasetoMaker(symmetricKeyHex string) (*PasetoMaker, error) {
	key, err := gopasseto.V2SymmetricKeyFromHex(symmetricKeyHex)
	if err != nil {
		if len(symmetricKeyHex) < minSymmetricKeySize {
			return nil, errors.New("invalid key size: must be at least 32 characters")
		}
		return nil, err
	}

	return &PasetoMaker{symmetricKey: key}, nil
}

func (m *PasetoMaker) CreateToken(userID, name, email, userType string, duration time.Duration) (string, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	now := time.Now()
	expiredAt := now.Add(duration)

	token := gopasseto.NewToken()
	token.SetString("id", tokenID.String())
	token.SetString("user_id", userID)
	token.SetString("name", name)
	token.SetString("email", email)
	token.SetString("user_type", userType)
	token.SetIssuedAt(now)
	token.SetExpiration(expiredAt)

	return token.V2Encrypt(m.symmetricKey), nil
}

func (m *PasetoMaker) VerifyToken(tokenStr string) (*Payload, error) {
	parser := gopasseto.NewParserWithoutExpiryCheck()

	token, err := parser.ParseV2Local(m.symmetricKey, tokenStr)
	if err != nil {
		return nil, ErrInvalidToken
	}

	userID, err := token.GetString("user_id")
	if err != nil {
		return nil, ErrInvalidToken
	}

	name, err := token.GetString("name")
	if err != nil {
		return nil, ErrInvalidToken
	}

	email, err := token.GetString("email")
	if err != nil {
		return nil, ErrInvalidToken
	}

	issuedAt, err := token.GetIssuedAt()
	if err != nil {
		return nil, ErrInvalidToken
	}

	expiredAt, err := token.GetExpiration()
	if err != nil {
		return nil, ErrInvalidToken
	}

	userType, err := token.GetString("user_type")
	if err != nil {
		return nil, ErrInvalidToken
	}

	payload := &Payload{
		UserID:    userID,
		Name:      name,
		Email:     email,
		UserType:  userType,
		IssuedAt:  issuedAt,
		ExpiredAt: expiredAt,
	}

	if err := payload.Valid(); err != nil {
		return nil, err
	}

	return payload, nil
}

func (p *Payload) Valid() error {
	if time.Now().After(p.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
