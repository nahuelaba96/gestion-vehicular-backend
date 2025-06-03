package controllers

import (
	"context"
	"encoding/json"
	"gestion-vehicular-backend/database"
	"gestion-vehicular-backend/models"
	"gestion-vehicular-backend/utils"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	// Validar si se envió vehiculo_id
	if !gasto.VehiculoID.IsZero() {
		// Si vino, validar que exista y pertenezca al usuario
		filtroVehiculo := bson.M{"_id": gasto.VehiculoID, "user_id": userID}
		count, err := database.VehiculosCollection.CountDocuments(ctx, filtroVehiculo)
		if err != nil || count == 0 {
			c.JSON(http.StatusForbidden, gin.H{"error": "Vehículo no válido o no pertenece al usuario"})
			return
		}
	} else {
		// Si no vino, dejar el campo vacío o nulo, para que sea opcional
		gasto.VehiculoID = primitive.NilObjectID
	}

	if strings.ToLower(gasto.Tipo) == "compra" {
	var total float64
	for _, item := range gasto.Items {
		total += item.Cantidad * item.PrecioUnitario
	}
	gasto.Total = total

	} else if strings.ToLower(gasto.Tipo) == "combustible" && gasto.Combustible != nil {
		if gasto.Combustible.PrecioLitro > 0 && gasto.Combustible.LitrosEstimados > 0 {
			gasto.Total = gasto.Combustible.PrecioLitro * gasto.Combustible.LitrosEstimados
		}

	} else if strings.ToLower(gasto.Tipo) == "mecanico" && gasto.Mecanico != nil {
		var total float64
		for _, item := range gasto.Items {
			total += item.Cantidad * item.PrecioUnitario
		}
		total += gasto.Mecanico.CostoManoObra
		gasto.Total = total
	}

	gasto.FechaInsert = time.Now()
	gasto.UserID = userID
	gasto.ID = primitive.NewObjectID()

	res, err := database.GastosCollection.InsertOne(ctx, gasto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo guardar el gasto"})
		return
	}

	gasto.ID = res.InsertedID.(primitive.ObjectID)

	c.JSON(http.StatusCreated, gasto)
}


func ListarGastos(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userID, err := utils.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	tipo := c.Query("tipo")
	fechaDesde := c.Query("fecha_desde")
	fechaHasta := c.Query("fecha_hasta")

	filtro := bson.M{"user_id": userID}
	if tipo != "" {
		filtro["tipo"] = strings.ToLower(tipo)
	}
	if fechaDesde != "" || fechaHasta != "" {
		rango := bson.M{}
		if fechaDesde != "" {
			rango["$gte"] = fechaDesde
		}
		if fechaHasta != "" {
			rango["$lte"] = fechaHasta
		}
		filtro["fecha"] = rango
	}

	cursor, err := database.GastosCollection.Find(ctx, filtro)
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
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)

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
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}

	camposActualizables := []string{"vehiculo_id", "fecha", "proveedor", "nota", "items", "total"}
	set := bson.M{}
	for _, campo := range camposActualizables {
		if val, ok := body[campo]; ok {
			set[campo] = val
		}
	}

	// Validación específica para tipo "compra", "combustible", "mecanico"
	if tipo, ok := body["tipo"].(string); ok {
		switch strings.ToLower(tipo) {
		case "compra":
			if items, ok := body["items"].([]interface{}); ok {
				var total float64
				for _, item := range items {
					if itemMap, ok := item.(map[string]interface{}); ok {
						cant, _ := toFloat(itemMap["cantidad"])
						precio, _ := toFloat(itemMap["precio_unitario"])
						total += cant * precio
					}
				}
				set["total"] = total
			}
		case "combustible":
			if comb, ok := body["combustible"].(map[string]interface{}); ok {
				precio, _ := toFloat(comb["precio_litro"])
				litros, _ := toFloat(comb["litros_estimados"])
				set["total"] = precio * litros
			}
		case "mecanico":
			var total float64
			if items, ok := body["items"].([]interface{}); ok {
				for _, item := range items {
					if itemMap, ok := item.(map[string]interface{}); ok {
						cant, _ := toFloat(itemMap["cantidad"])
						precio, _ := toFloat(itemMap["precio_unitario"])
						total += cant * precio
					}
				}
			}
			if mec, ok := body["mecanico"].(map[string]interface{}); ok {
				manoObra, _ := toFloat(mec["costo_mano_obra"])
				total += manoObra
			}
			set["total"] = total
		}
	}

	filter := bson.M{
		"_id":      objID,
		"user_id":  userID,
	}

	update := bson.M{"$set": body}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var gastoActualizado models.Gasto
	err = database.GastosCollection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&gastoActualizado)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Gasto no encontrado o no pertenece al usuario"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar gasto"})
		}
		return
	}

	c.JSON(http.StatusOK, gastoActualizado)
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

func toFloat(v interface{}) (float64, bool) {
	switch val := v.(type) {
	case float64:
		return val, true
	case float32:
		return float64(val), true
	case int:
		return float64(val), true
	case int32:
		return float64(val), true
	case int64:
		return float64(val), true
	case json.Number:
		f, err := val.Float64()
		return f, err == nil
	default:
		return 0, false
	}
}

func ActualizarKilometrajeVehiculo(ctx context.Context, vehiculoID, userID primitive.ObjectID) error {
	filter := bson.M{
		"vehiculo_id":        vehiculoID,
		"user_id":            userID,
		"tipo":               "mecanico",
		"mecanico.kilometraje": bson.M{"$exists": true},
	}
	opts := options.FindOne().SetSort(bson.D{{"fecha_insert", -1}})

	var gasto models.Gasto
	err := database.GastosCollection.FindOne(ctx, filter, opts).Decode(&gasto)

	var nuevoKm int
	if err == mongo.ErrNoDocuments {
		// No hay historial mecánico, usar kilometraje de registro
		var vehiculo models.Vehiculo
		err = database.VehiculosCollection.FindOne(ctx, bson.M{"_id": vehiculoID}).Decode(&vehiculo)
		if err != nil {
			return err
		}
		nuevoKm = int(vehiculo.KilometrosRegistro)
	} else if err == nil && gasto.Mecanico != nil {
		nuevoKm = int(gasto.Mecanico.Kilometraje)
	} else {
		return err
	}

	_, err = database.VehiculosCollection.UpdateOne(ctx, bson.M{"_id": vehiculoID}, bson.M{
		"$set": bson.M{"kilometraje_actual": nuevoKm},
	})
	return err
}
