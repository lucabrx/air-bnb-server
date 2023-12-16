package main

import (
	"errors"
	"github.com/air-bnb/internal/data"
	"github.com/air-bnb/internal/validator"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type Location struct {
	Flag   string    `json:"flag"`
	Label  string    `json:"label"`
	LatLng []float64 `json:"latlng"`
	Region string    `json:"region"`
	Value  string    `json:"value"`
}

func (app *application) createListingHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Bathrooms   int64    `json:"bathrooms"`
		Bedrooms    int64    `json:"bedrooms"`
		Category    string   `json:"category"`
		Guests      int64    `json:"guests"`
		Price       string   `json:"price"`
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Images      []string `json:"images"`
		Location    Location `json:"location"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}

	price, err := strconv.ParseInt(input.Price, 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	location := &data.Location{
		Flag:   input.Location.Flag,
		Label:  input.Location.Label,
		Lat:    input.Location.LatLng[0],
		Lng:    input.Location.LatLng[1],
		Region: input.Location.Region,
	}

	listing := &data.Listing{
		Bathrooms:   input.Bathrooms,
		Bedrooms:    input.Bedrooms,
		Category:    input.Category,
		Guests:      input.Guests,
		Price:       price,
		Title:       input.Title,
		Description: input.Description,
		OwnerID:     app.contextGetUser(r).ID,
		Location:    *location,
	}

	v := validator.New()
	data.ValidateListing(v, listing)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Listings.Insert(listing)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	for _, image := range input.Images {
		image := &data.Image{
			ListingID: listing.ID,
			Url:       image,
		}
		err = app.models.Images.Insert(image)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"id": listing.ID}, nil)
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
