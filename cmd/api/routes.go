package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (app *application) routes() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		res := map[string]string{
			"status": "ok",
		}

		app.writeJSON(w, http.StatusOK, res, nil)
	})

	return r
}
