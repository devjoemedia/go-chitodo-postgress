package api_response

import (
	"github.com/devjoemedia/scrumpilot-go-api/models"
	"github.com/devjoemedia/scrumpilot-go-api/types"
)

type GetTicketsResponse struct {
	Success    bool             `json:"success"`
	Status     int              `json:"status"`
	Message    string           `json:"message"`
	Tickets    []models.Ticket  `json:"tickets"`
	Pagination types.Pagination `json:"pagination"`
}

type GetTicketResponse struct {
	Success bool           `json:"success"`
	Status  int            `json:"status"`
	Message string         `json:"message"`
	Ticket  *models.Ticket `json:"ticket"`
}

type CreateTicketResponse struct {
	Success bool           `json:"success"`
	Status  int            `json:"status"`
	Message string         `json:"message"`
	Ticket  *models.Ticket `json:"ticket"`
}

type UpdateTicketResponse struct {
	Success bool           `json:"success"`
	Status  int            `json:"status"`
	Message string         `json:"message"`
	Ticket  *models.Ticket `json:"ticket"`
}

type DeleteTicketResponse struct {
	Success bool   `json:"success"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}
