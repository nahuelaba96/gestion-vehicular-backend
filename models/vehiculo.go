package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Vehiculo struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Tipo            string             `json:"tipo" bson:"tipo"`
	Patente         string             `json:"patente" bson:"patente"`
	Marca           string             `json:"marca" bson:"marca"`
	Modelo          string             `json:"modelo" bson:"modelo"`
	Anio            string             `json:"anio" bson:"anio"`
	TipoCombustible string             `json:"tipo_combustible" bson:"tipo_combustible"`
	Kilometros      int64              `json:"kilometros" bson:"kilometros"`
	Nota            string             `json:"nota" bson:"nota"`
	FechaCreacion   time.Time          `json:"fecha_creacion" bson:"fecha_creacion"`
}
