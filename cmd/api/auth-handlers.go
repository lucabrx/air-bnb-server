package main

import (
	"context"
	"encoding/json"
	"errors"
	"golang.org/x/oauth2"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/air-bnb/internal/data"
	"github.com/air-bnb/internal/random"
	"github.com/air-bnb/internal/validator"
	"github.com/go-chi/chi/v5"
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

	err = app.writeJSON(w, http.StatusCreated, envelope{"userId": user.ID}, nil)
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

	token, err := app.models.Tokens.New(user.ID, 30*24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	cookie := app.sessionCookie(token.Plaintext, token.Expiry)
	http.SetCookie(w, cookie)

	err = app.writeJSON(w, http.StatusOK, envelope{"cookie value": cookie.Value}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
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

	if !user.Activated {
		app.inactiveAccountResponse(w, r)
		return
	}
	match, err := user.Password.Matches(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if !match {
		app.invalidAuthenticationTokenResponse(w, r)
		return
	}
	token, err := app.models.Tokens.New(user.ID, 30*24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	cookie := app.sessionCookie(token.Plaintext, token.Expiry)
	http.SetCookie(w, cookie)
	err = app.writeJSON(w, http.StatusOK, envelope{"cookie value": cookie.Value}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) logoutHandler(w http.ResponseWriter, r *http.Request) {
	session := app.contextGetUser(r)

	err := app.models.Tokens.DeleteAllForUser(data.ScopeAuthentication, session.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	cookie := app.sessionCookie("", time.Unix(0, 0))
	http.SetCookie(w, cookie)

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "you have been logged out"}, nil)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}
}

func (app *application) githubLoginHandler(w http.ResponseWriter, r *http.Request) {
	session := app.contextGetUser(r)

	if !session.IsAnonymous() {
		app.alreadyHaveSessionResponse(w, r)
		return
	}

	url := app.githubConfig().AuthCodeURL("state", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (app *application) githubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	session := app.contextGetUser(r)

	if !session.IsAnonymous() {
		app.alreadyHaveSessionResponse(w, r)
		return
	}

	code := r.URL.Query().Get("code")
	token, err := app.githubConfig().Exchange(r.Context(), code)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	client := app.githubConfig().Client(context.Background(), token)

	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var userDetails struct {
		Name      string `json:"name"`
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := json.Unmarshal(body, &userDetails); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if userDetails.Email == "" {
		respEmails, err := client.Get("https://api.github.com/user/emails")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer respEmails.Body.Close()

		var emails []struct {
			Email   string `json:"email"`
			Primary bool   `json:"primary"`
		}

		bodyEmails, err := io.ReadAll(respEmails.Body)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		if err := json.Unmarshal(bodyEmails, &emails); err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		for _, e := range emails {
			if e.Primary {
				userDetails.Email = e.Email
				break
			}
		}
	}

	user := &data.User{
		Name:      userDetails.Name,
		Email:     userDetails.Email,
		Image:     userDetails.AvatarURL,
		Activated: true,
	}

	dbUser, err := app.models.Users.Get(0, user.Email)
	if dbUser == nil {
		err = app.models.Users.Insert(user)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	} else {
		user.ID = dbUser.ID
	}

	sessionToken, err := app.models.Tokens.New(user.ID, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	cookie := app.sessionCookie(sessionToken.Plaintext, token.Expiry)
	http.SetCookie(w, cookie)

	http.Redirect(w, r, "http://localhost:3000", http.StatusTemporaryRedirect)
}

func (app *application) googleLoginHandler(w http.ResponseWriter, r *http.Request) {
	session := app.contextGetUser(r)

	if !session.IsAnonymous() {
		app.alreadyHaveSessionResponse(w, r)
		return
	}

	url := app.googleConfig().AuthCodeURL("state", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (app *application) googleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	session := app.contextGetUser(r)

	if !session.IsAnonymous() {
		app.alreadyHaveSessionResponse(w, r)
		return
	}

	code := r.URL.Query().Get("code")
	token, err := app.googleConfig().Exchange(r.Context(), code)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	client := app.googleConfig().Client(context.Background(), token)

	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo?fields=email,name,picture")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var userDetails struct {
		Name      string `json:"name"`
		Email     string `json:"email"`
		AvatarURL string `json:"picture"`
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := json.Unmarshal(body, &userDetails); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := &data.User{
		Name:      userDetails.Name,
		Email:     userDetails.Email,
		Image:     userDetails.AvatarURL,
		Activated: true,
	}

	dbUser, err := app.models.Users.Get(0, user.Email)
	if dbUser == nil {
		err = app.models.Users.Insert(user)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	} else {
		user.ID = dbUser.ID
	}

	sessionToken, err := app.models.Tokens.New(user.ID, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	cookie := app.sessionCookie(sessionToken.Plaintext, token.Expiry)
	http.SetCookie(w, cookie)

	http.Redirect(w, r, "http://localhost:3000", http.StatusTemporaryRedirect)
}
