package routes

import assetsapi "github.com/praveennagaraj97/shopee/api/assets"

func (r *Router) assetsRoutes() {

	assetApi := assetsapi.AssetsAPI{}
	assetApi.Initialize(r.repos.GetAssetRepo(), r.conf)

	router := r.engine.Group("/api/v1/assets")

	router.POST("/upload", r.middlewares.IsAuthorized(), assetApi.UploadSingleAsset())
	router.GET("", r.middlewares.IsAuthorized(), r.middlewares.UserRole([]string{"admin", "super_admin"}), assetApi.GetAllAssets())
	router.GET("/:id", r.middlewares.IsAuthorized(), r.middlewares.UserRole([]string{"admin", "super_admin"}), assetApi.GetAssetByID())

}
