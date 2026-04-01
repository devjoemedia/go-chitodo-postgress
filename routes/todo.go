package routes

import (
	"github.com/devjoemedia/go-ticketing-api/handlers"
	"github.com/go-chi/chi/v5"
)

func TodoRoute() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", handlers.GetTodos)
	r.Get("/{id}", handlers.GetTodoByID)
	r.Post("/", handlers.CreateTodo)
	r.Patch("/{id}", handlers.UpdateTodo)
	r.Delete("/{id}", handlers.DeleteTodo)

	return r
}
