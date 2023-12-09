package main

import (
	"errors"
	"github.com/air-bnb/internal/data"
	"github.com/air-bnb/internal/random"
	"github.com/air-bnb/internal/validator"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (app *application) registerUserEmailHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Name     string `json:"name"`
		Password string `json:"password"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &data.User{
		Email: input.Email,
		Name:  input.Name,
	}
	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	user.VerificationToken = random.RandString(3) + "-" + random.RandString(3)

	v := validator.New()
	data.ValidateUser(v, user)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	emailData := struct {
		Name              string
		VerificationToken string
	}{
		Name:              user.Name,
		VerificationToken: user.VerificationToken,
	}

	err = app.sendEmail(
		"./templates/email-code.tmpl",
		emailData,
		user.Email,
		"Air BnB Clone - Email Verification",
	)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) verificationUserHandler(w http.ResponseWriter, r *http.Request) {
	params := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(params, 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	var input struct {
		Code string `json:"code"`
	}
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	user, err := app.models.Users.Get(id, "")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	v := validator.New()
	if user.Activated {
		v.AddError("code", "user already activated")
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	if user.VerificationToken != input.Code {
		v.AddError("code", "invalid verification code")
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user.Activated = true
	user.VerificationToken = ""
	err = app.models.Users.Update(user)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "ok"}, nil)
}
