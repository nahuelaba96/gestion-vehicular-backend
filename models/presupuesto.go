package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Presupuesto struct {
	ID         primitive.ObjectID `bson:"_id" json:"id"`
	UserID     primitive.ObjectID `bson:"user_id" json:"-"`
	VehiculoID primitive.ObjectID `bson:"vehiculo_id" json:"vehiculo_id"`
	MontoMax   float64            `bson:"monto_max" json:"monto_max"`
	Desde      time.Time          `bson:"desde" json:"desde"`
	Hasta      time.Time          `bson:"hasta" json:"hasta"`
	Notas      string             `bson:"notas,omitempty" json:"notas,omitempty"`
}
