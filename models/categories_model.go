package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CategoriesModel struct {
	ID          *primitive.ObjectID   `json:"_id"`
	Slug        *string               `json:"slug"`
	Locale      *string               `json:"locale"`
	Title       *string               `json:"title"`
	Description *string               `json:"description"`
	Icon        *primitive.ObjectID   `json:"icon"`
	Published   *bool                 `json:"published"`
	PublishedAt *string               `json:"published_at"`
	CreatedAt   primitive.DateTime    `json:"created_at"`
	UpdatedAt   *string               `json:"updated_at"`
	Parent      *primitive.ObjectID   `json:"parent"`
	Children    []*primitive.ObjectID `json:"children"`
}
