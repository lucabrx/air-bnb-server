package main

import (
	"errors"
	"fmt"
	"github.com/air-bnb/internal/data"
	"github.com/air-bnb/internal/validator"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
	"time"
)

func (app *application) createBookingHandler(w http.ResponseWriter, r *http.Request) {
	session := app.contextGetUser(r)

	var input struct {
		ListingID int64     `json:"listingId"`
		StartDate time.Time `json:"startDate"`
		EndDate   time.Time `json:"endDate"`
		Pricing   int64     `json:"pricing"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	total := input.Pricing * int64(input.EndDate.Sub(input.StartDate).Hours()/24)

	booking := &data.Booking{
		ListingID: input.ListingID,
		GuestID:   session.ID,
		CheckIn:   input.StartDate,
		CheckOut:  input.EndDate,
		Price:     input.Pricing,
		Total:     total,
	}

	v := validator.New()
	data.ValidateBooking(v, booking)
	fmt.Println(v.Errors)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Bookings.Insert(booking)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"booking": booking}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) getBookingHandler(w http.ResponseWriter, r *http.Request) {
	params := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(params, 10, 64)
	if err != nil || id < 1 {
		app.notFoundResponse(w, r)
		return
	}

	booking, err := app.models.Bookings.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"booking": booking}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteBookingHandler(w http.ResponseWriter, r *http.Request) {
	session := app.contextGetUser(r)
	params := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(params, 10, 64)
	if err != nil || id < 1 {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Bookings.Delete(id, session.ID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "booking successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getUserBookingsHandler(w http.ResponseWriter, r *http.Request) {
	session := app.contextGetUser(r)

	bookings, err := app.models.Bookings.GetForUser(session.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	for _, booking := range bookings {
		image, err := app.models.Images.GetForListing(booking.ListingID)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		booking.Listing.Images = image
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"bookings": bookings}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getPropertyBookingsHandler(w http.ResponseWriter, r *http.Request) {
	params := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(params, 10, 64)
	if err != nil || id < 1 {
		app.notFoundResponse(w, r)
		return
	}

	bookings, err := app.models.Bookings.GetForListing(id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	for _, booking := range bookings {
		image, err := app.models.Images.GetForListing(booking.ListingID)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		booking.Listing.Images = image
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"bookings": bookings}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
