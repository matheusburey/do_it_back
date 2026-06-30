package pkg

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const TIME_TOKEN_EXPIRATION = 24 * time.Hour

type Claims struct {
	UserID uuid.UUID `json:"user_id"`

	jwt.RegisteredClaims
}

func GenerateAccessToken(userID uuid.UUID, secret string, expires_in time.Duration) (string, error) {
	claims := Claims{
		UserID: userID,

		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(expires_in),
			),

			IssuedAt: jwt.NewNumericDate(
				time.Now(),
			),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString([]byte(secret))
}

func Validate(tokenString string, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,

		&Claims{},

		func(token *jwt.Token) (any, error) {
			return []byte(secret), nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)

	if !ok {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
