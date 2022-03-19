package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CategoriesTranslationModel struct {
	Locale      string              `bson:"locale"`
	Title       string              `bson:"title" json:"title" form:"title"`
	Description string              `bson:"description" json:"description" form:"description"`
	CategoryId  *primitive.ObjectID `bson:"category_id"`
}

type CategoriesModel struct {
	ID          *primitive.ObjectID   `json:"_id" bson:"_id"`
	Slug        *string               `json:"slug" bson:"slug"`
	Locale      *string               `json:"locale" bson:"locale"`
	Title       *string               `json:"title" bson:"title"`
	Description *string               `json:"description" bson:"description"`
	Icon        *AssetModel           `json:"icon" bson:"icon"`
	Published   *bool                 `json:"published" bson:"published"`
	PublishedAt *primitive.DateTime   `json:"published_at" bson:"published_at"`
	CreatedAt   *primitive.DateTime   `json:"created_at" bson:"created_at"`
	UpdatedAt   *primitive.DateTime   `json:"updated_at" bson:"updated_at"`
	Parent      *primitive.ObjectID   `json:"parent" bson:"parent"`
	Children    []*primitive.ObjectID `json:"children" bson:"children"`
}
