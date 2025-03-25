package config

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoClient este clientul global pentru MongoDB
var MongoClient *mongo.Client

// InitMongo conectează aplicația la MongoDB
func InitMongo() {
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Eroare la conectarea cu MongoDB: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Eroare la ping-ul MongoDB: %v", err)
	}

	MongoClient = client
	log.Println("Conexiune cu MongoDB realizată cu succes!")
}