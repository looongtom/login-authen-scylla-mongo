package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

func ConnectUser(collectionName string) *mongo.Collection {
	clientOptions := options.Client().ApplyURI(os.Getenv("MongoURI"))
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		fmt.Println(err.Error())
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("Connected to MongoDB!")

	collection := client.Database(os.Getenv("DatabaseMongo")).Collection(collectionName)

	return collection
}
