package main

import (
	"github.com/gin-gonic/gin"
	"gestion-vehicular-backend/routes"
)

func main() {
	r := gin.Default()
	routes.SetupGestionVehicularRoutes(r)
	r.Run(":8080")
}
