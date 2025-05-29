package authentication

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func VerifyToken(c *gin.Context) {
	// Reusás tu lógica del middleware
	tokenStr, err := c.Cookie("jwt")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado"})
		return
	}

	jwtSecreto := os.Getenv("JWT_SECRET")

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("firma inválida")
		}
		return []byte(jwtSecreto), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Token válido"})
}
