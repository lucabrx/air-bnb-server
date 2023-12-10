package data

import (
	"context"
	"database/sql"
	"errors"
	"github.com/air-bnb/internal/validator"
	"time"
)

type ListingsModel struct {
	DB *sql.DB
}

type Listing struct {
	ID            int64  `json:"id"`
	CreatedAt     string `json:"created_at"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	Category      string `json:"category"`
	RoomCount     int64  `json:"room_count"`
	BathroomCount int64  `json:"bathroom_count"`
	GuestCount    int64  `json:"guest_count"`
	Location      string `json:"location"`
	Price         int64  `json:"price"`
	OwnerID       int64  `json:"owner_id"`
	OwnerName     string `json:"owner_name"`
	OwnerPhoto    string `json:"owner_photo,omitempty"`
	Images        []*Image
}

func ValidateListing(v *validator.Validator, listing *Listing) {
	v.Check(listing.Title != "", "title", "must be provided")
	v.Check(len(listing.Title) <= 500, "title", "must not be more than 500 bytes long")
	v.Check(listing.Description != "", "description", "must be provided")
	v.Check(len(listing.Description) <= 5000, "description", "must not be more than 5000 bytes long")
	v.Check(listing.Category != "", "category", "must be provided")
	v.Check(len(listing.Category) <= 255, "category", "must not be more than 255 bytes long")
	v.Check(listing.RoomCount > 0, "room_count", "must be greater than zero")
	v.Check(listing.BathroomCount > 0, "bathroom_count", "must be greater than zero")
	v.Check(listing.GuestCount > 0, "guest_count", "must be greater than zero")
	v.Check(listing.Location != "", "location", "must be provided")
	v.Check(len(listing.Location) <= 255, "location", "must not be more than 255 bytes long")
	v.Check(listing.Price > 0, "price", "must be greater than zero")
	v.Check(listing.OwnerID > 0, "owner_id", "must be greater than zero")
}

func (m ListingsModel) Insert(listing *Listing) error {
	query := `INSERT INTO listings (title, description, category, room_count, bathroom_count,
              guest_count, location, price, owner_id)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id, created_at`
	args := []interface{}{
		listing.Title,
		listing.Description,
		listing.Category,
		listing.RoomCount,
		listing.BathroomCount,
		listing.GuestCount,
		listing.Location,
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
	query := `SELECT l.id, l.created_at, l.title, l.description, l.category, l.room_count,
			  l.bathroom_count, l.guest_count, l.location, l.price, l.owner_id, u.name, COALESCE(u.image, '')
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
		&listing.RoomCount,
		&listing.BathroomCount,
		&listing.GuestCount,
		&listing.Location,
		&listing.Price,
		&listing.OwnerID,
		&listing.OwnerName,
		&listing.OwnerPhoto,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &listing, nil
}

func (m ListingsModel) AllUserListings(userID int64) ([]*Listing, error) {
	query := `SELECT l.id, l.created_at, l.title, l.description, l.category, l.room_count,
			  l.bathroom_count, l.guest_count, l.location, l.price, l.owner_id, u.name, COALESCE(u.image, '')
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
			&listing.RoomCount,
			&listing.BathroomCount,
			&listing.GuestCount,
			&listing.Location,
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

	_, err := m.DB.ExecContext(ctx, query, id, ownerId)
	if err != nil {
		return err
	}

	return nil
}
