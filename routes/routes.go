package routes

import (
	"gestion-vehicular-backend/controllers"

	"github.com/gin-gonic/gin"
)

func SetupGestionVehicularRoutes(r *gin.Engine) {
	bitacora := r.Group("/gestion-vehicular")
	{
		bitacora.GET("/", controllers.GetDatos)
		bitacora.POST("/", controllers.CreateGestionVehicular)
		bitacora.GET("/:id", controllers.GetGestionVehicularByID)
	}
}
