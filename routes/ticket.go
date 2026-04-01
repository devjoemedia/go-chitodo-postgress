package routes

import (
	"github.com/devjoemedia/go-ticketing-api/handlers"
	"github.com/go-chi/chi/v5"
)

func TicketRoute() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", handlers.GetTickets)
	// r.Get("/{id}", handlers.GetTicketByID)
	r.Post("/", handlers.CreateTicket)
	// r.Put("/assign", handlers.AssignTicket)
	// r.Patch("/{id}", handlers.UpdateTicket)
	// r.Delete("/{id}", handlers.DeleteTicket)

	return r
}
