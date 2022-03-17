package utils

import (
	"context"
	"os"
	"time"

	"github.com/h2non/bimg"
	logger "github.com/praveennagaraj97/shoppers-gocommerce/pkg/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateFolder(dirName string) error {
	// check if folder already exists
	_, err := os.Stat(dirName)
	if os.IsNotExist(err) {
		// create if it doesn't exist
		errDir := os.MkdirAll(dirName, 0755)
		if errDir != nil {
			return errDir
		}
	}
	return nil
}

func CreateBlurDataForImages(buffer []byte, quality int, width int, height int) ([]byte, error) {

	processed, err := bimg.NewImage(buffer).Process(bimg.Options{
		Quality:       quality,
		Force:         true,
		StripMetadata: true,
		Crop:          true,
		Lossless:      false,
		Width:         width,
		Height:        height,
	})
	if err != nil {
		return nil, err
	}

	return processed, nil
}

type IndexFieldsOptions struct {
	Field string
	Order int8
}

// CreateIndex - creates an index for a collection
func CreateIndex(collection *mongo.Collection, keys bson.D, indexName string, unique bool) bool {

	var indexOptions *options.IndexOptions = &options.IndexOptions{}

	indexOptions.Unique = &unique
	indexOptions.Name = options.Index().SetName(indexName).Name

	// 1. Field key
	mod := mongo.IndexModel{
		Keys:    keys,
		Options: indexOptions,
	}

	// 2. Create the context for this operation
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 3. Create a single index
	_, err := collection.Indexes().CreateOne(ctx, mod)
	if err != nil {

		logger.ErrorLogFatal(err.Error())
		return false
	}

	return true
}
