package serialize

import "github.com/praveennagaraj97/shoppers-gocommerce/models"

type Response struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

type ErrorResponse struct {
	Errors interface{} `json:"errors,omitempty"`
	Response
}

type RefreshResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	Response
}

type AuthResponse struct {
	User         interface{} `json:"user"`
	Token        string      `json:"token"`
	RefreshToken string      `json:"refresh_token"`
	Response
}

type RegisterResponse struct {
	User *models.User `json:"user"`
	Response
}

type DataResponse struct {
	Data interface{} `json:"data"`
	Response
}

type PaginatedDataResponse struct {
	Count            *uint64 `json:"count"`
	Next             *bool   `json:"next"`
	Prev             *bool   `json:"prev"`
	PaginateKeySetID *string `json:"paginate_id"`
	DataResponse
}
