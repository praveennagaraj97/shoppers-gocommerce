package validators

import (
	"github.com/gin-gonic/gin"
	"github.com/praveennagaraj97/shoppers-gocommerce/models/dto"
	"github.com/praveennagaraj97/shoppers-gocommerce/pkg/i18n"
)

// Validate SignUp Input checks fields for value and also validates email.
func ValidateSignUpData(payload *dto.CreateUserDTO, localizer *i18n.Internationalization, c *gin.Context) map[string]interface{} {

	var errors map[string]interface{} = make(map[string]interface{})

	// check for empty field
	if payload.Email == "" || payload.Password == "" || payload.LastName == "" || payload.FirstName == "" {
		errors["fields"] = localizer.GetMessage("one_or_more_fileds_missing", c)
	}

	// validate email
	if err := validateEmail(payload.Email); err != nil {
		errors["email"] = []string{localizer.GetMessage("provided_email_is_invalid", c)}
	}

	if len(errors) == 0 {
		return nil
	}

	return errors
}

// Validate Login Input checks fields for value and also validates email.
func ValidateLoginInput(payload *dto.LoginDTO, localizer *i18n.Internationalization, c *gin.Context) map[string]interface{} {
	var errors map[string]interface{} = make(map[string]interface{})

	// check for empty field
	if payload.Email == "" || payload.Password == "" {
		errors["fields"] = localizer.GetMessage("one_or_more_fileds_missing", c)
	}

	// validate email
	if err := validateEmail(payload.Email); err != nil {
		errors["email"] = []string{localizer.GetMessage("provided_email_is_invalid", c)}
	}

	if len(errors) == 0 {
		return nil
	}

	return errors
}

func ValidateUpdatePasswordDTO(payload *dto.UpdatePasswordDTO, localizer *i18n.Internationalization, c *gin.Context) map[string]interface{} {
	var errors map[string]interface{} = make(map[string]interface{})

	// check for empty field
	if payload.CurrentPassword == "" || payload.NewPassword == "" {
		errors["fields"] = localizer.GetMessage("one_or_more_fileds_missing", c)
		return errors
	}

	// same password
	if payload.CurrentPassword == payload.NewPassword {
		errors["fields"] = localizer.GetMessage("old_password_and_new_password_cannot_be_same", c)
		return errors
	}

	return nil

}
