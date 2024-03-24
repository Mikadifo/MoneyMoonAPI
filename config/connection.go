package config

import (
	"context"
	"fmt"
	"mikadifo/money-moon/utily"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

func init() {
	if err := connectToMongoDB(); err != nil {
		fmt.Print("Could not connect to MongoDB")
	}
}

func connectToMongoDB() error {
	DB_URL := utily.GetEnvVar("DB_URL")
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(DB_URL).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	err = client.Ping(context.TODO(), nil)
	MongoClient = client

	return err
}
