package data

import (
	"context"
	"database/sql"
	"errors"
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
