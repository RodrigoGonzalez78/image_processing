package utils

import (
	"time"

	"github.com/RodrigoGonzalez78/models"
	"github.com/golang-jwt/jwt"
)

func GeneringJwt(t models.User) (string, error) {
	myClave := []byte("Twittor_clone")

	payload := jwt.MapClaims{
		"user_name": t.UserName,
		"user_id":   t.UserId,
		"exp":       time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	tokenStrin, err := token.SignedString(myClave)

	if err != nil {
		return tokenStrin, err
	}

	return tokenStrin, nil
}
