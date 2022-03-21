package categoriesrepository

import (
	"context"
	"errors"
	"time"

	"github.com/praveennagaraj97/shoppers-gocommerce/models"
	"github.com/praveennagaraj97/shoppers-gocommerce/models/dto"
	"github.com/praveennagaraj97/shoppers-gocommerce/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	utils.CreateIndex(col, bson.D{{Key: "published", Value: 1}}, "Published", false)
	utils.CreateIndex(r.translationCollection, bson.D{{Key: "locale", Value: 1}}, "Locale", false)
	utils.CreateIndex(r.translationCollection, bson.D{{Key: "category_id", Value: 1}}, "Category Reference", false)

}

func (r *CategoriesRepository) Create(data *dto.CreateCategoryDTO, defaultLocale string) (*models.CategoriesModel, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := data.SetData(); err != nil {
		return nil, err
	}

	//Check if exists
	exists := r.FindBySlug(data.Slug)

	if exists {
		return nil, errors.New("category_with_this_title_already_exist")
	}

	_, err := r.repo.InsertOne(ctx, data)
	if err != nil {
		return nil, err
	}

	// Add Translation for Default Locale
	err = r.AddTranslations(
		&dto.CategoriesTranslationDTO{
			Locale:      defaultLocale,
			Title:       data.Title,
			Description: data.Description,
			CategoryId:  &data.ID,
		})

	if err != nil {
		r.repo.DeleteOne(ctx, bson.M{"_id": data.ID})
		return nil, err
	}

	return &models.CategoriesModel{
		ID:          &data.ID,
		Slug:        &data.Slug,
		Locale:      &defaultLocale,
		Title:       &data.Title,
		Description: &data.Description,
		Published:   &data.Published,
		PublishedAt: data.PublishedAt,
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.UpdatedAt,
		Parent:      data.Parent,
		Children:    data.Children,
	}, nil

}

func (r *CategoriesRepository) FindBySlug(slug string) bool {
	res := r.repo.FindOne(context.Background(), bson.D{{Key: "slug", Value: slug}})

	if res.Err() != nil {
		return false
	}
	return true
}

func (r *CategoriesRepository) FindAll(locale string) ([]models.CategoriesModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	matchPublishedStage := bson.D{{Key: "$match", Value: bson.M{"published": true}}}
	iconLookUpStage := bson.D{{Key: "$lookup", Value: bson.M{
		"from":         "assets",
		"localField":   "icon",
		"foreignField": "_id",
		"as":           "icon",
	}}}
	iconUnWindStage := bson.D{{Key: "$unwind", Value: bson.M{
		"preserveNullAndEmptyArrays": true,
		"path":                       "$icon",
	}}}

	translationLookUpStage := bson.D{{Key: "$lookup", Value: bson.M{
		"from":         "categories_translations",
		"localField":   "_id",
		"foreignField": "category_id",
		"as":           "translation",
		"pipeline": bson.A{
			bson.M{
				"$match": bson.M{"locale": locale},
			}},
	}}}
	translationUnWindStage := bson.D{{Key: "$unwind", Value: bson.M{
		"path":                       "$translation",
		"preserveNullAndEmptyArrays": false,
	}}}

	projectFieldsStage := bson.D{{Key: "$project", Value: bson.M{
		"_id":          1,
		"slug":         1,
		"locale":       "$translation.locale",
		"title":        "$translation.title",
		"description":  "$translation.description",
		"icon":         1,
		"published":    1,
		"published_at": 1,
		"updated_at":   1,
		"created_at":   1,
		"children":     1,
	}}}

	cur, err := r.repo.Aggregate(ctx, mongo.Pipeline{
		matchPublishedStage,
		iconLookUpStage,
		iconUnWindStage,
		translationLookUpStage,
		translationUnWindStage,
		projectFieldsStage,
	})

	if err != nil {
		return nil, err
	}

	var results []models.CategoriesModel

	if err = cur.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil

}

func (r *CategoriesRepository) MarkPublishedStatus(id *primitive.ObjectID, status bool) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	updateRes, err := r.repo.UpdateByID(ctx, id, bson.M{"$set": bson.M{"published": status}})
	if err != nil {
		return err
	}

	if updateRes.ModifiedCount == 0 {
		return errors.New("document_not_modified")
	}

	return nil

}
