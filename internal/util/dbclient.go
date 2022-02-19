package util

import (
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DBConnect is a function to connect to the database
func DbConnect() *mongo.Client {

	connectionString := "mongodb://" + MongoDBHost + ":" + MongoDBPort + "/" + MongoDBName
	log.Debug("Using connection string: ", connectionString)
	// Set client options
	mongoOptions := options.Client().ApplyURI(connectionString)

	// Connect to MongoDB
	mongoClient, err := mongo.Connect(ctx, mongoOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = mongoClient.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Connected to MongoDB")
	return mongoClient

}

func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database(MongoDBName).Collection(collectionName)
	return collection
}
