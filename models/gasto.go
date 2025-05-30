package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ItemGasto struct {
	Descripcion    string  `bson:"descripcion" json:"descripcion"`
	Cantidad       int     `bson:"cantidad" json:"cantidad"`
	PrecioUnitario float64 `bson:"precio_unitario" json:"precio_unitario"`
}

type Gasto struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      primitive.ObjectID `bson:"user_id" json:"-"`
	VehiculoID  primitive.ObjectID `bson:"vehiculo_id" json:"vehiculo_id"`
	FechaInsert time.Time          `bson:"fecha_insert" json:"fecha_insert"`
	Fecha       string             `bson:"fecha" json:"fecha"`
	Proveedor   string             `bson:"proveedor" json:"proveedor"`
	Items       []ItemGasto        `bson:"items" json:"items"`
	Total       float64            `bson:"total" json:"total"`
}
