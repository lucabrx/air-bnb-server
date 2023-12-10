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
	r.Use(app.authenticate)
	r.Use(app.enableCORS)

	r.NotFound(app.notFoundResponse)
	r.MethodNotAllowed(app.methodNotAllowedResponse)

	r.Route("/v1/auth", func(r chi.Router) {
		r.Post("/register", app.registerUserEmailHandler)
		r.Post("/verify/{id}", app.verificationUserHandler)
		r.Post("/login", app.loginUserHandler)
		r.Delete("/logout", app.requireActivatedUser(app.logoutHandler))
	})

	r.Route("/v1/user", func(r chi.Router) {
		r.Get("/", app.requireActivatedUser(app.getUserHandler))
		r.Post("/reset-password", app.resetPasswordHandler)
		r.Post("/new-password/{email}", app.resetPasswordConfirmHandler)
		r.Delete("/", app.requireActivatedUser(app.deleteUserHandler))
		r.Patch("/", app.requireActivatedUser(app.updateUserHandler))
		r.Patch("/password", app.requireActivatedUser(app.updatePasswordHandler))
		r.Post("/change-email", app.requireActivatedUser(app.changeEmailHandler))
		r.Post("/change-email/verify/{email}", app.verifyChangeEmailHandler)
	})

	r.Route("/v1/upload", func(r chi.Router) {
		r.Post("/image", app.requireAuthenticatedUser(app.uploadImageHandler))
	})

	r.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		res := map[string]string{
			"status": "ok",
		}

		app.writeJSON(w, http.StatusOK, res, nil)
	})

	return r
}
