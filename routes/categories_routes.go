package routes

import (
	categoriesapi "github.com/praveennagaraj97/shoppers-gocommerce/api/categories"
)

func (r *Router) categoriesRoutes() {
	router := r.engine.Group("/api/v1/categories")

	// categories api
	categoryApi := categoriesapi.CategoriesAPI{}
	categoryApi.Initialize(r.conf, r.repos.GetCategoriesRepo())

	router.POST("", r.middlewares.IsAuthorized(), r.middlewares.UserRole([]string{"admin"}), categoryApi.AddCategory())
	router.POST("/translations/:locale/:category", r.middlewares.IsAuthorized(), r.middlewares.UserRole([]string{"admin"}), categoryApi.AddCategoryTranslations())
	router.GET("", categoryApi.GetAllCategories())

}
