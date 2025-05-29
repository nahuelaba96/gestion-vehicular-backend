package routes

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		fmt.Println("[AuthMiddleware] Authorization header:", authHeader)

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			fmt.Println("[AuthMiddleware] Token no provisto o mal formado")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token no provisto"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		jwtSecreto := os.Getenv("JWT_SECRET")

		if jwtSecreto == "" {
			fmt.Println("[AuthMiddleware] JWT_SECRET no está seteado en el entorno")
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("método de firma inesperado")
			}
			return []byte(jwtSecreto), nil
		})

		if err != nil {
			fmt.Printf("[AuthMiddleware] Error al verificar token: %v\n", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
			return
		}

		if !token.Valid {
			fmt.Println("[AuthMiddleware] Token inválido (pero sin error)")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
			return
		}

		fmt.Println("[AuthMiddleware] Token válido")
		c.Next()
	}
}


