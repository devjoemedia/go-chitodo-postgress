package api_response

import (
	"github.com/devjoemedia/go-ticketing-api/models"
	"github.com/devjoemedia/go-ticketing-api/types"
)

type GetUsersResponse struct {
	Success    bool             `json:"success"`
	Status     int              `json:"status"`
	Message    string           `json:"message"`
	Users      []models.User    `json:"users"`
	Pagination types.Pagination `json:"pagination"`
}

type GetUserResponse struct {
	Success bool         `json:"success"`
	Status  int          `json:"status"`
	Message string       `json:"message"`
	User    *models.User `json:"user"`
}

type UpdateUserResponse struct {
	Success    bool             `json:"success"`
	Status     int              `json:"status"`
	Message    string           `json:"message"`
	User       *models.User     `json:"user"`
	Pagination types.Pagination `json:"pagination"`
}

type DeleteUserResponse struct {
	Success bool   `json:"success"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}
