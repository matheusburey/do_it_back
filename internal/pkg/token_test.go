package pkg

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestGenerateAndValidate(t *testing.T) {
	secret := "test-secret"
	userID := uuid.New()

	token, err := GenerateAccessToken(userID, secret, time.Hour)
	if err != nil {
		t.Fatalf("GenerateAccessToken: erro inesperado: %v", err)
	}

	claims, err := Validate(token, secret)
	if err != nil {
		t.Fatalf("Validate: erro inesperado: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("UserID = %v, quero %v", claims.UserID, userID)
	}
}

func TestValidateExpired(t *testing.T) {
	secret := "secret"

	token, err := GenerateAccessToken(uuid.New(), secret, -time.Hour)

	if err != nil {
		t.Fatalf("GenerateAccessToken: %v", err)
	}

	_, err = Validate(token, secret)

	if err == nil {
		t.Errorf("esperava erro de token expirado, veio %v", err)
	}
}
