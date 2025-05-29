package routes

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Intentar obtener el token desde la cookie
		tokenString, err := c.Cookie("jwt")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token no provisto"})
			c.Abort()
			return
		}

		// Obtener la clave secreta desde las variables de entorno
		jwtSecreto := os.Getenv("JWT_SECRET") 

		// Verificar el token con la misma clave secreta que usaste al generarlo
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validar algoritmo
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("método de firma inesperado")
			}
			return []byte(jwtSecreto), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
			return
		}

		c.Next()
	}
}

