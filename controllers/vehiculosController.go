package controllers

import (
	"context"
	"gestion-vehicular-backend/database"
	"gestion-vehicular-backend/models"
	"gestion-vehicular-backend/utils"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetVehiculos(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userID, err := utils.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userID, err := utils.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var v models.Vehiculo
	if err := c.ShouldBindJSON(&v); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Asociar usuario, ID y fecha creación
	v.ID = primitive.NewObjectID()
	v.UserID = userID
	v.FechaCreacion = time.Now()

	res, err := database.VehiculosCollection.InsertOne(ctx, v)
	if err != nil {
		log.Println("Error al insertar vehículo:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo guardar el vehículo"})
		return
	}

	v.ID = res.InsertedID.(primitive.ObjectID)

	c.JSON(http.StatusCreated, v)
}

func EliminarVehiculo(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Obtener user_id seguro con chequeo
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

	// Obtener y validar el ID del vehículo
	idStr := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Eliminar solo si el vehículo pertenece al usuario autenticado
	filter := bson.M{
		"_id":     objID,
		"user_id": userID,
	}

	result, err := database.VehiculosCollection.DeleteOne(ctx, filter)
	if err != nil {
		log.Println("Error al eliminar vehículo:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo eliminar el vehículo"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vehículo no encontrado o no pertenece al usuario"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensaje": "Vehículo eliminado con éxito"})
}

func ActualizarVehiculo(c *gin.Context) {
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

	idStr := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	delete(body, "user_id")
	delete(body, "_id")

	update := bson.M{"$set": body}
	filter := bson.M{"_id": objID, "user_id": userID}

	var vehiculoActualizado models.Vehiculo
	err = database.VehiculosCollection.FindOneAndUpdate(
		ctx,
		filter,
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&vehiculoActualizado)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Vehículo no encontrado o no pertenece al usuario"})
		} else {
			log.Println("Error al actualizar vehículo:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo actualizar el vehículo"})
		}
		return
	}

	c.JSON(http.StatusOK, vehiculoActualizado)
}
