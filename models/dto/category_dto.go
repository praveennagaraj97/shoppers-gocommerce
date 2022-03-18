package dto

import (
	"time"

	"github.com/praveennagaraj97/shoppers-gocommerce/pkg/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateCategoryDTO struct {
	Title       string `json:"title" form:"title" bson:"title"`
	Description string `json:"description" form:"description" bson:"description"`
	Published   bool   `json:"published" form:"published" bson:"published"`

	ParentID string `json:"parent_id" form:"parent_id" bson:"-"`
	IconID   string `json:"icon_id" form:"icon_id" bson:"-"`

	PublishedAt *primitive.DateTime `json:"-" form:"-" bson:"published_at"`
	CreatedAt   *primitive.DateTime `json:"-" form:"-" bson:"created_at"`
	UpdatedAt   *primitive.DateTime `json:"-" form:"-" bson:"updated_at"`

	Slug     string                `json:"-" form:"-" bson:"slug"`
	Children []*primitive.ObjectID `json:"-" form:"-" bson:"children"`
	Icon     *primitive.ObjectID   `json:"-" form:"-" bson:"icon"`
	Parent   *primitive.ObjectID   `json:"-" form:"-" bson:"parent"`
}

func (c *CreateCategoryDTO) SetData() error {
	c.Slug = utils.Slugify(c.Title)
	c.Children = []*primitive.ObjectID{}
	ico, err := primitive.ObjectIDFromHex(c.IconID)
	if err != nil {
		return err
	}
	c.Icon = &ico

	if c.ParentID != "" {
		prId, err := primitive.ObjectIDFromHex(c.ParentID)
		if err != nil {
			return err
		}
		c.Parent = &prId
	}

	currentTime := primitive.NewDateTimeFromTime(time.Now())

	c.CreatedAt = &currentTime
	c.UpdatedAt = &currentTime
	c.PublishedAt = nil

	return nil

}

type CategoriesTranslationDTO struct {
	Locale      string              `bson:"locale"`
	Title       string              `bson:"title" json:"title" form:"title"`
	Description string              `bson:"description" json:"description" form:"description"`
	CategoryId  *primitive.ObjectID `bson:"category_id"`
}
