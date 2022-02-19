package util

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
)

const (
	VarPrefix = "CLIENTDB_"
)

var (
	// config, _               = LoadConfig()
	MongoDBHost               = GetEnv(VarPrefix+"MONGODB_HOST", "localhost")
	MongoDBPort               = GetEnv(VarPrefix+"MONGODB_PORT", "27017")
	MongoDBName               = GetEnv(VarPrefix+"MONGODB_NAME", "test")
	LogLevel									= GetEnv(VarPrefix+"LOG_LEVEL", "info")
	ctx                       = context.TODO()
	DB          *mongo.Client = DbConnect()
)

// func GetEnv is a function to get environment variables to be able to set default values if not set
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
