package main

import (
	"errors"
	"github.com/air-bnb/internal/data"
	"github.com/air-bnb/internal/random"
	"github.com/air-bnb/internal/validator"
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

func (app *application) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)

	err := app.models.Users.Delete(user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "success"}, nil)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}
}

func (app *application) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)

	var input struct {
		Name  string `json:"name"`
		Image string `json:"image"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	v.Check(input.Name != "", "name", "must be provided")
	v.Check(len(input.Name) >= 3, "name", "must be at least 3 bytes long")
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user.Name = input.Name

	if input.Image == "" {
		user.Image = ""
	} else {
		user.Image = input.Image
	}

	err = app.models.Users.Update(user)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}
}

func (app *application) updatePasswordHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)

	var input struct {
		OldPassword string `json:"oldPassword"`
		NewPassword string `json:"newPassword"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	match, err := user.Password.Matches(input.OldPassword)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if !match {
		app.invalidAuthenticationTokenResponse(w, r)
		return
	}

	err = user.Password.Set(input.NewPassword)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.models.Users.Update(user)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "success"}, nil)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}
}

func (app *application) requestChangeEmailHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)

	user.ResetEmailToken = random.RandString(3) + "-" + random.RandString(3)

	err := app.models.Users.Update(user)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	emailData := struct {
		Name       string
		ResetToken string
	}{
		Name:       user.Name,
		ResetToken: user.ResetEmailToken,
	}

	err = app.sendEmail(
		"./templates/change-email-code.tmpl",
		emailData,
		user.Email,
		"Air BnB Clone - Change Email",
	)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "ok"}, nil)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}
}

func (app *application) changeEmailHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)

	var input struct {
		ResetToken string `json:"resetToken"`
		NewEmail   string `json:"newEmail"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if user.ResetEmailToken != input.ResetToken {
		app.invalidAuthenticationTokenResponse(w, r)
		return
	}

	user.Email = input.NewEmail
	user.ResetEmailToken = ""

	err = app.models.Users.Update(user)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "success"}, nil)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}
}
