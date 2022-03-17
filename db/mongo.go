package db

import (
	"context"
	"log"
	"time"

	"github.com/praveennagaraj97/shopee/pkg/color"
	logger "github.com/praveennagaraj97/shopee/pkg/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Initailizes database connection and returns mongo client
func InitializeDatabaseConnection(mongoUrl, dbName string) *mongo.Client {

	clientOption := options.Client().ApplyURI(mongoUrl)

	client, err := mongo.NewClient(clientOption)
	if err != nil {
		log.Fatal(err)
	}

	// provide connection timout limit for connection to get established.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err = client.Connect(ctx)

	if err != nil {
		log.Fatal(err)
	}

	logger.PrintLog("Connected to MongoDB 🗄", color.Cyan)

	return client
}

func OpenCollection(mgoClient *mongo.Client, dbName, collectionName string) *mongo.Collection {
	return mgoClient.Database(dbName).Collection(collectionName)
}
