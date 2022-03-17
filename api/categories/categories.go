package categoriesapi

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/praveennagaraj97/shopee/api"
	conf "github.com/praveennagaraj97/shopee/config"
	"github.com/praveennagaraj97/shopee/constants"
	"github.com/praveennagaraj97/shopee/models/dto"
	"github.com/praveennagaraj97/shopee/models/serialize"
	categoriesrepository "github.com/praveennagaraj97/shopee/repository/categories"
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

		// payload.Locale = a.conf.DefaultLocale

		res, err := a.repo.Create(&payload)
		if err != nil {
			api.SendErrorResponse(a.conf.Localize, c, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		c.JSON(http.StatusCreated, serialize.DataResponse{Data: res, Response: serialize.Response{
			StatusCode: http.StatusCreated,
			Message:    a.conf.Localize.GetMessage("category_added_successfully", c),
		}})

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

		res, err := a.repo.FindAll(locale, a.conf.DefaultLocale)

		if err != nil {
			api.SendErrorResponse(a.conf.Localize, ctx, err.Error(), 500, nil)
			return
		}

		ctx.JSON(200, map[string]interface{}{
			"data": res,
		})

	}
}
