package authentication

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

func CrearJWT(id, email, name string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET no configurado")
	}

	claims := jwt.MapClaims{
		"user_id": id,
		"email":   email,
		"name":    name,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

