package database

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

// InitializeMongoDB initializes and returns a MongoDB client
func InitializeMongoDB() {
	// MongoDB Atlas connection string
	connectionString := "mongodb+srv://shivaditya2024:sLuRn1KXvm2iKEHk@cluster0.xnmo4.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"

	clientOptions := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB!")
	MongoClient = client
}

// GetMongoDBCollection is a helper function to get a collection from the database
func GetMongoDBCollection(databaseName, collectionName string) *mongo.Collection {
	return MongoClient.Database(databaseName).Collection(collectionName)
}

// DisconnectMongoDB is used to disconnect the MongoDB client when the application closes
func DisconnectMongoDB() {
	if MongoClient != nil {
		if err := MongoClient.Disconnect(context.TODO()); err != nil {
			log.Fatalf("Failed to disconnect from MongoDB: %v", err)
		}
		log.Println("Disconnected from MongoDB.")
	}
}
