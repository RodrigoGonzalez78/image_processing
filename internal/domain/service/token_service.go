package service

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenClaims struct {
	UserName string `json:"user_name"`
	jwt.RegisteredClaims
}

type TokenService struct {
	secretKey []byte
	expiry    time.Duration
}

func NewTokenService(secret string, expiry time.Duration) *TokenService {
	return &TokenService{
		secretKey: []byte(secret),
		expiry:    expiry,
	}
}

func (s *TokenService) Generate(userName string) (string, error) {
	claims := TokenClaims{
		UserName: userName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

func (s *TokenService) Validate(tokenStr string) (*TokenClaims, error) {
	if !strings.HasPrefix(tokenStr, "Bearer ") {
		return nil, errors.New("formato de token inválido")
	}

	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

	var claims TokenClaims
	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(t *jwt.Token) (interface{}, error) {
		return s.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("token inválido")
	}

	return &claims, nil
}
