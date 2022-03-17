package routes

import globalapi "github.com/praveennagaraj97/shoppers-gocommerce/api/global"

func (r *Router) globalRoutes() {
	router := r.engine.Group("/api/v1/global")

	apis := globalapi.GlobalAPI{}
	apis.Initialize(r.conf)

	router.GET("/locales", apis.GetLocales())
}
