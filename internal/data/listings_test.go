package data

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestListingsModel_Insert(t *testing.T) {
	user := CreateRandomUser(t)
	CreateRandomListing(t, user)
}

func TestListingsModel_Get(t *testing.T) {
	user := CreateRandomUser(t)
	listing := CreateRandomListing(t, user)

	listingFromDB, err := testQueries.Listings.Get(listing.ID)
	require.NoError(t, err)
	require.NotEmpty(t, listingFromDB)

	require.Equal(t, listing.ID, listingFromDB.ID)
	require.Equal(t, listing.CreatedAt, listingFromDB.CreatedAt)
	require.Equal(t, listing.Title, listingFromDB.Title)
	require.Equal(t, listing.Description, listingFromDB.Description)
	require.Equal(t, listing.Category, listingFromDB.Category)
	require.Equal(t, listing.RoomCount, listingFromDB.RoomCount)
	require.Equal(t, listing.BathroomCount, listingFromDB.BathroomCount)
	require.Equal(t, listing.GuestCount, listingFromDB.GuestCount)
	require.Equal(t, listing.Location, listingFromDB.Location)
	require.Equal(t, listing.Price, listingFromDB.Price)
	require.Equal(t, listing.OwnerID, listingFromDB.OwnerID)
	require.Equal(t, listing.OwnerName, listingFromDB.OwnerName)
	require.Equal(t, listing.OwnerPhoto, listingFromDB.OwnerPhoto)

}

func TestListingsModel_AllUserListings(t *testing.T) {
	user := CreateRandomUser(t)
	CreateRandomListing(t, user)
	CreateRandomListing(t, user)
	CreateRandomListing(t, user)

	listings, err := testQueries.Listings.AllUserListings(user.ID)
	require.NoError(t, err)
	require.NotEmpty(t, listings)
	require.Len(t, listings, 3)
}

func TestListingsModel_Delete(t *testing.T) {
	user := CreateRandomUser(t)
	listing := CreateRandomListing(t, user)

	err := testQueries.Listings.Delete(listing.ID, user.ID)
	require.NoError(t, err)

	listingFromDB, err := testQueries.Listings.Get(listing.ID)
	require.Error(t, err)
	require.Empty(t, listingFromDB)
}
