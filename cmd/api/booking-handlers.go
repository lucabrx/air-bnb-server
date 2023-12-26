package main

import "net/http"

// create booking
func (app *application) createBookingHandler(w http.ResponseWriter, r *http.Request) {
	session := app.contextGetUser(r)

	var input struct {
		PropertyID int    `json:"propertyId"`
		UserID     int    `json:"userId"`
		StartDate  string `json:"startDate"`
		EndDate    string `json:"endDate"`
		Pricing    int    `json:"pricing"`
		Total      int    `json:"total"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

}

// get booking
// delete booking
// get user bookings
// get property bookings
