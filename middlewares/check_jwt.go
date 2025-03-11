package middlewares

import (
	"context"
	"net/http"

	"github.com/RodrigoGonzalez78/utils"
)

func CheckJwt(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		claim, _, err := utils.ProcessToken(r.Header.Get("Authorization"))

		if err != nil {
			http.Error(w, "Erro en el token!"+err.Error(), http.StatusBadRequest)
			return
		}
		ctx := context.WithValue(r.Context(), "userData", claim)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
