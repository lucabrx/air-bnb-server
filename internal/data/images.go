package data

import (
	"context"
	"database/sql"
	"time"
)

type ImageModel struct {
	DB *sql.DB
}

type Image struct {
	ID        int64  `json:"id"`
	ListingID int64  `json:"listingId"`
	Url       string `json:"url"`
}

func (m *ImageModel) Insert(image *Image) error {
	query := `INSERT INTO images (listing_id, url) VALUES($1, $2) RETURNING id`
	args := []interface{}{image.ListingID, image.Url}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&image.ID)
	if err != nil {
		return err
	}

	return nil
}

func (m *ImageModel) Delete(id int64) error {
	query := `DELETE FROM images WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
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

func (m *ImageModel) GetForListing(listingID int64) ([]*Image, error) {
	query := `SELECT id, listing_id, url FROM images WHERE listing_id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, listingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []*Image
	for rows.Next() {
		var image Image
		err := rows.Scan(&image.ID, &image.ListingID, &image.Url)
		if err != nil {
			return nil, err
		}
		images = append(images, &image)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return images, nil
}
