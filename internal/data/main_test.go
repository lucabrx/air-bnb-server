package data

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"github.com/air-bnb/internal/random"
	"github.com/air-bnb/internal/validator"
	"github.com/stretchr/testify/require"
	"time"

	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var testQueries Models

func TestMain(m *testing.M) {
	conn, err := sql.Open("pgx", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		panic(err)
	}

	testQueries = NewModels(conn)

	m.Run()
}

func CreateRandomUser(t *testing.T) User {
	user := &User{
		Email: random.RandString(10) + "@gmail.com",
		Name:  random.RandString(10),
	}

	err := testQueries.Users.Insert(user)
	require.NoError(t, err)

	require.Regexp(t, validator.EmailRX, user.Email)

	require.NotZero(t, user.ID)
	require.NotZero(t, user.CreatedAt)

	return *user
}

func CreateTokenForUser(t *testing.T, user User) Token {
	token := &Token{
		UserID: user.ID,
		Expiry: time.Now().Add(24 * time.Hour),
		Scope:  ScopeAuthentication,
	}
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	require.NoError(t, err)
	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]

	err = testQueries.Tokens.Insert(token)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	return *token
}

func CreateRandomListing(t *testing.T, user User) Listing {
	location := Location{
		Flag:   "HR",
		Lat:    23.5,
		Lng:    23.5,
		Label:  random.RandString(5),
		Region: random.RandString(5),
		Value:  random.RandString(5),
	}
	listing := Listing{
		OwnerID:     user.ID,
		OwnerName:   user.Name,
		Title:       random.RandString(10),
		Description: random.RandString(10),
		Category:    random.RandString(10),
		Price:       random.RandInt(1, 100),
		Guests:      random.RandInt(1, 10),
		Bedrooms:    random.RandInt(1, 10),
		Bathrooms:   random.RandInt(1, 10),
		Location:    location,
	}

	err := testQueries.Listings.Insert(&listing)
	require.NoError(t, err)

	require.NotZero(t, listing.ID)
	require.NotZero(t, listing.CreatedAt)

	return listing
}

func CreateRandomImage(t *testing.T, listingId int64) Image {
	image := Image{
		ListingID: listingId,
		Url:       random.RandString(10),
	}

	err := testQueries.Images.Insert(&image)
	require.NoError(t, err)

	require.NotZero(t, image.ID)

	return image
}
