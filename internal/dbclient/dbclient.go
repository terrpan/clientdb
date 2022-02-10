package dbclient

import (
	"context"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	// DB variable to hold the database connection
	addr string        = "localhost"
	port string        = "27017"
	db   string        = "test"
	ctx                = context.TODO()
	DB   *mongo.Client = DbConnect()
)

// DBConnect is a function to connect to the database
func DbConnect() *mongo.Client {

	connectionString := "mongodb://" + addr + ":" + port + "/" + db

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
	collection := client.Database(db).Collection(collectionName)
	return collection
}
