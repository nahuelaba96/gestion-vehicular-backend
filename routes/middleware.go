package routes

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		// Intentar obtener token desde header Authorization
		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
			// Validar tokenString para que no sea vacío o "null"
			if tokenString == "" || tokenString == "null" || tokenString == "undefined" {
				// Token inválido en header, intentar cookie
				tokenString = ""
			}
		}

		// Si token no viene o es inválido en header, probar cookie
		if tokenString == "" {
			cookieToken, err := c.Cookie("jwt")
			if err != nil || cookieToken == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token no provisto"})
				c.Abort()
				return
			}
			tokenString = cookieToken
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

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Claims inválidos"})
			return
		}

		//  Guardás el user_id en el contexto
		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user_id no presente"})
			return
		}

		userID, err := primitive.ObjectIDFromHex(userIDStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user_id inválido"})
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}

