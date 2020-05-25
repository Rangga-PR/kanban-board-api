package config

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

//ConnectDB : open connection to mongodb atlas database
func ConnectDB(ctx context.Context) (*mongo.Database, context.Context) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://1st-user:apaya123@0-cluster-pxgd6.mongodb.net/kanban?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal("Couldn't connect to the database", err)
	} else {
		log.Println("Connected!")
	}

	return client.Database("kanban"), ctx
}
