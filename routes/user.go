package routes

import (
	"github.com/devjoemedia/scrumpilot-go-api/handlers"
	"github.com/go-chi/chi/v5"
)

func UserRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/", handlers.GetUsers)
	r.Get("/{id}", handlers.GetUser)
	r.Patch("/{id}", handlers.UpdateUser)
	r.Delete("/{id}", handlers.DeleteUser)
	return r
}
