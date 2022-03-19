package dto

import (
	"time"

	"github.com/praveennagaraj97/shoppers-gocommerce/pkg/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateCategoryDTO struct {
	ID          primitive.ObjectID `bson:"_id"`
	Title       string             `json:"title" form:"title" bson:"-"`
	Description string             `json:"description" form:"description" bson:"-"`
	IconID      string             `json:"icon_id" form:"icon_id" bson:"-"`

	Slug        string                `json:"-" form:"-" bson:"slug"`
	Children    []*primitive.ObjectID `json:"-" form:"-" bson:"children"`
	ParentID    string                `json:"parent_id" form:"parent_id" bson:"-"`
	Icon        *primitive.ObjectID   `json:"-" form:"-" bson:"icon"`
	Parent      *primitive.ObjectID   `json:"-" form:"-" bson:"parent"`
	Published   bool                  `json:"-" form:"-" bson:"published"`
	PublishedAt *primitive.DateTime   `json:"-" form:"-" bson:"published_at"`
	CreatedAt   *primitive.DateTime   `json:"-" form:"-" bson:"created_at"`
	UpdatedAt   *primitive.DateTime   `json:"-" form:"-" bson:"updated_at"`
}

func (c *CreateCategoryDTO) SetData() error {
	c.ID = primitive.NewObjectID()
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

	return nil

}

type CategoriesTranslationDTO struct {
	Locale      string              `bson:"locale"`
	Title       string              `bson:"title" json:"title" form:"title"`
	Description string              `bson:"description" json:"description" form:"description"`
	CategoryId  *primitive.ObjectID `bson:"category_id"`
}
