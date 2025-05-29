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

	val, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No se pudo obtener el usuario autenticado"})
		return
	}

	userID, ok := val.(primitive.ObjectID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ID de usuario inválido"})
		return
	}

	cursor, err := database.VehiculosCollection.Find(ctx, bson.M{"user_id": userID})
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

	// Obtener el user_id del contexto (puesto por el AuthMiddleware)
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	v.ID = primitive.NewObjectID()          // ID del vehículo
	v.UserID = userID.(primitive.ObjectID)  // Asociar al usuario autenticado
	v.FechaCreacion = time.Now()            // Fecha actual

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
	userID := c.MustGet("user_id").(primitive.ObjectID)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	result, err := database.VehiculosCollection.DeleteOne(ctx, bson.M{"_id": objID, "user_id": userID})
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
	userID := c.MustGet("user_id").(primitive.ObjectID)

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

	result, err := database.VehiculosCollection.UpdateOne(ctx, bson.M{"_id": objID, "user_id": userID}, update)
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
