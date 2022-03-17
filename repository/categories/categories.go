package categoriesrepository

import (
	"context"
	"errors"
	"time"

	"github.com/praveennagaraj97/shoppers-gocommerce/models"
	"github.com/praveennagaraj97/shoppers-gocommerce/models/dto"
	"github.com/praveennagaraj97/shoppers-gocommerce/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CategoriesRepository struct {
	repo                  *mongo.Collection
	translationCollection *mongo.Collection
}

func (r *CategoriesRepository) Initialize(col *mongo.Collection) {
	r.repo = col
	r.initTranslationsCollection()

	utils.CreateIndex(col, bson.D{{Key: "slug", Value: 1}}, "Category Slug", true)

}

func (r *CategoriesRepository) Create(data *dto.CreateCategoryDTO) (*models.CategoriesModel, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Check if exists
	exists := r.FindBySlug(data.Slug)

	if exists {
		return nil, errors.New("category_with_this_title_already_exist")
	}

	_, err := r.repo.InsertOne(ctx, data)
	if err != nil {
		return nil, err
	}

	// Add Translation for Default Locale
	// err = r.AddTranslations(
	// 	&dto.CategoriesTranslationDTO{
	// 		Locale:      data.Locale,
	// 		Title:       data.Title,
	// 		Description: data.Description,
	// 		CategoryId:  &dataModel.ID,
	// 	})

	// if err != nil {
	// 	r.repo.DeleteOne(ctx, bson.M{"_id": dataModel.ID})
	// 	return nil, err
	// }

	return nil, nil

}

func (r *CategoriesRepository) FindBySlug(slug string) bool {
	res := r.repo.FindOne(context.Background(), bson.D{{Key: "slug", Value: slug}})

	if res.Err() != nil {
		return false
	}
	return true
}

func (r *CategoriesRepository) FindAll(locale string, defaultLocale string) ([]models.CategoriesModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	matchStage := bson.D{{Key: "$match", Value: bson.M{"locale": locale}}}

	lookUpStage := bson.D{{Key: "$lookup", Value: bson.M{
		"from":         "categories",
		"localField":   "category_id",
		"foreignField": "_id",
		"as":           "category",
	},
	}}

	unWindStage := bson.D{{Key: "$unwind", Value: "$category"}, {Key: "preserveNullAndEmptyArrays", Value: true}}

	projectStage := bson.D{{Key: "$project", Value: bson.M{}}}

	cur, err := r.translationCollection.Aggregate(ctx, mongo.Pipeline{
		matchStage,
		lookUpStage,
		unWindStage,
		projectStage,
	})

	if err != nil {
		return nil, err
	}

	var resp []models.CategoriesModel

	err = cur.All(ctx, &resp)
	if err != nil {
		return nil, err
	}

	cur.Close(ctx)

	return resp, nil

}
