package routes

import (
	"gestion-vehicular-backend/authentication"
	"gestion-vehicular-backend/controllers"

	"github.com/gin-gonic/gin"
)

func SetupVehiculosRoutes(vehiculos *gin.RouterGroup) {
	{
		vehiculos.GET("/", controllers.GetVehiculos)
		vehiculos.POST("/", controllers.CreateVehiculo)
		vehiculos.DELETE("/:id", controllers.EliminarVehiculo)
		vehiculos.PUT("/:id", controllers.ActualizarVehiculo)
	}
}

func SetupGastosRoutes(gastos *gin.RouterGroup) {
	{
		gastos.GET("/", controllers.ListarGastos)
		gastos.GET("/:id", controllers.ObtenerGastosPorVehiculo)
		gastos.POST("/", controllers.CrearGasto)
		gastos.DELETE("/:id", controllers.EliminarGasto)
		gastos.PUT("/:id", controllers.ActualizarGasto)
	}
}

func SetupPublicRoutes(r *gin.Engine) {
	auth := r.Group("/auth")
	{
		auth.POST("/google", authentication.GoogleLogin)
		auth.GET("/verify", authentication.VerifyToken)
	}

}

func SetupProtectedRoutes(r *gin.Engine, prefix string) *gin.RouterGroup {
	protected := r.Group(prefix)
	protected.Use(AuthMiddleware())
	return protected
}
