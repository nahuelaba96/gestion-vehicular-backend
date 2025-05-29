package authentication

import (
	"context"
	"log"
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
	if oauthClient == "" {
		log.Println("OAUTH_CLIENT no seteado")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Configuración incorrecta del servidor"})
		return
	}

	// Validar token con Google
	payload, err := idtoken.Validate(context.Background(), req.Token, oauthClient)
	if err != nil {
		log.Println("Error al validar token de Google:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token no válido"})
		return
	}

	email, ok := payload.Claims["email"].(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Email no encontrado en el token"})
		return
	}
	name, _ := payload.Claims["name"].(string)

	// Crear JWT propio
	tokenStr, err := CrearJWT(email, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo generar token"})
		return
	}

	// Setear cookie compatible con cross-origin
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "jwt",
		Value:    tokenStr,
		Path:     "/",
		Domain:   "", // sin dominio para que funcione en Railway (o usar el dominio final si tenés uno)
		MaxAge:   3600 * 24,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode, // esto es clave
	})

	// También podés enviar el token por JSON para testing/local
	c.JSON(http.StatusOK, gin.H{"token": tokenStr})
}


