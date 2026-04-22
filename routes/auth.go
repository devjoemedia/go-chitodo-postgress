package routes

import (
	"github.com/devjoemedia/scrumpilot-go-api/handlers"
	"github.com/go-chi/chi/v5"
)

func AuthRoute() *chi.Mux {
	r := chi.NewRouter()

	r.Post("/register", handlers.Register)
	r.Post("/login", handlers.Login)
	r.Post("/refresh", handlers.RefreshToken)
	r.Post("/logout", handlers.Logout)

	return r
}
