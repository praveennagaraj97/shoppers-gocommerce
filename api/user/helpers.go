package userapi

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/praveennagaraj97/shoppers-gocommerce/api"
	"github.com/praveennagaraj97/shoppers-gocommerce/constants"
	"github.com/praveennagaraj97/shoppers-gocommerce/pkg/tokens"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// generateAccessAndRefreshToken  returs access, refresh, error
func (a *UserAPI) generateAccessAndRefreshToken(userId primitive.ObjectID, c *gin.Context, role string) (string, string, error) {
	var er error = nil

	token, err := tokens.GenerateTokenWithExpiryTimeAndType(userId.Hex(),
		time.Now().Local().Add(time.Minute*constants.JWT_AccessTokenExpiry).Unix(), "access", role)
	if err != nil {
		er = err
	}
	refreshToken, err := tokens.GenerateNoExpiryTokenWithCustomType(userId.Hex(), "refresh", role)
	if err != nil {
		er = err
	}

	return token, refreshToken, er
}

// returns access, refresh and error
func (a *UserAPI) getAccessAndRefreshTokenFromRequest(c *gin.Context) (string, string, error) {
	var token string
	var refresh_token string

	// Get auth token either from cookie nor Header
	cookie, err := c.Request.Cookie(string(constants.AUTH_TOKEN))
	if err != nil {
		// check in auth header as bearer
		authHeader := c.Request.Header.Get("Authorization")
		containsBearerToken := strings.HasPrefix(authHeader, "Bearer")
		if !containsBearerToken {
			token = ""
		} else {
			token = strings.Split(authHeader, "Bearer ")[1]
		}
	} else {
		token = cookie.Value
	}

	// Get refresh token either from cookie nor params
	refreshCookie, err := c.Request.Cookie(string(constants.REFRESH_TOKEN))
	if err != nil {
		// check in auth header as bearer
		refresh_token = c.Request.URL.Query().Get("refresh_token")
		if refresh_token == "" {
			return "", "", errors.New("refresh token is missing")
		}

	} else {
		refresh_token = refreshCookie.Value
	}

	return token, refresh_token, nil
}

func (a *UserAPI) generateAndSetTokens(userId primitive.ObjectID, c *gin.Context, role string) (string, string, error) {
	token, refreshToken, err := a.generateAccessAndRefreshToken(userId, c, role)
	if err != nil {
		api.SendErrorResponse(a.config.Localize, c, err.Error(), http.StatusUnauthorized, nil)
		return "", "", err
	}

	err = a.userRepo.UpdateByField(userId, "refresh_token", refreshToken)
	if err != nil {
		api.SendErrorResponse(a.config.Localize, c, "something_went_wrong", http.StatusInternalServerError, nil)
		return "", "", err
	}

	c.SetCookie(string(constants.AUTH_TOKEN), token, constants.CookieAccessExpiryTime, "/", a.config.Domain, false, true)
	c.SetCookie(string(constants.REFRESH_TOKEN), refreshToken, constants.CookieRefreshExpiryTime, "/", a.config.Domain, false, true)

	return token, refreshToken, nil
}
