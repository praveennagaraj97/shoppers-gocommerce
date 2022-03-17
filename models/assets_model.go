package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FileMetaModel struct {
	Key         string `json:"-" bson:"key"`
	Title       string `json:"title" bson:"title"`
	Description string `json:"description" bson:"description"`
}

type AssetLinkedWithInfo struct {
	Title string
	ID    *primitive.ObjectID
	Model string
}

type AssetModel struct {
	ID           primitive.ObjectID     `json:"_id" bson:"_id"`
	OriginalURL  string                 `json:"original_url" bson:"original_url"`
	BlurDataURL  string                 `json:"blur_data_url,omitempty" bson:"blur_data_url,omitempty"`
	MetaData     *FileMetaModel         `json:"meta" bson:"meta"`
	ContentType  string                 `json:"content_type" bson:"content_type"`
	CreatedAt    primitive.DateTime     `json:"created_at" bson:"created_at"`
	Published    bool                   `json:"published" bson:"published"`
	LinkedEntity []*AssetLinkedWithInfo `json:"-" bson:"linked_entities"`
}
