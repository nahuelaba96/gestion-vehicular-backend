package routes

import (
	"gestion-vehicular-backend/controllers"

	"github.com/gin-gonic/gin"
)

func SetupGestionVehicularRoutes(r *gin.Engine) {
	gestionVehicular := r.Group("/vehicles")
	{
		gestionVehicular.GET("/", controllers.GetVehiculos)
		gestionVehicular.POST("/create", controllers.CreateVehiculo)
		gestionVehicular.DELETE("/delete/:id", controllers.EliminarVehiculo)
		gestionVehicular.PUT("/update/:id", controllers.ActualizarVehiculo)
	}
}
