package data

import (
	"github.com/air-bnb/internal/random"
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

func TestListingsModel_AllUserListings_Empty(t *testing.T) {
	user := CreateRandomUser(t)
	CreateRandomListing(t, user)
	CreateRandomListing(t, user)
	CreateRandomListing(t, user)

	listings, err := testQueries.Listings.AllUserListings(user.ID + 1)
	require.NoError(t, err)
	require.Empty(t, listings)
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

func TestListingsModel_Delete_NotFound(t *testing.T) {
	user := CreateRandomUser(t)
	listing := CreateRandomListing(t, user)

	err := testQueries.Listings.Delete(listing.ID, user.ID+1)
	require.Error(t, err)
	require.EqualError(t, err, ErrRecordNotFound.Error())
}

func TestListingsModel_Update(t *testing.T) {
	user := CreateRandomUser(t)
	listing := CreateRandomListing(t, user)

	listing.Title = random.RandString(6)
	listing.Description = random.RandString(16)
	listing.Price = random.RandInt(100, 1000)

	err := testQueries.Listings.Update(&listing)
	require.NoError(t, err)

	listingFromDB, err := testQueries.Listings.Get(listing.ID)
	require.NoError(t, err)
	require.NotEmpty(t, listingFromDB)

	require.Equal(t, listing.ID, listingFromDB.ID)
	require.Equal(t, listing.CreatedAt, listingFromDB.CreatedAt)
	require.Equal(t, listing.Title, listingFromDB.Title)
	require.Equal(t, listing.Description, listingFromDB.Description)
	require.Equal(t, listing.Category, listingFromDB.Category)
}

func TestListingModel_Update_NothingChanged(t *testing.T) {
	listing := Listing{}

	err := testQueries.Listings.Update(&listing)
	require.Error(t, err)
	require.EqualError(t, err, ErrRecordNotFound.Error())
}

func TestListingsModel_GetAll(t *testing.T) {
	user := CreateRandomUser(t)

	for i := 0; i < 10; i++ {
		CreateRandomListing(t, user)
	}

	var filters Filters
	filters.Page = 1
	filters.PageSize = 10
	filters.SortSafelist = []string{"id", "-id"}
	filters.Sort = "id"

	listings, _, err := testQueries.Listings.GetAll("", filters)
	require.NoError(t, err)
	require.NotEmpty(t, listings)
	require.Len(t, listings, 10)
}
