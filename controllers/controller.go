package controllers

import (
	"context"
	"gestion-vehicular-backend/database"
	"gestion-vehicular-backend/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

var datos = []models.Vehiculo{}

func GetDatos(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := database.VehiculosCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Println("Error al buscar vehículos:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo obtener la lista de vehículos"})
		return
	}
	defer cursor.Close(ctx)

	var vehiculos []models.Vehiculo
	if err := cursor.All(ctx, &vehiculos); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al parsear datos"})
		return
	}

	c.JSON(http.StatusOK, vehiculos)
}

func CreateVehiculo(c *gin.Context) {
	var v models.Vehiculo
	if err := c.BindJSON(&v); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	v.FechaCreacion = time.Now() // Asignar la fecha de creación actual

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := database.VehiculosCollection.InsertOne(ctx, v)
	if err != nil {
		log.Println("Error al insertar:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo guardar el vehículo"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"mensaje": "Vehículo creado con éxito"})
}

