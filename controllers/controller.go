package controllers

import (
	"gestion-vehicular-backend/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var datos = []models.Datos{}
var nextID int64 = 1

func GetDatos(c *gin.Context) {
	datos := models.Datos{
		ID:                  1,
		Auto:                "BMW 323I 2008",
		Patente:             "GWO040",
		TipoBitacora:        "Compra",
		ComponenteRecambio:  "Cubiertas traseras",
		ComponenteInstalado: "SI",
		Marca:               "Sailum",
		Fecha:               "2024-08-27",
		Vendedor:            "Saracho Neumaticos",
		Kilometro:           171191,
		Costo:               320000.0,
		Nota:                "Compra 2 cubiertas 235/40/18 marca Sailum",
		FechaProximo:        "2034-08-27",
		KilometrosProximo:   180000,
	}

	c.JSON(http.StatusOK, datos)
}

func CreateGestionVehicular(c *gin.Context) {
	var newBitacora models.Datos
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
