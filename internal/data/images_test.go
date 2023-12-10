package data

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestImageModel_Insert(t *testing.T) {
	user := CreateRandomUser(t)
	listing := CreateRandomListing(t, user)
	CreateRandomImage(t, listing.ID)
}

func TestImageModel_Delete(t *testing.T) {
	user := CreateRandomUser(t)
	listing := CreateRandomListing(t, user)
	image := CreateRandomImage(t, listing.ID)
	err := testQueries.Images.Delete(image.ID)
	require.NoError(t, err)

}

func TestImageModel_GetForListing(t *testing.T) {
	user := CreateRandomUser(t)
	listing := CreateRandomListing(t, user)
	CreateRandomImage(t, listing.ID)
	CreateRandomImage(t, listing.ID)
	CreateRandomImage(t, listing.ID)
	images, err := testQueries.Images.GetForListing(listing.ID)
	require.NoError(t, err)
	require.Len(t, images, 3)
}
