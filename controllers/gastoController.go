package controllers

import (
	"context"
	"gestion-vehicular-backend/database"
	"gestion-vehicular-backend/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CrearGasto(c *gin.Context) {
	var gasto models.Gasto
	if err := c.ShouldBindJSON(&gasto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Calcular total
	var total float64
	for _, item := range gasto.Items {
		total += float64(item.Cantidad) * item.PrecioUnitario
	}
	gasto.Total = total
	gasto.FechaInsert = time.Now()

	res, err := database.GastosCollection.InsertOne(context.TODO(), gasto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al insertar gasto"})
		return
	}
	c.JSON(http.StatusOK, res.InsertedID)
}

func ListarGastos(c *gin.Context) {
	cursor, err := database.GastosCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar gastos"})
		return
	}
	defer cursor.Close(context.TODO())

	var gastos []models.Gasto
	if err = cursor.All(context.TODO(), &gastos); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al parsear gastos"})
		return
	}

	c.JSON(http.StatusOK, gastos)
}

func ActualizarGasto(c *gin.Context) {
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

	// Si el body contiene items, podés recalcular el total
	if itemsRaw, ok := body["items"]; ok {
		items, ok := itemsRaw.([]interface{})
		if ok {
			var total float64
			for _, item := range items {
				itemMap, ok := item.(map[string]interface{})
				if !ok {
					continue
				}
				cant, _ := itemMap["cantidad"].(float64)
				precio, _ := itemMap["precio_unitario"].(float64)
				total += cant * precio
			}
			body["total"] = total
		}
	}

	update := bson.M{"$set": body}

	res, err := database.GastosCollection.UpdateOne(context.TODO(), bson.M{"_id": objID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar gasto"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"matched": res.MatchedCount, "modified": res.ModifiedCount})
}


func EliminarGasto(c *gin.Context) {
	idStr := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	_, err = database.GastosCollection.DeleteOne(context.TODO(), bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar gasto"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensaje": "Gasto eliminado"})
}

