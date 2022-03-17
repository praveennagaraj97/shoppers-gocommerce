package useraddressrepository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/praveennagaraj97/shoppers-gocommerce/api"
	"github.com/praveennagaraj97/shoppers-gocommerce/models"
	"github.com/praveennagaraj97/shoppers-gocommerce/models/dto"
	"github.com/praveennagaraj97/shoppers-gocommerce/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserAddressRepository struct {
	collection *mongo.Collection
}

// Method to initialize user repository
func (r *UserAddressRepository) InitializeRepository(colln *mongo.Collection) {
	r.collection = colln

	utils.CreateIndex(colln, bson.D{{Key: "uid", Value: 1}, {Key: "_id", Value: 1}}, "Uid and _id", false)
	utils.CreateIndex(colln, bson.D{{Key: "uid", Value: 1}, {Key: "created_at", Value: 1}}, "Uid and Created At", false)
}

func (r *UserAddressRepository) Create(data *dto.UserAddressDTO, uid *primitive.ObjectID) (*models.AddressModel, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	id := primitive.NewObjectID()
	addressData := r.serializeAddressInput(&id, uid, data)

	_, err := r.collection.InsertOne(ctx, addressData)

	if err != nil {
		return nil, errors.New("something_went_wrong")
	}

	return addressData, nil

}

func (r *UserAddressRepository) FindById(addressId, uid *primitive.ObjectID) (*models.AddressModel, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var data models.AddressModel
	res := r.collection.FindOne(ctx, bson.M{"uid": uid, "_id": addressId})
	if res.Err() != nil {
		return nil, errors.New("no_results_found")
	}
	err := res.Decode(&data)

	if res.Err() != nil {
		return nil, err
	}
	return &data, nil
}

func (r *UserAddressRepository) FindAllByUser(uid *primitive.ObjectID,
	pgnOpt *api.PaginationOptions,
	sortOpts *map[string]int8,
	filterOpts *map[string]primitive.M,
	keySetSortby string,
) ([]models.AddressModel, error) {
	opt := &options.FindOptions{}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	filters := map[string]bson.M{
		"uid": {"$eq": uid},
	}

	if pgnOpt != nil {
		opt.Limit = options.Find().SetLimit(int64(pgnOpt.PerPage)).Limit
		opt.Skip = options.Find().SetSkip(int64((pgnOpt.PageNum - 1) * int(pgnOpt.PerPage))).Skip
	}

	if sortOpts != nil {
		opt.Sort = options.Find().SetSort(sortOpts).Sort
	} else {
		opt.Sort = options.Find().SetSort(bson.M{"created_at": -1}).Sort
	}

	if filterOpts != nil {
		for key, value := range *filterOpts {
			filters[key] = value
		}
	}

	if pgnOpt.PaginateId != nil {
		filters["_id"] = bson.M{keySetSortby: pgnOpt.PaginateId}
	} else {
		opt.Skip = options.Find().SetSkip(int64((pgnOpt.PageNum - 1) * int(pgnOpt.PerPage))).Skip
	}

	cur, err := r.collection.Find(ctx, filters, opt)
	if err != nil {
		return nil, err
	}

	var results []models.AddressModel
	// check for errors in the conversion
	if err = cur.All(ctx, &results); err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	return results, nil
}

func (r *UserAddressRepository) GetDocumentsCount(uid *primitive.ObjectID, filterOpts *map[string]primitive.M) (int64, error) {

	filters := map[string]bson.M{
		"uid": {"$eq": uid},
	}

	if filterOpts != nil {
		for key, value := range *filterOpts {
			filters[key] = value
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	return r.collection.CountDocuments(ctx, filters)
}

func (r *UserAddressRepository) UpdateAddress(uid, addrId *primitive.ObjectID, input *dto.UserAddressDTO) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	addressData := r.serializeAddressInput(addrId, uid, input)
	res, err := r.collection.UpdateOne(ctx, bson.M{"uid": uid, "_id": addrId}, bson.M{"$set": addressData})

	if err != nil || res.ModifiedCount == 0 {
		return errors.New("not_changed")
	}

	return nil

}

func (r *UserAddressRepository) DeleteAddress(uid, addrId *primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	res, err := r.collection.DeleteOne(ctx, bson.M{"uid": uid, "_id": addrId})

	if err != nil || res.DeletedCount == 0 {
		return errors.New("not_deleted")
	}
	return nil
}

func (r *UserAddressRepository) serializeAddressInput(addressId, uid *primitive.ObjectID, data *dto.UserAddressDTO) *models.AddressModel {
	addressData := &models.AddressModel{
		ID:       *addressId,
		Uid:      *uid,
		UserName: fmt.Sprintf("%s %s", data.FirstName, data.LastName),
		Building: data.Building,
		Country: models.Country{
			Name: data.CountryName,
			Code: data.CountryCode,
		},
		Locality: models.Locality{
			Name:          data.LocalityName,
			PostalCode:    data.LocalityPostalCode,
			StreetAddress: data.LocalityStreetAddress,
			City:          data.LocalityCity,
		},
		State:       data.State,
		Phone:       data.Phone,
		AddressType: data.AddressType,
		CreatedAt:   primitive.NewDateTimeFromTime(time.Now()),
		UpdatedAt:   primitive.NewDateTimeFromTime(time.Now()),
	}

	return addressData

}
