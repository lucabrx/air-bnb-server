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
		r.Get("/github/login", app.githubLoginHandler)
		r.Get("/github/callback", app.githubCallbackHandler)
		r.Get("/google/login", app.googleLoginHandler)
		r.Get("/google/callback", app.googleCallbackHandler)
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

	r.Route("/v1/listings", func(r chi.Router) {
		r.Get("/user-listings", app.requireActivatedUser(app.getAllUserListingsHandler))
		r.Get("/{listingId}", app.getListingHandler)
		r.Get("/", app.getAllListingsHandler)
		r.Post("/", app.requireActivatedUser(app.createListingHandler))
		r.Patch("/{listingId}", app.requireActivatedUser(app.updateListingHandler))
		r.Delete("/delete/{listingId}", app.requireActivatedUser(app.deleteListingHandler))
		r.Post("/{listingId}/images", app.requireActivatedUser(app.addImageToListingGalleryHandler))
		r.Delete("/images/{imageId}", app.requireActivatedUser(app.removeImageFromListingGalleryHandler))
		r.Post("/images/{listingId}", app.requireActivatedUser(app.uploadImagesToListingHandler))
	})

	r.Route("/v1/bookings", func(r chi.Router) {
		r.Post("/", app.requireActivatedUser(app.createBookingHandler))
		r.Get("/{id}", app.requireActivatedUser(app.getBookingHandler))
		r.Delete("/{id}", app.requireActivatedUser(app.deleteBookingHandler))
		r.Get("/user-bookings", app.requireActivatedUser(app.getUserBookingsHandler))
		r.Get("/property-bookings/{id}", app.requireActivatedUser(app.getPropertyBookingsHandler))
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
