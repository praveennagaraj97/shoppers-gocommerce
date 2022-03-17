package userapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/praveennagaraj97/shoppers-gocommerce/api"
	conf "github.com/praveennagaraj97/shoppers-gocommerce/config"
	"github.com/praveennagaraj97/shoppers-gocommerce/models/dto"
	"github.com/praveennagaraj97/shoppers-gocommerce/models/serialize"
	"github.com/praveennagaraj97/shoppers-gocommerce/pkg/validators"
	userrepository "github.com/praveennagaraj97/shoppers-gocommerce/repository/user"
	useraddress "github.com/praveennagaraj97/shoppers-gocommerce/repository/useraddress"
)

type UserAPI struct {
	config      *conf.GlobalConfiguration
	userRepo    *userrepository.UserRepository
	addressRepo *useraddress.UserAddressRepository
}

// Intialize User API by passing repos and app config.
func (a *UserAPI) InitializeUserAPI(config *conf.GlobalConfiguration,
	userRepo *userrepository.UserRepository,
	addressRepo *useraddress.UserAddressRepository,
) {
	if a.config == nil {
		a.config = config
	}

	if a.userRepo == nil {
		a.userRepo = userRepo
	}

	if a.addressRepo == nil {
		a.addressRepo = addressRepo
	}
}

//Update User
func (a *UserAPI) UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		var payload dto.UpdateUserDTO

		userId, err := api.GetUserIdFromContext(c)
		if err != nil {
			api.SendErrorResponse(a.config.Localize, c, err.Error(), http.StatusUnauthorized, nil)
			return
		}

		// Parse payload
		err = c.ShouldBind(&payload)
		if err != nil {
			api.SendErrorResponse(a.config.Localize, c, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		defer c.Request.Body.Close()

		// check email unique.
		if payload.Email != "" {
			_, err = a.userRepo.FindUserByEmail(payload.Email)
			if err == nil {
				api.SendErrorResponse(a.config.Localize, c, "email_already_taken", http.StatusUnprocessableEntity, nil)
				return
			}
		}

		_, err = a.userRepo.UpdateUser(userId, &payload)

		if err != nil {
			if err != nil {
				api.SendErrorResponse(a.config.Localize, c, err.Error(), http.StatusBadRequest, nil)
				return
			}
		}
		c.JSON(http.StatusNoContent, nil)
	}
}

// Update Password with current password context.
func (a *UserAPI) UpdatePassword() gin.HandlerFunc {
	return func(c *gin.Context) {

		var payload dto.UpdatePasswordDTO

		err := c.ShouldBind(&payload)

		if err != nil {
			api.SendErrorResponse(a.config.Localize, c, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		// validate payload
		if errors := validators.ValidateUpdatePasswordDTO(&payload, a.config.Localize, c); errors != nil {
			api.SendErrorResponse(a.config.Localize, c, "invalid_input", http.StatusUnprocessableEntity, errors)
			return
		}
		defer c.Request.Body.Close()

		userId, err := api.GetUserIdFromContext(c)
		if err != nil {
			api.SendErrorResponse(a.config.Localize, c, "something_went_wrong", http.StatusUnauthorized, nil)
			return
		}

		// check user exists with email
		user, err := a.userRepo.FindUserByID(userId)

		if err != nil {
			api.SendErrorResponse(a.config.Localize, c, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		// check for password match
		err = user.CompareHashedPassword(payload.CurrentPassword)
		if err != nil {
			api.SendErrorResponse(a.config.Localize, c, "invalid_password_entered", http.StatusUnauthorized, nil)
			return
		}

		user.Password = payload.NewPassword
		user.HashPassword()

		err = a.userRepo.UpdateByField(*userId, "password", user.Password)

		if err != nil {

			api.SendErrorResponse(a.config.Localize, c, err.Error(), http.StatusInternalServerError, nil)
			return
		}

		c.JSON(http.StatusOK, serialize.Response{
			StatusCode: http.StatusOK,
			Message:    a.config.Localize.GetMessage("password_updated_successsfully", c),
		})

	}
}

// Get Current User Details
func (a *UserAPI) GetMe() gin.HandlerFunc {
	return func(c *gin.Context) {

		userId, err := api.GetUserIdFromContext(c)
		if err != nil {
			api.SendErrorResponse(a.config.Localize, c, err.Error(), http.StatusUnauthorized, nil)
		}

		user, err := a.userRepo.FindUserByID(userId)
		if err != nil {
			api.SendErrorResponse(a.config.Localize, c, err.Error(), http.StatusNotFound, nil)
		}

		c.JSON(http.StatusOK, serialize.DataResponse{
			Data: user,
			Response: serialize.Response{
				StatusCode: http.StatusOK,
				Message:    "retrieved_user_details",
			},
		})

	}
}
