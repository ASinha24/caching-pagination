package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Dbconnect() *mongo.Collection {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("failed to create client err: %s", err)
	}

	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatalf("failed to connect client err: %s", err)
	}
	//defer client.Disconnect(context.TODO())

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("failed to create ping err: %s", err)

	}
	log.Println("document created.....")
	collection := client.Database("mydb").Collection("mart")
	return collection
}
