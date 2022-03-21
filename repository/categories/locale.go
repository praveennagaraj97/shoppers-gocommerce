package categoriesrepository

import (
	"context"
	"time"

	"github.com/praveennagaraj97/shoppers-gocommerce/models/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (r *CategoriesRepository) initTranslationsCollection() {
	r.translationCollection = r.repo.Database().Collection("categories_translations")

}

func (r *CategoriesRepository) AddTranslations(data *dto.CategoriesTranslationDTO) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_, err := r.translationCollection.InsertOne(ctx, data)
	return err
}

func (r *CategoriesRepository) GetAvailableTranslationsCount(category *primitive.ObjectID) (uint8, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	count, err := r.translationCollection.CountDocuments(ctx, bson.M{"category_id": bson.M{"$eq": category}})
	if err != nil {
		return 0, err
	}

	return uint8(count), nil
}
