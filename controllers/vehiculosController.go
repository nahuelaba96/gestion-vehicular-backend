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
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetVehiculos(c *gin.Context) {
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

func ObtenerGastosPorVehiculo(c *gin.Context) {
	idStr := c.Param("id")
	vehiculoID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	filter := bson.M{"vehiculo_id": vehiculoID}
	cursor, err := database.GastosCollection.Find(context.TODO(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar gastos"})
		return
	}
	defer cursor.Close(context.TODO())

	var gastos []models.Gasto
	if err := cursor.All(context.TODO(), &gastos); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al procesar gastos"})
		return
	}

	c.JSON(http.StatusOK, gastos)
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

func EliminarVehiculo(c *gin.Context) {
	id := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	result, err := database.VehiculosCollection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		log.Println("Error al eliminar vehículo:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo eliminar el vehículo"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vehículo no encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensaje": "Vehículo eliminado con éxito"})
}

func ActualizarVehiculo(c *gin.Context) {
	id := c.Param("id")

	// Parsear ID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Leer solo los campos enviados
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{"$set": body}

	result, err := database.VehiculosCollection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		log.Println("Error al actualizar vehículo:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo actualizar el vehículo"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vehículo no encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensaje": "Vehículo actualizado con éxito"})
}
