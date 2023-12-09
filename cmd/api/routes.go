package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func (app *application) routes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.NotFound(app.notFoundResponse)
	r.MethodNotAllowed(app.methodNotAllowedResponse)

	r.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		res := map[string]string{
			"status": "ok",
		}

		app.writeJSON(w, http.StatusOK, res, nil)
	})

	return r
}
