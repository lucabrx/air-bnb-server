package main

import "net/http"

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	session := app.contextGetUser(r)

	err := app.writeJSON(w, http.StatusOK, envelope{"user": session}, nil)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}
}
