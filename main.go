package main

import (
	"github.com/gin-gonic/gin"
	"gestion-vehicular-backend/routes"
	"github.com/gin-contrib/cors"
	"time"
)

func main() {
   r := gin.Default()

	// Middleware CORS
 	r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"}, // Asegúrate que Authorization esté si lo usas
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
    }))

	// Ahora las rutas
	routes.SetupGestionVehicularRoutes(r)

	r.Run(":8080")

}
