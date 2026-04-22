package api_response

import "github.com/devjoemedia/scrumpilot-go-api/models"

type AuthTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"` // always "Bearer"
	ExpiresIn    int    `json:"expires_in"` // seconds
}

type RegisterResponse struct {
	Success bool         `json:"success"`
	Status  int          `json:"status"`
	Message string       `json:"message"`
	User    *models.User `json:"user"`
	Tokens  AuthTokens   `json:"tokens"`
}

type LoginResponse struct {
	Success bool         `json:"success"`
	Status  int          `json:"status"`
	Message string       `json:"message"`
	User    *models.User `json:"user"`
	Tokens  AuthTokens   `json:"tokens"`
}

type RefreshResponse struct {
	Success bool       `json:"success"`
	Status  int        `json:"status"`
	Message string     `json:"message"`
	Tokens  AuthTokens `json:"tokens"`
}

type LogoutResponse struct {
	Success bool   `json:"success"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}
