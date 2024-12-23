package server

import (
	"net/http"

	"github.com/anunayjoshi29/token-server/internal/handlers"
	"github.com/anunayjoshi29/token-server/internal/routecalc"
	"github.com/go-chi/chi/v5"
)

func NewServer(finder *routecalc.Finder, cache *routecalc.Cache) *http.Server {
	r := chi.NewRouter()

	r.Post("/routes", handlers.RoutesHandler(finder, cache))

	return &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
}
