package main

import (
	"gestion-vehicular-backend/database"
	"gestion-vehicular-backend/routes"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	r := gin.Default()

	_ = godotenv.Load()

	database.ConnectMongo()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://127.0.0.1:5173","https://c196-45-234-34-114.ngrok-free.app",},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Rutas públicas
	routes.SetupPublicRoutes(r)

	// Grupo de rutas protegidas con middleware de autenticación
	vehiculosGroup := routes.SetupProtectedRoutes(r, "/vehicles")
	routes.SetupVehiculosRoutes(vehiculosGroup)

	gastosGroup := routes.SetupProtectedRoutes(r, "/expenses")
	routes.SetupGastosRoutes(gastosGroup)


	r.Run(":8080")
}

