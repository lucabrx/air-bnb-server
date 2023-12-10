package main

import (
	"errors"
	"github.com/air-bnb/internal/data"
	"github.com/air-bnb/internal/validator"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (app *application) createListingHandler(w http.ResponseWriter, r *http.Request) {
	session := app.contextGetUser(r)
	var input struct {
		Title         string   `json:"title"`
		Description   string   `json:"description"`
		Category      string   `json:"category"`
		RoomCount     int64    `json:"room_count"`
		BathroomCount int64    `json:"bathroom_count"`
		GuestCount    int64    `json:"guest_count"`
		Location      string   `json:"location"`
		Price         int64    `json:"price"`
		Images        []string `json:"images"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	listing := &data.Listing{
		Title:         input.Title,
		Description:   input.Description,
		Category:      input.Category,
		RoomCount:     input.RoomCount,
		BathroomCount: input.BathroomCount,
		GuestCount:    input.GuestCount,
		Location:      input.Location,
		Price:         input.Price,
		OwnerID:       session.ID,
		OwnerName:     session.Name,
		OwnerPhoto:    session.Image,
	}
	v := validator.New()
	if data.ValidateListing(v, listing); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.models.Listings.Insert(listing)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	for _, inputImage := range input.Images {
		dbImage := &data.Image{
			ListingID: listing.ID,
			Url:       inputImage,
		}
		err = app.models.Images.Insert(dbImage)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	}
	// return a response
	err = app.writeJSON(w, http.StatusCreated, envelope{"listingId": listing.ID}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getListingHandler(w http.ResponseWriter, r *http.Request) {
	params := chi.URLParam(r, "listingId")
	id, err := strconv.ParseInt(params, 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if id < 1 {
		app.notFoundResponse(w, r)
		return
	}
	listing, err := app.models.Listings.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)

		}
		return
	}

	images, err := app.models.Images.GetForListing(listing.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"listing": listing, "listingImages": images}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getAllUserListingsHandler(w http.ResponseWriter, r *http.Request) {
	session := app.contextGetUser(r)
	listings, err := app.models.Listings.AllUserListings(session.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	for _, listing := range listings {
		images, err := app.models.Images.GetForListing(listing.ID)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		listing.Images = images
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"listings": listings}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
