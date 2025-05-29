package authentication

import (
	"context"
	"gestion-vehicular-backend/database"
	"gestion-vehicular-backend/models"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

	user, err := FindUserByEmail(email, name, c)
	if err != nil {
		log.Println("Error al buscar o crear usuario:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar o crear usuario"})
		return
	}

	// Crear JWT propio
	tokenStr, err := CrearJWT(user.ID.Hex(), email, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo generar token"})
		return
	}

	// Setear cookie compatible con cross-origin
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "jwt",
		Value:    tokenStr,
		Path:     "/",
		Domain:   "",
		MaxAge:   3600 * 24,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})

	// También podés enviar el token por JSON para testing/local
	c.JSON(http.StatusOK, gin.H{"token": tokenStr})
}

func FindUserByEmail(email, name string, c *gin.Context) (models.User, error) {
	var user models.User
    err := database.UserCollection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
    if err == mongo.ErrNoDocuments {
        // No existe: lo creo
        user = models.User{
            ID:     primitive.NewObjectID(),
            Email:  email,
            Name:   name,
        }
        _, err = database.UserCollection.InsertOne(context.Background(), user)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creando el usuario"})
            return user, err 
        }
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error buscando el usuario"})
		return user, err 
	}
	// Si el usuario existe o fue creado correctamente, retornar su ID
	return user, err 
}


func Logout(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "jwt",
		Value:    "",
		Path:     "/",
		MaxAge:   -1, // Expira inmediatamente
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})
	c.JSON(http.StatusOK, gin.H{"mensaje": "Sesión cerrada"})
}
