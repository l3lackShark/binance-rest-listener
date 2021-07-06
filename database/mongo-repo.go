package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//This repository implementation is not really reliable if the user wants to use multiple drivers (postgresql/etc) Made for simplicity

type (
	MongoString string //to prevent accidental string input, those should be pre defined or explicitly typecasted.

	Repository interface {
	}

	database struct { //internal struct
		client *mongo.Client
	}
)

var (
	DaatabaseName  MongoString = "blockchain"
	CollectionName MongoString = "binance"
)

func New(connURL string) (Repository, error) {

	clientOptions := options.Client().ApplyURI(connURL)
	clientOptions = clientOptions.SetMaxPoolSize(50)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	dbClient, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return &database{}, fmt.Errorf("failed to connect to mongodb: %e", err)
	}
	log.Println("Connected to mongodb database") //TODO: replace with custom logger
	return &database{
		client: dbClient,
	}, nil
}
