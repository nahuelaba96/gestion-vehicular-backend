package controllers

import (
	"context"
	"gestion-vehicular-backend/database"
	"gestion-vehicular-backend/models"
	"gestion-vehicular-backend/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CrearGasto(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userID, err := utils.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var gasto models.Gasto
	if err := c.ShouldBindJSON(&gasto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if gasto.VehiculoID.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID del vehículo es requerido"})
		return
	}

	filtroVehiculo := bson.M{"_id": gasto.VehiculoID, "user_id": userID}
	count, err := database.VehiculosCollection.CountDocuments(ctx, filtroVehiculo)
	if err != nil || count == 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Vehículo no válido o no pertenece al usuario"})
		return
	}

	var total float64
	for _, item := range gasto.Items {
		total += float64(item.Cantidad) * item.PrecioUnitario
	}
	gasto.Total = total
	gasto.FechaInsert = time.Now()
	gasto.UserID = userID

	res, err := database.GastosCollection.InsertOne(ctx, gasto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo guardar el gasto"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": res.InsertedID})
}

func ListarGastos(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userID, err := utils.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	cursor, err := database.GastosCollection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo obtener la lista de gastos"})
		return
	}
	defer cursor.Close(ctx)

	var gastos []models.Gasto
	if err := cursor.All(ctx, &gastos); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al parsear datos"})
		return
	}

	c.JSON(http.StatusOK, gastos)
}

func ObtenerGastosPorVehiculo(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userID, err := utils.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	idStr := c.Param("id")
	vehiculoID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de vehículo inválido"})
		return
	}

	filterVehiculo := bson.M{"_id": vehiculoID, "user_id": userID}
	count, err := database.VehiculosCollection.CountDocuments(ctx, filterVehiculo)
	if err != nil || count == 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Vehículo no válido o no pertenece al usuario"})
		return
	}

	filterGastos := bson.M{"vehiculo_id": vehiculoID, "user_id": userID}
	cursor, err := database.GastosCollection.Find(ctx, filterGastos)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar gastos"})
		return
	}
	defer cursor.Close(ctx)

	var gastos []models.Gasto
	if err := cursor.All(ctx, &gastos); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al procesar gastos"})
		return
	}

	c.JSON(http.StatusOK, gastos)
}

func ActualizarGasto(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userID, err := utils.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
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

	if itemsRaw, ok := body["items"]; ok {
		if items, ok := itemsRaw.([]interface{}); ok {
			var total float64
			for _, item := range items {
				if itemMap, ok := item.(map[string]interface{}); ok {
					cant, _ := itemMap["cantidad"].(float64)
					precio, _ := itemMap["precio_unitario"].(float64)
					total += cant * precio
				}
			}
			body["total"] = total
		}
	}

	delete(body, "user_id")

	filter := bson.M{"_id": objID, "user_id": userID}
	update := bson.M{"$set": body}

	res, err := database.GastosCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar gasto"})
		return
	}

	if res.MatchedCount == 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Gasto no encontrado o no pertenece al usuario"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"matched": res.MatchedCount, "modified": res.ModifiedCount})
}

func EliminarGasto(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userID, err := utils.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	idStr := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	filter := bson.M{"_id": objID, "user_id": userID}
	res, err := database.GastosCollection.DeleteOne(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar gasto"})
		return
	}

	if res.DeletedCount == 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Gasto no encontrado o no pertenece al usuario"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensaje": "Gasto eliminado"})
}