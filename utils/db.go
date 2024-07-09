package utils

import (
	"context"
	"os"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	clientInstance      *mongo.Client
	clientInstanceError error
	mongoOnce           sync.Once
)

const (
	CONNECTIONSTRING = "MONGO_URI"
)

// GetMongoClient initializes and returns the MongoDB client instance
func GetMongoClient() (*mongo.Client, error) {
	mongoOnce.Do(func() {
		clientOptions := options.Client().ApplyURI(os.Getenv(CONNECTIONSTRING))
		clientInstance, clientInstanceError = mongo.Connect(context.TODO(), clientOptions)
	})

	return clientInstance, clientInstanceError
}
