package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/air-bnb/internal/validator"
	"github.com/jackc/pgx/v5"
	"time"
)

type ListingsModel struct {
	DB *sql.DB
}

type Location struct {
	Flag   string  `json:"flag"`
	Label  string  `json:"label"`
	Lat    float64 `json:"lat"`
	Lng    float64 `json:"lng"`
	Region string  `json:"region"`
	Value  string  `json:"value"`
}

type Listing struct {
	ID          int64    `json:"id"`
	CreatedAt   string   `json:"created_at"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Category    string   `json:"category"`
	Bedrooms    int64    `json:"bedrooms"`
	Bathrooms   int64    `json:"bathrooms"`
	Guests      int64    `json:"guests"`
	Location    Location `json:"location"`
	Price       int64    `json:"price"`
	OwnerID     int64    `json:"ownerId"`
	OwnerName   string   `json:"ownerName"`
	OwnerPhoto  string   `json:"ownerPhoto,omitempty"`
	Images      []*Image `json:"images,omitempty"`
}

func ValidateListing(v *validator.Validator, listing *Listing) {
	v.Check(listing.Title != "", "title", "must be provided")
	v.Check(len(listing.Title) <= 500, "title", "must not be more than 500 characters long")

	v.Check(listing.Description != "", "description", "must be provided")
	v.Check(len(listing.Description) <= 5000, "description", "must not be more than 5000 characters long")

	v.Check(listing.Category != "", "category", "must be provided")
	v.Check(len(listing.Category) <= 255, "category", "must not be more than 255 characters long")

	v.Check(listing.Bedrooms > 0, "bedrooms", "must be greater than zero")
	v.Check(listing.Bathrooms > 0, "bathrooms", "must be greater than zero")
	v.Check(listing.Guests > 0, "guests", "must be greater than zero")

	v.Check(listing.Location.Flag != "", "location.flag", "must be provided")
	v.Check(len(listing.Location.Flag) <= 255, "location.flag", "must not be more than 255 characters long")

	v.Check(len(listing.Location.Label) <= 255, "location.label", "must not be more than 255 characters long")

	v.Check(len(listing.Location.Region) <= 255, "location.region", "must not be more than 255 characters long")

	v.Check(listing.Location.Lat != 0, "location.lat", "must be provided")
	v.Check(listing.Location.Lng != 0, "location.lng", "must be provided")

	v.Check(listing.Price > 0, "price", "must be greater than zero")
	v.Check(listing.OwnerID > 0, "owner_id", "must be greater than zero")
}

func (m ListingsModel) Insert(listing *Listing) error {
	query := `INSERT INTO listings (title, description, category, bedrooms, bathrooms,
              guests, location_flag, location_label, location_lat, location_lng, location_region, location_value,
              price, owner_id)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) RETURNING id, created_at`
	args := []interface{}{
		listing.Title,
		listing.Description,
		listing.Category,
		listing.Bedrooms,
		listing.Bathrooms,
		listing.Guests,
		listing.Location.Flag,
		listing.Location.Label,
		listing.Location.Lat,
		listing.Location.Lng,
		listing.Location.Region,
		listing.Location.Value,
		listing.Price,
		listing.OwnerID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&listing.ID, &listing.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (m ListingsModel) Get(id int64) (*Listing, error) {
	query := `SELECT l.id, l.created_at, l.title, l.description, l.category, l.bedrooms,
			  l.bathrooms, l.guests, l.location_flag, l.location_label, l.location_lat, l.location_lng,
			  l.location_region, l.location_value, l.price, l.owner_id, u.name, COALESCE(u.image, '')
			  FROM listings l
			  INNER JOIN users u ON u.id = l.owner_id
			  WHERE l.id = $1`

	var listing Listing
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&listing.ID,
		&listing.CreatedAt,
		&listing.Title,
		&listing.Description,
		&listing.Category,
		&listing.Bedrooms,
		&listing.Bathrooms,
		&listing.Guests,
		&listing.Location.Flag,
		&listing.Location.Label,
		&listing.Location.Lat,
		&listing.Location.Lng,
		&listing.Location.Region,
		&listing.Location.Value,
		&listing.Price,
		&listing.OwnerID,
		&listing.OwnerName,
		&listing.OwnerPhoto,
	)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &listing, nil
}

func (m ListingsModel) AllUserListings(userID int64) ([]*Listing, error) {
	query := `SELECT l.id, l.created_at, l.title, l.description, l.category, l.bedrooms,
			  l.bathrooms, l.guests, l.location_flag, l.location_label, l.location_lat, l.location_lng,
			  l.location_region, l.location_value, l.price, l.owner_id, u.name, COALESCE(u.image, '')
			  FROM listings l
			  INNER JOIN users u ON u.id = l.owner_id
			  WHERE l.owner_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var listings []*Listing
	for rows.Next() {
		var listing Listing
		err := rows.Scan(
			&listing.ID,
			&listing.CreatedAt,
			&listing.Title,
			&listing.Description,
			&listing.Category,
			&listing.Bedrooms,
			&listing.Bathrooms,
			&listing.Guests,
			&listing.Location.Flag,
			&listing.Location.Label,
			&listing.Location.Lat,
			&listing.Location.Lng,
			&listing.Location.Region,
			&listing.Location.Value,
			&listing.Price,
			&listing.OwnerID,
			&listing.OwnerName,
			&listing.OwnerPhoto,
		)
		if err != nil {
			return nil, err
		}

		listings = append(listings, &listing)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return listings, nil
}

func (m ListingsModel) Delete(id, ownerId int64) error {
	query := `DELETE FROM listings WHERE id = $1 AND owner_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id, ownerId)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (m ListingsModel) GetAll(search string, filters Filters) ([]*Listing, Metadata, error) {
	baseQuery := `SELECT count(*) OVER(), l.id, l.created_at, l.title, l.description, l.category, l.bedrooms,
				 l.bathrooms, l.guests, l.location_flag, l.location_label, l.location_lat, l.location_lng,
				 l.location_region, l.location_value, l.price, l.owner_id, u.name, COALESCE(u.image, '')
				 FROM listings l INNER JOIN users u ON u.id = l.owner_id`

	if search != "" {
		baseQuery += ` WHERE (l.title ILIKE '%' || $3 || '%'
 					   OR l.category ILIKE '%' || $3 || '%'
					   OR l.location_region ILIKE '%' || $3 || '%'
					   OR l.location_label ILIKE '%' || $3 || '%')`
	}

	advQuery := fmt.Sprintf(` ORDER BY %s %s `, filters.sortColumn(), filters.sortDirection())

	query := baseQuery + advQuery + ` LIMIT $1 OFFSET $2`

	args := []interface{}{filters.limit(), filters.offset()}

	if search != "" {
		args = append(args, search)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	var listings []*Listing

	for rows.Next() {
		var listing Listing

		err := rows.Scan(
			&totalRecords,
			&listing.ID,
			&listing.CreatedAt,
			&listing.Title,
			&listing.Description,
			&listing.Category,
			&listing.Bedrooms,
			&listing.Bathrooms,
			&listing.Guests,
			&listing.Location.Flag,
			&listing.Location.Label,
			&listing.Location.Lat,
			&listing.Location.Lng,
			&listing.Location.Region,
			&listing.Location.Value,
			&listing.Price,
			&listing.OwnerID,
			&listing.OwnerName,
			&listing.OwnerPhoto,
		)
		if err = rows.Err(); err != nil {
			return nil, Metadata{}, err
		}

		listings = append(listings, &listing)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return listings, metadata, nil
}

func (m ListingsModel) Update(listing *Listing) error {
	query := `UPDATE listings SET title = $1, description = $2, category = $3, bedrooms = $4,
			  bathrooms = $5, guests = $6, location_flag = $7, location_label = $8, location_lat = $9,
			  location_lng = $10, location_region = $11, location_value = $12, price = $13
			  WHERE id = $14`

	args := []interface{}{
		listing.Title,
		listing.Description,
		listing.Category,
		listing.Bedrooms,
		listing.Bathrooms,
		listing.Guests,
		listing.Location.Flag,
		listing.Location.Label,
		listing.Location.Lat,
		listing.Location.Lng,
		listing.Location.Region,
		listing.Location.Value,
		listing.Price,
		listing.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
