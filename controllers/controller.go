package controllers

import (
	"gestion-vehicular-backend/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var datos = []models.Vehiculo{}
var nextID int64 = 1

func GetDatos(c *gin.Context) {
	vehiculo := models.Vehiculo{
		ID:                  1,
		Tipo:                "Sedán",
		Patente:             "GWO040",
		Marca:               "BMW",
		Modelo:              "323i",
		Anio:                "2008",
		TipoCombustible:     "Gasolina",
		Kilometros:          173020,
		Nota:                "Vehículo en buen estado, mantenimiento reciente.",
	}
	datos = []models.Vehiculo{} // Reiniciar el slice de Vehiculo
	datos := []models.Vehiculo{vehiculo}


	c.JSON(http.StatusOK, datos)
}

func CreateGestionVehicular(c *gin.Context) {
	var newBitacora models.Vehiculo
	if err := c.ShouldBindJSON(&newBitacora); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newBitacora.ID = nextID
	nextID++
	datos = append(datos, newBitacora)
	c.JSON(http.StatusCreated, newBitacora)
}

func GetGestionVehicularByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	for _, b := range datos {
		if b.ID == id {
			c.JSON(http.StatusOK, b)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "Registro no encontrado"})
}
