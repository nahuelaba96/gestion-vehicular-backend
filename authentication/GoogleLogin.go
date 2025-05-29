package authentication

import (
	"context"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/idtoken"
)

func GoogleLogin(c *gin.Context) {
	var req struct {
		Token string `json:"token"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token inválido"})
		return
	}

	oauthClient := os.Getenv("OAUTH_CLIENT")

	// Validar token con Google
	payload, err := idtoken.Validate(context.Background(), req.Token, oauthClient)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token no válido"})
		return
	}

	email := payload.Claims["email"].(string)
	name := payload.Claims["name"].(string)

	// Podés guardar usuario en tu base si no existe

	// Crear JWT propio
	tokenStr, err := CrearJWT(email, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo generar token"})
		return
	}

	c.SetCookie("jwt", tokenStr, 3600*24, "/", "localhost", false, false)

	c.JSON(http.StatusOK, gin.H{"token": tokenStr})
}

