package assetrepository

import (
	"context"
	"errors"
	"time"

	"github.com/praveennagaraj97/shoppers-gocommerce/api"
	"github.com/praveennagaraj97/shoppers-gocommerce/models"
	awspkg "github.com/praveennagaraj97/shoppers-gocommerce/pkg/aws"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AssetRepository struct {
	repo     *mongo.Collection
	s3Client *awspkg.AWSConfiguration
}

func (a *AssetRepository) Initialize(rp *mongo.Collection) {
	a.repo = rp

	// Index Options
	// utils.CreateIndex(a.repo, "Created At - DESC", "created_at", -1, false)

}

func (a *AssetRepository) AddNewAsset(data *models.AssetModel) (*models.AssetModel, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	id := primitive.NewObjectID()

	data.ID = id
	_, err := a.repo.InsertOne(ctx, data)

	if err != nil {
		return nil, err
	}
	return data, nil
}

//Supports Last Id based pagination, Sorting by fields and filter by field.
func (a *AssetRepository) FindAll(
	pgOpt *api.PaginationOptions,
	sortOption *map[string]int8,
	filterOptions *map[string]primitive.M,
	// $gt | $lt
	keySetSortBy string,
) ([]*models.AssetModel, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	queryOptions := options.FindOptions{}
	filters := map[string]bson.M{}

	if pgOpt.PaginateId != nil {
		filters["_id"] = bson.M{keySetSortBy: pgOpt.PaginateId}
	} else {
		queryOptions.Skip = options.Find().SetSkip(int64((pgOpt.PageNum - 1) * pgOpt.PerPage)).Skip
	}

	if pgOpt != nil {
		queryOptions.Limit = options.Find().SetLimit(int64(pgOpt.PerPage)).Limit
	}

	if filterOptions != nil {
		for key, value := range *filterOptions {
			filters[key] = value
		}
	}

	if sortOption != nil {
		queryOptions.Sort = queryOptions.SetSort(sortOption).Sort
	}

	cursor, err := a.repo.Find(ctx, filters, &queryOptions)

	if err != nil {
		return nil, err
	}

	var data []*models.AssetModel
	cursor.All(ctx, &data)
	defer cursor.Close(ctx)
	return data, nil

}

func (a *AssetRepository) GetDocumentsCount(filterOptions *map[string]primitive.M) (int64, error) {
	var filters map[string]primitive.M = make(map[string]primitive.M)

	if filterOptions != nil {
		for key, value := range *filterOptions {
			filters[key] = value
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)

	defer cancel()
	return a.repo.CountDocuments(ctx, filterOptions)

}

func (a *AssetRepository) FindByID(id *primitive.ObjectID) (*models.AssetModel, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	res := a.repo.FindOne(ctx, bson.M{"_id": id})
	var data *models.AssetModel
	err := res.Decode(&data)
	if err != nil {
		return nil, errors.New("no_results_found")
	}
	return data, nil
}

func (a *AssetRepository) AddNewLinkedEntity(assetID *primitive.ObjectID, entity *models.AssetLinkedWithInfo) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	return a.repo.UpdateOne(ctx, bson.M{"_id": assetID}, bson.M{"$push": bson.M{"linked_entities": entity}})

}
