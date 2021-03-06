package categoriesapi

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/praveennagaraj97/shoppers-gocommerce/api"
	conf "github.com/praveennagaraj97/shoppers-gocommerce/config"
	"github.com/praveennagaraj97/shoppers-gocommerce/constants"
	"github.com/praveennagaraj97/shoppers-gocommerce/models/dto"
	"github.com/praveennagaraj97/shoppers-gocommerce/models/serialize"
	categoriesrepository "github.com/praveennagaraj97/shoppers-gocommerce/repository/categories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CategoriesAPI struct {
	conf *conf.GlobalConfiguration
	repo *categoriesrepository.CategoriesRepository
}

func (a *CategoriesAPI) Initialize(cfg *conf.GlobalConfiguration, repo *categoriesrepository.CategoriesRepository) {
	a.conf = cfg
	a.repo = repo

}

// Add new category - requires admin rights to add new category.
func (a *CategoriesAPI) AddCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		var payload dto.CreateCategoryDTO

		if err = c.ShouldBind(&payload); err != nil {
			api.SendErrorResponse(a.conf.Localize, c, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}
		defer c.Request.Body.Close()

		if payload.Title == "" {
			api.SendErrorResponse(a.conf.Localize, c, "category_title_is_required", http.StatusUnprocessableEntity, nil)
			return
		}

		if payload.IconID != "" {
			icon, err := primitive.ObjectIDFromHex(payload.IconID)
			if err != nil {
				api.SendErrorResponse(a.conf.Localize, c, err.Error(), http.StatusUnprocessableEntity, nil)
				return
			}
			payload.Icon = &icon
		}

		res, err := a.repo.Create(&payload, a.conf.DefaultLocale)
		if err != nil {
			api.SendErrorResponse(a.conf.Localize, c, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		c.JSON(http.StatusCreated, serialize.DataResponse{
			Data: res, Response: serialize.Response{
				StatusCode: http.StatusCreated,
				Message:    a.conf.Localize.GetMessage("category_added_successfully", c),
			},
		})
	}
}

// Add Translation for Category
func (a *CategoriesAPI) AddCategoryTranslations() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		var payload dto.CategoriesTranslationDTO

		locale, _ := c.Params.Get("locale")
		category, _ := c.Params.Get("category")

		// Check if default locale is selected
		if locale == a.conf.DefaultLocale {
			api.SendErrorResponse(a.conf.Localize, c, "default_locale_is_not_allowed", http.StatusUnprocessableEntity, nil)
			return
		}

		if err = c.ShouldBind(&payload); err != nil {
			api.SendErrorResponse(a.conf.Localize, c, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		defer c.Request.Body.Close()

		if payload.Title == "" {
			api.SendErrorResponse(a.conf.Localize, c, "category_title_is_required", http.StatusUnprocessableEntity, nil)
			return
		}

		if category == "" {
			api.SendErrorResponse(a.conf.Localize, c, "categoryid_is_required", http.StatusUnprocessableEntity, nil)
			return
		}

		categoryID, err := primitive.ObjectIDFromHex(category)

		if err != nil {
			api.SendErrorResponse(a.conf.Localize, c, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		payload.CategoryId = &categoryID
		payload.Locale = locale

		err = a.repo.AddTranslations(&payload)
		if err != nil {
			api.SendErrorResponse(a.conf.Localize, c, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		c.JSON(http.StatusCreated, serialize.Response{
			StatusCode: http.StatusCreated,
			Message:    a.conf.Localize.GetMessage("translation_added_successfully", c),
		})

	}
}

func (a *CategoriesAPI) GetAllCategories() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		loc, _ := ctx.Get(string(constants.CUSTOME_HEADER_LANG_KEY))

		locale := fmt.Sprint(loc)

		res, err := a.repo.FindAll(locale)

		if err != nil {
			api.SendErrorResponse(a.conf.Localize, ctx, err.Error(), 500, nil)
			return
		}

		ctx.JSON(http.StatusOK, serialize.DataResponse{
			Data: res,
			Response: serialize.Response{
				StatusCode: http.StatusOK,
				Message:    "list_of_categories_retrieved",
			},
		})
	}
}

func (a *CategoriesAPI) MarkPublishStatus(status bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cId, _ := ctx.Params.Get("category")

		categoryId, err := primitive.ObjectIDFromHex(cId)

		if err != nil {
			api.SendErrorResponse(a.conf.Localize, ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		// check if translations exist for all locales
		c, err := a.repo.GetAvailableTranslationsCount(&categoryId)

		if err != nil {
			api.SendErrorResponse(a.conf.Localize, ctx, err.Error(), http.StatusBadRequest, nil)
			return
		}

		if int(c) != len(*a.conf.Locales) {
			api.SendErrorResponse(a.conf.Localize, ctx, "one_or_more_translations_are_missing", http.StatusBadRequest, nil)
			return
		}

		err = a.repo.MarkPublishedStatus(&categoryId, status)

		if err != nil {
			api.SendErrorResponse(a.conf.Localize, ctx, err.Error(), http.StatusBadRequest, nil)
			return
		}

	}
}
