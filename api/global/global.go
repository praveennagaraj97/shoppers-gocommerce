package globalapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	conf "github.com/praveennagaraj97/shoppers-gocommerce/config"
	"github.com/praveennagaraj97/shoppers-gocommerce/models/serialize"
)

type GlobalAPI struct {
	conf *conf.GlobalConfiguration
}

func (a *GlobalAPI) Initialize(cfg *conf.GlobalConfiguration) {
	a.conf = cfg
}

func (a *GlobalAPI) GetLocales() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		ctx.JSON(http.StatusOK, serialize.DataResponse{
			Data: a.conf.Locales,
			Response: serialize.Response{
				StatusCode: http.StatusOK,
				Message:    "locales_retrieved",
			},
		})

	}
}
