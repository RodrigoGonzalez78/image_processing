package middleware

import (
	"context"
	"net/http"

	"github.com/RodrigoGonzalez78/internal/domain/service"
)

type JWTMiddleware struct {
	tokenService *service.TokenService
}

func NewJWTMiddleware(tokenService *service.TokenService) *JWTMiddleware {
	return &JWTMiddleware{tokenService: tokenService}
}

type contextKey string

const userDataKey contextKey = "userData"

func (m *JWTMiddleware) Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Token requerido", http.StatusUnauthorized)
			return
		}

		claims, err := m.tokenService.Validate(authHeader)
		if err != nil {
			http.Error(w, "Token inválido: "+err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userDataKey, claims)
		next(w, r.WithContext(ctx))
	}
}

func GetUserFromContext(ctx context.Context) (*service.TokenClaims, bool) {
	claims, ok := ctx.Value(userDataKey).(*service.TokenClaims)
	return claims, ok
}
