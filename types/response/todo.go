package api_response

import (
	"github.com/devjoemedia/scrumpilot-go-api/models"
	"github.com/devjoemedia/scrumpilot-go-api/types"
)

type GetTodosResponse struct {
	Success    bool             `json:"success"`
	Status     int              `json:"status"`
	Message    string           `json:"message"`
	Todos      []models.Todo    `json:"todos"`
	Pagination types.Pagination `json:"pagination"`
}

type CreateTodoResponse struct {
	Success bool         `json:"success"`
	Status  int          `json:"status"`
	Message string       `json:"message"`
	Todo    *models.Todo `json:"todo"`
}

type GetTodoResponse struct {
	Success bool         `json:"success"`
	Status  int          `json:"status"`
	Message string       `json:"message"`
	Todo    *models.Todo `json:"todo"`
}

type UpdateTodoResponse struct {
	Success bool         `json:"success"`
	Status  int          `json:"status"`
	Message string       `json:"message"`
	Todo    *models.Todo `json:"todo"`
}

type DeleteTodoResponse struct {
	Success bool   `json:"success"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}
