package auth_test

import (
	"testing"
	"time"

	pasetoauth "github.com/5gMurilo/helptrix-api/adapter/auth"
)

const testSymmetricKey = "3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d5e6f7a8b9c0d1e2f3a4b"

func newTestMaker(t *testing.T) *pasetoauth.PasetoMaker {
	t.Helper()
	maker, err := pasetoauth.NewPasetoMaker(testSymmetricKey)
	if err != nil {
		t.Fatalf("failed to create PasetoMaker: %v", err)
	}
	return maker
}

func TestNewPasetoMaker(t *testing.T) {
	t.Run("valid 64-char hex key", func(t *testing.T) {
		// Arrange
		key := testSymmetricKey

		// Act
		maker, err := pasetoauth.NewPasetoMaker(key)

		// Assert
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if maker == nil {
			t.Fatal("expected non-nil maker")
		}
	})

	t.Run("invalid key: too short", func(t *testing.T) {
		// Arrange
		key := "short"

		// Act
		_, err := pasetoauth.NewPasetoMaker(key)

		// Assert
		if err == nil {
			t.Fatal("expected error for short key, got nil")
		}
	})

	t.Run("invalid key: non-hex characters", func(t *testing.T) {
		// Arrange
		key := "zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz"

		// Act
		_, err := pasetoauth.NewPasetoMaker(key)

		// Assert
		if err == nil {
			t.Fatal("expected error for non-hex key, got nil")
		}
	})
}

func TestPasetoMaker_CreateToken(t *testing.T) {
	t.Run("creates non-empty token", func(t *testing.T) {
		// Arrange
		maker := newTestMaker(t)

		// Act
		token, err := maker.CreateToken("user-123", "John", "john@example.com", "helper", 8*time.Hour)

		// Assert
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if token == "" {
			t.Error("expected non-empty token")
		}
	})
}

func TestPasetoMaker_VerifyToken(t *testing.T) {
	t.Run("valid token round-trip", func(t *testing.T) {
		// Arrange
		maker := newTestMaker(t)
		userID := "user-abc"
		email := "john@example.com"
		token, err := maker.CreateToken(userID, "John", email, "helper", 8*time.Hour)
		if err != nil {
			t.Fatalf("failed to create token: %v", err)
		}

		// Act
		payload, err := maker.VerifyToken(token)

		// Assert
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if payload.UserID != userID {
			t.Errorf("expected UserID %q, got %q", userID, payload.UserID)
		}
		if payload.Email != email {
			t.Errorf("expected Email %q, got %q", email, payload.Email)
		}
	})

	t.Run("tampered token returns ErrInvalidToken", func(t *testing.T) {
		// Arrange
		maker := newTestMaker(t)

		// Act
		_, err := maker.VerifyToken("not.a.valid.token")

		// Assert
		if err == nil {
			t.Fatal("expected error for tampered token, got nil")
		}
		if err != pasetoauth.ErrInvalidToken {
			t.Errorf("expected ErrInvalidToken, got: %v", err)
		}
	})

	t.Run("expired token returns ErrExpiredToken", func(t *testing.T) {
		// Arrange
		maker := newTestMaker(t)
		token, err := maker.CreateToken("user-xyz", "John", "john@example.com", "helper", -1*time.Second)
		if err != nil {
			t.Fatalf("failed to create expired token: %v", err)
		}

		// Act
		_, err = maker.VerifyToken(token)

		// Assert
		if err == nil {
			t.Fatal("expected error for expired token, got nil")
		}
		if err != pasetoauth.ErrExpiredToken {
			t.Errorf("expected ErrExpiredToken, got: %v", err)
		}
	})
}

func TestPayload_Valid(t *testing.T) {
	t.Run("not expired", func(t *testing.T) {
		// Arrange
		p := &pasetoauth.Payload{ExpiredAt: time.Now().Add(time.Hour)}

		// Act
		err := p.Valid()

		// Assert
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
	})

	t.Run("already expired", func(t *testing.T) {
		// Arrange
		p := &pasetoauth.Payload{ExpiredAt: time.Now().Add(-time.Second)}

		// Act
		err := p.Valid()

		// Assert
		if err == nil {
			t.Fatal("expected error for expired payload, got nil")
		}
		if err != pasetoauth.ErrExpiredToken {
			t.Errorf("expected ErrExpiredToken, got: %v", err)
		}
	})
}
