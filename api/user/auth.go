package userapi

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/praveennagaraj97/shoppers-gocommerce/api"
	"github.com/praveennagaraj97/shoppers-gocommerce/constants"
	"github.com/praveennagaraj97/shoppers-gocommerce/models"
	"github.com/praveennagaraj97/shoppers-gocommerce/models/dto"
	"github.com/praveennagaraj97/shoppers-gocommerce/models/serialize"
	"github.com/praveennagaraj97/shoppers-gocommerce/pkg/tokens"
	"github.com/praveennagaraj97/shoppers-gocommerce/pkg/validators"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Register User Takes User payload and returns token and referesh token along with saved user entity
func (a *UserAPI) Register(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error

		var userPayload dto.CreateUserDTO
		// Decode json payload from client and store the values to userPayload
		if err = c.ShouldBind(&userPayload); err != nil {
			api.SendErrorResponse(a.config.Localize, c, "invalid_input", http.StatusUnprocessableEntity, nil)
			return
		}

		defer c.Request.Body.Close()

		// validate input
		if errs := validators.ValidateSignUpData(&userPayload, a.config.Localize, c); errs != nil {
			api.SendErrorResponse(a.config.Localize, c, "invalid_input", http.StatusUnprocessableEntity, errs)
			return
		}

		response, err := a.userRepo.Create(&userPayload, role)

		if err != nil {
			api.SendErrorResponse(a.config.Localize, c, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		token, err := tokens.GenerateNoExpiryTokenWithCustomType(response.ID.Hex(), "verify", response.UserRole)
		if err != nil {
			api.SendErrorResponse(a.config.Localize, c, "something_went_wrong", http.StatusInternalServerError, nil)
			return
		}

		err = a.config.Mailer.SendNoReplyMail([]string{response.Email},
			a.config.Localize.GetMessage("confirm_email_address", c),
			"register",
			models.WelcomeEmailAndVerifyEmailTemplate(
				a.config.Localize, c,
				a.config.FrontendBaseUrl+"/user/confirm-email?token="+token))

		if err != nil {
			api.SendErrorResponse(a.config.Localize, c, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		c.JSON(http.StatusCreated, &serialize.RegisterResponse{
			Response: serialize.Response{
				StatusCode: http.StatusCreated,
				Message:    a.config.Localize.GetMessage("user_registered", c),
			},
			User: response,
		})
	}
}

// Confirm email address from mail using JWT.
func (a *UserAPI) ConfirmEmailAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get the token from request
		token := c.Request.URL.Query().Get("token")

		// parse refresh token
		claimedToken, err := tokens.DecodeJSONWebToken(token)
		if err != nil {
			api.SendErrorResponse(a.config.Localize, c, err.Error(), http.StatusNotAcceptable, nil)
			return
		}

		if claimedToken.Type != "verify" {
			api.SendErrorResponse(a.config.Localize, c, "token_malformed", http.StatusNotAcceptable, nil)
			return

		}

		userId, err := primitive.ObjectIDFromHex(claimedToken.ID)
		if err != nil {
			api.SendErrorResponse(a.config.Localize, c, err.Error(), http.StatusNotAcceptable, nil)
			return
		}

		// cross check refresh token with db.
		user, err := a.userRepo.FindUserByID(&userId)
		if err != nil {
			api.SendErrorResponse(a.config.Localize, c, err.Error(), http.StatusUnauthorized, nil)
			return
		}

		if user.EmailVerified {
			api.SendErrorResponse(a.config.Localize, c, "email_is_already_verified", http.StatusUnauthorized, nil)
			return
		}

		// generate new token and refresh token
		token, refreshToken, err := a.generateAndSetTokens(user.ID, c, user.UserRole)
		if err != nil {
			api.SendErrorResponse(a.config.Localize, c, "something_went_wrong", http.StatusInternalServerError, nil)
			return
		}

		err = a.userRepo.ActivateUserAccount(userId, refreshToken)
		if err != nil {
			api.SendErrorResponse(a.config.Localize, c, "something_went_wrong", http.StatusInternalServerError, nil)
			return
		}

		c.JSON(http.StatusOK, &serialize.AuthResponse{
			Response: serialize.Response{
				StatusCode: http.StatusOK,
				Message:    a.config.Localize.GetMessage("success_login", c),
			},
			Token:        token,
			RefreshToken: refreshToken,
			User:         user,
		})
	}
}

// Login User Checks User with Email and Password and returns token and referesh token along with user entity
func (a *UserAPI) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error

		var loginPayload dto.LoginDTO

		if err = c.ShouldBind(&loginPayload); err != nil {
			api.SendErrorResponse(a.config.Localize, c, "invalid_input", http.StatusUnprocessableEntity, nil)
			return
		}

		defer c.Request.Body.Close()

		// Validate input
		if errors := validators.ValidateLoginInput(&loginPayload, a.config.Localize, c); errors != nil {
			api.SendErrorResponse(a.config.Localize, c, "login_cred_req", http.StatusUnprocessableEntity, errors)
			return
		}

		// check user exists with email
		user, err := a.userRepo.FindUserByEmail(loginPayload.Email)

		if err != nil {
			api.SendErrorResponse(a.config.Localize, c, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		// check for password match
		err = user.CompareHashedPassword(loginPayload.Password)
		if err != nil {
			api.SendErrorResponse(a.config.Localize, c, "invalid_password_entered", http.StatusUnauthorized, nil)
			return
		}

		if !user.EmailVerified {
			api.SendErrorResponse(a.config.Localize, c, "email_not_verified", http.StatusUnauthorized, nil)
			return
		}

		// generate new token and refresh token
		token, refreshToken, err := a.generateAndSetTokens(user.ID, c, user.UserRole)
		if err != nil {
			api.SendErrorResponse(a.config.Localize, c, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		c.JSON(http.StatusOK, &serialize.AuthResponse{
			Response: serialize.Response{
				StatusCode: http.StatusOK,
				Message:    a.config.Localize.GetMessage("success_login", c),
			},
			User:         user,
			Token:        token,
			RefreshToken: refreshToken,
		})
	}
}

// Refresh access token using refresh token.
// Can be force refreshed by passing force param set to true.
func (a *UserAPI) RefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, refreshToken, err := a.getAccessAndRefreshTokenFromRequest(c)
		if err != nil {
			api.SendErrorResponse(a.config.Localize, c, err.Error(), http.StatusUnauthorized, nil)
			return
		}

		// check if force refresh is requested
		isForce, _ := strconv.ParseBool(c.Request.URL.Query().Get("force"))

		// parse auth Token
		_, err = tokens.DecodeJSONWebToken(token)
		if err == nil && !isForce {
			api.SendErrorResponse(a.config.Localize, c, "token_is_not_expired", http.StatusNotAcceptable, nil)
			return
		}

		// parse refresh token
		claimedRefreshToken, err := tokens.DecodeJSONWebToken(refreshToken)
		if err != nil {
			api.SendErrorResponse(a.config.Localize, c, "revalidate_token_malformed", http.StatusNotAcceptable, nil)
			return
		}

		userId, err := primitive.ObjectIDFromHex(claimedRefreshToken.ID)
		if err != nil {
			api.SendErrorResponse(a.config.Localize, c, "something_went_wrong", http.StatusNotAcceptable, nil)
			return
		}

		// cross check refresh token with db.
		user, err := a.userRepo.FindUserByID(&userId)
		if err != nil {
			api.SendErrorResponse(a.config.Localize, c, err.Error(), http.StatusNotFound, nil)
			return
		}

		if user.RefreshToken != refreshToken {
			api.SendErrorResponse(a.config.Localize, c, "revalidate_token_malformed", http.StatusInternalServerError, nil)
			return
		}

		// generate new token and refresh token
		token, refreshToken, err = a.generateAndSetTokens(user.ID, c, user.UserRole)
		if err != nil {
			api.SendErrorResponse(a.config.Localize, c, "something_went_wrong", http.StatusInternalServerError, nil)
			return
		}

		c.JSON(http.StatusOK, &serialize.RefreshResponse{
			Response: serialize.Response{
				StatusCode: http.StatusOK,
				Message:    a.config.Localize.GetMessage("token_refreshed_successfully", c),
			},
			Token:        token,
			RefreshToken: refreshToken,
		})
	}
}

// Logout - Removes user token from Entity and disables token from used further.
func (a *UserAPI) Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := api.GetUserIdFromContext(c)
		if err != nil {
			api.SendErrorResponse(a.config.Localize, c, err.Error(), http.StatusUnauthorized, nil)
			return
		}

		err = a.userRepo.UpdateByField(*id, "refresh_token", "")

		if err != nil {
			api.SendErrorResponse(a.config.Localize, c, err.Error(), http.StatusUnauthorized, nil)
			return
		}

		c.SetCookie(string(constants.AUTH_TOKEN), "", 0, "/", a.config.Domain, false, true)
		c.SetCookie(string(constants.REFRESH_TOKEN), "", 0, "/", a.config.Domain, false, true)

		c.JSON(http.StatusOK, serialize.Response{
			StatusCode: http.StatusOK,
			Message:    a.config.Localize.GetMessage("logged_out_successfully", c),
		})

	}
}

// Send reset email to user's email.
func (a *UserAPI) ForgotPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get email from body.
		var email struct {
			Email string `json:"email"`
		}

		if err := c.ShouldBind(&email); err != nil {
			api.SendErrorResponse(a.config.Localize, c, "invalid_input", http.StatusUnprocessableEntity, nil)
			return
		}
		defer c.Request.Body.Close()

		if email.Email == "" {
			api.SendErrorResponse(a.config.Localize, c, "invalid_input", http.StatusUnprocessableEntity, nil)
			return
		}

		// find user by email.
		user, err := a.userRepo.FindUserByEmail(email.Email)

		if err != nil {
			api.SendErrorResponse(a.config.Localize, c, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		token, err := tokens.GenerateTokenWithExpiryTimeAndType(user.ID.Hex(), int64(time.Minute*5), "reset", user.UserRole)

		if err != nil {
			api.SendErrorResponse(a.config.Localize, c, "something_went_wrong", http.StatusInternalServerError, nil)
			return
		}

		err = a.userRepo.UpdateByField(user.ID, "reset_password_token", token)

		if err != nil {
			api.SendErrorResponse(a.config.Localize, c, "something_went_wrong", http.StatusInternalServerError, nil)
			return
		}

		a.config.Mailer.SendNoReplyMail([]string{email.Email},
			a.config.Localize.GetMessage("request_for_change_password", c),
			"forgot_password",
			models.ForgotPasswordEmailTemplateData(a.config.Localize, c,
				a.config.FrontendBaseUrl+"/user/reset-password?token="+token))

		c.JSON(200, map[string]interface{}{
			"status_code": http.StatusOK,
			"message":     a.config.Localize.GetMessage("reset_email_sent_successfully", c),
		})

	}
}

// reset password from email.
func (a *UserAPI) ResetPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get token from query
		token := c.Request.URL.Query().Get("token")

		// parse refresh token
		claimedToken, err := tokens.DecodeJSONWebToken(token)
		if err != nil {
			api.SendErrorResponse(a.config.Localize, c, err.Error(), http.StatusNotAcceptable, nil)
			return
		}

		if claimedToken.Type != "reset" {
			api.SendErrorResponse(a.config.Localize, c, "token_malformed", http.StatusNotAcceptable, nil)
			return

		}
		// get new password from req body
		var password struct {
			Password string `json:"password"`
		}

		if err := c.ShouldBind(&password); err != nil || password.Password == "" {
			api.SendErrorResponse(a.config.Localize, c, "invalid_input", http.StatusUnprocessableEntity, nil)
			return
		}
		defer c.Request.Body.Close()

		userId, err := primitive.ObjectIDFromHex(claimedToken.ID)
		if err != nil {
			api.SendErrorResponse(a.config.Localize, c, err.Error(), http.StatusNotAcceptable, nil)
			return
		}

		// check user exists with id
		user, err := a.userRepo.FindUserByID(&userId)
		if err != nil {
			api.SendErrorResponse(a.config.Localize, c, err.Error(), http.StatusNotFound, nil)
			return
		}

		// // check for token matches
		if user.ResetPasswordToken != token {
			api.SendErrorResponse(a.config.Localize, c, "token_malformed", http.StatusNotAcceptable, nil)
			return
		}

		user.Password = password.Password
		user.HashPassword()

		err = a.userRepo.UpdateByField(userId, "password", user.Password)
		if err != nil {
			api.SendErrorResponse(a.config.Localize, c, err.Error(), http.StatusInternalServerError, nil)
			return
		}

		err = a.userRepo.UpdateByField(user.ID, "reset_password_token", "")

		if err != nil {
			api.SendErrorResponse(a.config.Localize, c, "something_went_wrong", http.StatusInternalServerError, nil)
			return
		}

		// generate new token and refresh token
		newToken, refreshToken, err := a.generateAndSetTokens(user.ID, c, user.UserRole)
		if err != nil {
			api.SendErrorResponse(a.config.Localize, c, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		c.JSON(http.StatusOK, &serialize.AuthResponse{
			Response: serialize.Response{
				StatusCode: http.StatusOK,
				Message:    a.config.Localize.GetMessage("password_reset_success_and_logged_in", c),
			},
			User:         user,
			Token:        newToken,
			RefreshToken: refreshToken,
		})

	}
}
