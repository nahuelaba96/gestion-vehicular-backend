// database/mongo.go
package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var VehiculosCollection *mongo.Collection

// mongo password RMc2axKooRJ5hGrT
// mongo username nahuelaba96

func ConnectMongo() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// compose
	//clientOptions := options.Client().ApplyURI("mongodb://mongo:27017")

	// local
	//clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// mongo
	clientOptions := options.Client().ApplyURI("mongodb+srv://nahuelaba96:RMc2axKooRJ5hGrT@cluster-gestion-vehicul.czzbkqo.mongodb.net/?retryWrites=true&w=majority&appName=cluster-gestion-vehicular")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Error al conectar con MongoDB:", err)
	}

	MongoClient = client
	VehiculosCollection = client.Database("miapp").Collection("vehiculos")

	log.Println("Conectado a MongoDB")
}
