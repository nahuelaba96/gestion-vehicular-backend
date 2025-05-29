// database/mongo.go
package database

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var VehiculosCollection *mongo.Collection
var GastosCollection *mongo.Collection
var UserCollection *mongo.Collection

func ConnectMongo() {
	mongoUser := os.Getenv("MONGO_USER")
	mongoPassword := os.Getenv("MONGO_PASSWORD")

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
  	opts := options.Client().ApplyURI("mongodb+srv://"+mongoUser+":"+mongoPassword+"@cluster-gestion-vehicul.czzbkqo.mongodb.net/?retryWrites=true&w=majority&appName=cluster-gestion-vehicular").SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
  	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Fatal("Error al conectar con MongoDB:", err)
	}

	// Ping para asegurarse de que se conecta
	if err := client.Database("miapp").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		log.Fatal("Ping falló:", err)
	}

	MongoClient = client
	VehiculosCollection = client.Database("miapp").Collection("vehiculos")
	GastosCollection = client.Database("miapp").Collection("gastos")
	UserCollection = client.Database("miapp").Collection("usuarios")

	log.Println("Conectado a MongoDB")
}
