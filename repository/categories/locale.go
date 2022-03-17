package categoriesrepository

import (
	"context"
	"time"

	"github.com/praveennagaraj97/shopee/models/dto"
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
