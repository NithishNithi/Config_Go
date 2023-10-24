package database

import (
	"Config/constants"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func ConnectDatabase() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	mongoConnection := options.Client().ApplyURI(constants.ConnectionString)
	MongoClient, err := mongo.Connect(ctx, mongoConnection)
	if err != nil {
		log.Fatal()
		return nil, err
	}

	err1 := MongoClient.Ping(ctx, readpref.Primary())
	if err1 != nil {
		return nil, err1
	}
	return MongoClient, nil
}

func GetCollection(client *mongo.Client, dbName string, collectionName string) *mongo.Collection {
	collection := client.Database(dbName).Collection(collectionName)
	return collection
}
