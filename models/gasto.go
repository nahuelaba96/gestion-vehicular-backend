package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ItemGasto struct {
	Descripcion    string  `bson:"descripcion" json:"descripcion"`
	Cantidad       float64 `bson:"cantidad" json:"cantidad"`
	PrecioUnitario float64 `bson:"precio_unitario" json:"precio_unitario"`
}

type Gasto struct {
	ID           primitive.ObjectID `bson:"_id" json:"id"`
	UserID       primitive.ObjectID `bson:"user_id" json:"-"`
	VehiculoID   primitive.ObjectID `bson:"vehiculo_id" json:"vehiculo_id"`
	FechaInsert  time.Time          `bson:"fecha_insert" json:"fecha_insert"`
	Fecha        string             `bson:"fecha" json:"fecha"`
	Tipo         string             `bson:"tipo" json:"tipo"` // "compra", "combustible", "mecanico"
	Proveedor    string             `bson:"proveedor,omitempty" json:"proveedor,omitempty"`
	Nota         string             `bson:"nota,omitempty" json:"nota,omitempty"`
	Items        []ItemGasto        `bson:"items,omitempty" json:"items,omitempty"`
	Total        float64            `bson:"total" json:"total"`

	// Combustible
	Combustible *GastoCombustible `bson:"combustible,omitempty" json:"combustible,omitempty"`

	// Mantenimiento/Mecánico
	Mecanico *TrabajoMecanico `bson:"mecanico,omitempty" json:"mecanico,omitempty"`
}

type GastoCombustible struct {
	TipoCombustible string  `bson:"tipo_combustible" json:"tipo_combustible"` // super, premium, diesel, etc.
	PrecioLitro     float64 `bson:"precio_litro,omitempty" json:"precio_litro,omitempty"`
	LitrosEstimados float64 `bson:"litros_estimados,omitempty" json:"litros_estimados,omitempty"`
	Ubicacion       string  `bson:"ubicacion,omitempty" json:"ubicacion,omitempty"`
	Estacion        string  `bson:"estacion,omitempty" json:"estacion,omitempty"` // YPF, Shell, etc.
}

type TrabajoMecanico struct {
	Descripcion   string    `bson:"descripcion" json:"descripcion"` // Ej: Cambio de amortiguador
	CostoManoObra float64   `bson:"costo_mano_obra,omitempty" json:"costo_mano_obra,omitempty"`
	HechoPor      string    `bson:"hecho_por,omitempty" json:"hecho_por,omitempty"` // Nombre del mecánico o taller
	FechaRealizado *string  `bson:"fecha_realizado,omitempty" json:"fecha_realizado,omitempty"` // null si está pendiente
	Estado        string    `bson:"estado" json:"estado"` // pendiente, realizado, cancelado, etc.
	Kilometraje   int64     `bson:"kilometraje,omitempty" json:"kilometraje,omitempty"` // Kilometraje al momento del trabajo
}
