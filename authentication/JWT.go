package authentication

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

func CrearJWT(email, name string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"name":  name,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET") // poné una variable de entorno segura
	return token.SignedString([]byte(secret))
}
