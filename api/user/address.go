package userapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/praveennagaraj97/shopee/api"
	"github.com/praveennagaraj97/shopee/models/dto"
	"github.com/praveennagaraj97/shopee/models/serialize"
	"github.com/praveennagaraj97/shopee/pkg/validators"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Add User User
func (a *UserAPI) AddNewAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get user id from req
		userId, err := api.GetUserIdFromContext(c)
		if err != nil {
			api.SendErrorResponse(a.config.Localize, c, err.Error(), http.StatusUnauthorized, nil)
			return
		}

		// get address body
		var payload dto.UserAddressDTO

		err = c.ShouldBind(&payload)

		if err != nil {
			api.SendErrorResponse(a.config.Localize, c, "invalid_input", http.StatusUnprocessableEntity, nil)
			return
		}
		defer c.Request.Body.Close()

		if errors := validators.ValidateAddressInput(&payload, a.config.Localize, c); errors != nil {
			api.SendErrorResponse(a.config.Localize, c, "invalid_input", http.StatusUnprocessableEntity, errors)
			return
		}

		res, err := a.addressRepo.Create(&payload, userId)
		if err != nil {
			api.SendErrorResponse(a.config.Localize, c, err.Error(), http.StatusInternalServerError, nil)
			return
		}

		c.JSON(http.StatusCreated, serialize.DataResponse{
			Data: res,
			Response: serialize.Response{
				StatusCode: http.StatusCreated,
				Message:    "address_added_successfully",
			},
		})

	}
}

// get individual address by id
func (a *UserAPI) GetAddressById() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, err := api.GetUserIdFromContext(c)
		if err != nil {
			api.SendErrorResponse(a.config.Localize, c, err.Error(), http.StatusUnauthorized, nil)
			return
		}

		// get request id from param
		addressId := c.Param("id")

		addId, err := primitive.ObjectIDFromHex(addressId)
		if err != nil {
			api.SendErrorResponse(a.config.Localize, c, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		res, err := a.addressRepo.FindById(&addId, userId)
		if err != nil {
			api.SendErrorResponse(a.config.Localize, c, err.Error(), http.StatusNotFound, nil)
			return
		}

		c.JSON(http.StatusOK, serialize.DataResponse{
			Data: res,
			Response: serialize.Response{
				StatusCode: http.StatusOK,
				Message:    "address_retrieved_successfully",
			},
		})

	}
}

// Get all address with filter,sorting and pagination
func (a *UserAPI) GetListOfUserAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		// user context
		userId, err := api.GetUserIdFromContext(c)
		if err != nil {
			api.SendErrorResponse(a.config.Localize, c, err.Error(), http.StatusUnauthorized, nil)
			return
		}

		// get pagination/sort/filter options.
		pgOpts := api.ParsePaginationOptions(c)
		srtOpts := api.ParseSortByOptions(c)
		filterOpts := api.ParseFilterByOptions(c)
		keySetSortby := "$gt"

		// Default options | sort by latest
		if len(*srtOpts) == 0 {
			srtOpts = &map[string]int8{"_id": -1}
		}
		// Key Set fix for created_at desc
		if pgOpts.PaginateId != nil {
			for key, value := range *srtOpts {
				if value == -1 && key == "_id" {
					keySetSortby = "$lt"
				}
			}
		}

		res, err := a.addressRepo.FindAllByUser(userId, pgOpts, srtOpts, filterOpts, keySetSortby)
		if err != nil {
			api.SendErrorResponse(a.config.Localize, c, err.Error(), http.StatusInternalServerError, nil)
			return
		}

		resLen := len(res)

		// Paginate Options
		var docCount int64
		var lastResId *primitive.ObjectID

		if pgOpts.PaginateId == nil {
			docCount, err = a.addressRepo.GetDocumentsCount(userId, filterOpts)
			if err != nil {
				api.SendErrorResponse(a.config.Localize, c, err.Error(), http.StatusInternalServerError, nil)
				return
			}
		}

		if resLen > 0 {
			lastResId = &res[resLen-1].ID
		}

		count, next, prev, paginateKeySetID := api.GetPaginateOptions(docCount, pgOpts, int64(resLen), lastResId)

		c.JSON(http.StatusOK, serialize.PaginatedDataResponse{
			Count:            count,
			Next:             next,
			Prev:             prev,
			PaginateKeySetID: paginateKeySetID,
			DataResponse: serialize.DataResponse{
				Data: res,
				Response: serialize.Response{
					StatusCode: http.StatusOK,
					Message:    "address_retrieved_successfully",
				},
			},
		})
	}
}

// Update address
func (a *UserAPI) UpdateUserAddressByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// user context
		userId, err := api.GetUserIdFromContext(c)
		if err != nil {
			api.SendErrorResponse(a.config.Localize, c, err.Error(), http.StatusUnauthorized, nil)
			return
		}

		// get request id from param
		addressId := c.Param("id")

		addrId, err := primitive.ObjectIDFromHex(addressId)
		if err != nil {
			api.SendErrorResponse(a.config.Localize, c, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		var payload dto.UserAddressDTO

		c.ShouldBind(&payload)

		err = a.addressRepo.UpdateAddress(userId, &addrId, &payload)

		if err != nil {
			api.SendErrorResponse(a.config.Localize, c, err.Error(), http.StatusBadRequest, nil)
			return
		}
		c.JSON(http.StatusNoContent, nil)
	}
}

// delete user address
func (a *UserAPI) DeleteUserAddressByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// user context
		userId, err := api.GetUserIdFromContext(c)
		if err != nil {
			api.SendErrorResponse(a.config.Localize, c, err.Error(), http.StatusUnauthorized, nil)
			return
		}

		// get request id from param
		addressId := c.Param("id")

		addrId, err := primitive.ObjectIDFromHex(addressId)
		if err != nil {
			api.SendErrorResponse(a.config.Localize, c, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		err = a.addressRepo.DeleteAddress(userId, &addrId)
		if err != nil {
			api.SendErrorResponse(a.config.Localize, c, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}
		c.JSON(http.StatusNoContent, nil)
	}
}
