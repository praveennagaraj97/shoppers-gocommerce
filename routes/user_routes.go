package routes

import (
	userapi "github.com/praveennagaraj97/shopee/api/user"
)

func (r *Router) userRoutes() {
	route := r.engine.Group("/api/v1/user")

	// user repo

	userAPI := userapi.UserAPI{}

	userAPI.InitializeUserAPI(r.conf, r.repos.GetUserRepo(), r.repos.GetAddressRepo())

	// Auth
	route.POST("/signup", userAPI.Register("user"))
	route.GET("/confirm-email", userAPI.ConfirmEmailAddress())
	route.POST("/login", userAPI.Login())
	route.GET("/refresh", userAPI.RefreshToken())
	// route.POST("/forgot-password", userAPI.ForgotPassword())
	route.GET("/logout", r.middlewares.IsAuthorized(), userAPI.Logout())
	// route.POST("/reset-password", userAPI.ResetPassword())
	route.POST("/add-admin", r.middlewares.IsAuthorized(), r.middlewares.UserRole([]string{"super_admin"}), userAPI.Register("admin"))

	// Profile
	route.Use(r.middlewares.IsAuthorized())
	route.GET("", userAPI.GetMe())
	route.PATCH("", userAPI.UpdateUser())
	route.POST("/update_password", userAPI.UpdatePassword())

	// Address
	addressRouteGroup := route.Group("/address")
	addressRouteGroup.Use(r.middlewares.UserRole([]string{"user"}))
	addressRouteGroup.POST("", userAPI.AddNewAddress())
	addressRouteGroup.GET("", userAPI.GetListOfUserAddress())
	addressRouteGroup.GET("/:id", userAPI.GetAddressById())
	addressRouteGroup.PUT("/:id", userAPI.UpdateUserAddressByID())
	addressRouteGroup.DELETE("/:id", userAPI.DeleteUserAddressByID())
}
