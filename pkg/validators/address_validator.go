package validators

import (
	"github.com/gin-gonic/gin"
	"github.com/praveennagaraj97/shoppers-gocommerce/models/dto"
	"github.com/praveennagaraj97/shoppers-gocommerce/pkg/i18n"
)

func ValidateAddressInput(payload *dto.UserAddressDTO,
	localizer *i18n.Internationalization,
	c *gin.Context) map[string]interface{} {

	var errors map[string]interface{} = make(map[string]interface{})

	if payload.AddressType == "" {
		errors["address_type"] = localizer.GetMessage("address_type_is_required", c)
	}

	if payload.CountryName == "" {
		errors["country_name"] = localizer.GetMessage("country_name_is_required", c)
	}

	if payload.CountryCode == "" {
		errors["country_code"] = localizer.GetMessage("country_code_is_required", c)
	}

	if payload.Phone == "" {
		errors["phone"] = localizer.GetMessage("phone_number_is_required", c)
	}

	if payload.FirstName == "" {
		errors["first_name"] = localizer.GetMessage("first_name_is_required", c)
	}

	if len(errors) == 0 {
		return nil
	}

	return errors
}
