package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type BookingModel struct {
	DB *sql.DB
}

type Booking struct {
	ID        int64     `json:"id"`
	CreatedAt string    `json:"createdAt"`
	ListingID int64     `json:"listingId"`
	GuestID   int64     `json:"guestId"`
	CheckIn   time.Time `json:"checkIn"`
	CheckOut  time.Time `json:"checkOut"`
	Price     int64     `json:"price"`
	Total     int64     `json:"total"`
	Listing   Listing   `json:"listing"`
}

func (m BookingModel) Insert(booking *Booking) error {
	query := `INSERT INTO bookings (listing_id, guest_id, check_in, check_out, price, total)
    	VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at`

	args := []interface{}{booking.ListingID, booking.GuestID, booking.CheckIn, booking.CheckOut, booking.Price, booking.Total}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&booking.ID, &booking.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (m BookingModel) Get(id int64) (*Booking, error) {
	query := `SELECT b.id, b.created_at, b.listing_id, b.guest_id, b.check_in, b.check_out,
			  b.price, b.total, l.id, l.title, l.description, l.category, l.bedrooms,
			  l.bathrooms, l.guests, l.location_flag, l.location_label, l.location_lat, l.location_lng,	
			  l.location_region, l.location_value, l.price, l.owner_id, u.name, COALESCE(u.image, '')
			  FROM bookings b
			  INNER JOIN listings l ON l.id = b.listing_id
			  INNER JOIN users u ON u.id = l.owner_id
			  WHERE b.id = $1`

	var booking Booking

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&booking.ID,
		&booking.CreatedAt,
		&booking.ListingID,
		&booking.GuestID,
		&booking.CheckIn,
		&booking.CheckOut,
		&booking.Price,
		&booking.Total,
		&booking.Listing.ID,
		&booking.Listing.Title,
		&booking.Listing.Description,
		&booking.Listing.Category,
		&booking.Listing.Bedrooms,
		&booking.Listing.Bathrooms,
		&booking.Listing.Guests,
		&booking.Listing.Location.Flag,
		&booking.Listing.Location.Label,
		&booking.Listing.Location.Lat,
		&booking.Listing.Location.Lng,
		&booking.Listing.Location.Region,
		&booking.Listing.Location.Value,
		&booking.Listing.Price,
		&booking.Listing.OwnerID,
		&booking.Listing.OwnerName,
		&booking.Listing.OwnerPhoto,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &booking, nil
}

func (m BookingModel) Delete(id int64) error {
	query := `DELETE FROM bookings WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (m BookingModel) GetForUser(userID int64) ([]*Booking, error) {
	query := `SELECT b.id, b.created_at, b.listing_id, b.guest_id, b.check_in, b.check_out,
			  b.price, b.total, l.id, l.title, l.description, l.category, l.bedrooms,
			  l.bathrooms, l.guests, l.location_flag, l.location_label, l.location_lat, l.location_lng,	
			  l.location_region, l.location_value, l.price, l.owner_id, u.name, COALESCE(u.image, '')
			  FROM bookings b
			  INNER JOIN listings l ON l.id = b.listing_id
			  INNER JOIN users u ON u.id = l.owner_id
			  WHERE b.guest_id = $1
			  ORDER BY b.created_at DESC`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, userID)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	var bookings []*Booking

	for rows.Next() {
		var booking Booking
		err := rows.Scan(
			&booking.ID,
			&booking.CreatedAt,
			&booking.ListingID,
			&booking.GuestID,
			&booking.CheckIn,
			&booking.CheckOut,
			&booking.Price,
			&booking.Total,
			&booking.Listing.ID,
			&booking.Listing.Title,
			&booking.Listing.Description,
			&booking.Listing.Category,
			&booking.Listing.Bedrooms,
			&booking.Listing.Bathrooms,
			&booking.Listing.Guests,
			&booking.Listing.Location.Flag,
			&booking.Listing.Location.Label,
			&booking.Listing.Location.Lat,
			&booking.Listing.Location.Lng,
			&booking.Listing.Location.Region,
			&booking.Listing.Location.Value,
			&booking.Listing.Price,
			&booking.Listing.OwnerID,
			&booking.Listing.OwnerName,
			&booking.Listing.OwnerPhoto,
		)
		if err != nil {
			return nil, err
		}

		bookings = append(bookings, &booking)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return bookings, nil
}

func (m BookingModel) GetForListing(listingID int64) ([]*Booking, error) {
	query := `SELECT b.id, b.created_at, b.listing_id, b.guest_id, b.check_in, b.check_out,
			  b.price, b.total, l.id, l.title, l.description, l.category, l.bedrooms,
			  l.bathrooms, l.guests, l.location_flag, l.location_label, l.location_lat, l.location_lng,	
			  l.location_region, l.location_value, l.price, l.owner_id, u.name, COALESCE(u.image, '')
			  FROM bookings b
			  INNER JOIN listings l ON l.id = b.listing_id
			  INNER JOIN users u ON u.id = b.guest_id
			  WHERE b.listing_id = $1
			  ORDER BY b.created_at DESC`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, listingID)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	var bookings []*Booking

	for rows.Next() {
		var booking Booking
		err := rows.Scan(
			&booking.ID,
			&booking.CreatedAt,
			&booking.ListingID,
			&booking.GuestID,
			&booking.CheckIn,
			&booking.CheckOut,
			&booking.Price,
			&booking.Total,
			&booking.Listing.ID,
			&booking.Listing.Title,
			&booking.Listing.Description,
			&booking.Listing.Category,
			&booking.Listing.Bedrooms,
			&booking.Listing.Bathrooms,
			&booking.Listing.Guests,
			&booking.Listing.Location.Flag,
			&booking.Listing.Location.Label,
			&booking.Listing.Location.Lat,
			&booking.Listing.Location.Lng,
			&booking.Listing.Location.Region,
			&booking.Listing.Location.Value,
			&booking.Listing.Price,
			&booking.Listing.OwnerID,
			&booking.Listing.OwnerName,
			&booking.Listing.OwnerPhoto,
		)
		if err != nil {
			return nil, err
		}

		bookings = append(bookings, &booking)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return bookings, nil
}
