package main

import (
	"errors"
	"github.com/air-bnb/internal/data"
	"github.com/air-bnb/internal/random"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	session := app.contextGetUser(r)

	err := app.writeJSON(w, http.StatusOK, envelope{"user": session}, nil)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}
}

func (app *application) resetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email string `json:"email"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user, err := app.models.Users.Get(0, input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	user.ResetToken = random.RandString(3) + "-" + random.RandString(3)

	err = app.models.Users.Update(user)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	emailData := struct {
		Name       string
		ResetToken string
	}{
		Name:       user.Name,
		ResetToken: user.ResetToken,
	}

	err = app.sendEmail(
		"./templates/reset-password-code.tmpl",
		emailData,
		user.Email,
		"Air BnB Clone - Reset Password",
	)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"email": user.Email}, nil)
}

func (app *application) resetPasswordConfirmHandler(w http.ResponseWriter, r *http.Request) {
	email := chi.URLParam(r, "email")

	var input struct {
		ResetToken  string `json:"resetToken"`
		NewPassword string `json:"newPassword"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user, err := app.models.Users.Get(0, email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if user.ResetToken != input.ResetToken {
		app.invalidAuthenticationTokenResponse(w, r)
		return
	}

	err = user.Password.Set(input.NewPassword)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	user.ResetToken = ""

	err = app.models.Users.Update(user)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "success"}, nil)
}
