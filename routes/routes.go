package routes

import (
	"gestion-vehicular-backend/controllers"

	"github.com/gin-gonic/gin"
)

func SetupGestionVehicularRoutes(r *gin.Engine) {
	vehiculos := r.Group("/vehicles")
	{
		vehiculos.GET("/", controllers.GetVehiculos)
		vehiculos.POST("/", controllers.CreateVehiculo)
		vehiculos.DELETE("/:id", controllers.EliminarVehiculo)
		vehiculos.PUT("/:id", controllers.ActualizarVehiculo)
	}

	gastos := r.Group("/expenses")
	{
		gastos.GET("/", controllers.ListarGastos)
		gastos.GET("/:id", controllers.ObtenerGastosPorVehiculo)
		gastos.POST("/", controllers.CrearGasto)
		gastos.DELETE("/:id", controllers.EliminarGasto)
		gastos.PUT("/:id", controllers.ActualizarGasto)
	}
}
