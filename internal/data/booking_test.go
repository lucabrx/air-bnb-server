package data

import (
	"github.com/air-bnb/internal/random"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestBookingModel_Insert(t *testing.T) {
	user := CreateRandomUser(t)
	listing := CreateRandomListing(t, user)
	randCheckIn := int(random.RandInt(1, 100))
	randCheckOut := int(random.RandInt(1, 100))
	booking := &Booking{
		ListingID: listing.ID,
		GuestID:   user.ID,
		CheckIn:   time.Now().AddDate(0, 0, randCheckIn),
		CheckOut:  time.Now().AddDate(0, 0, randCheckOut),
		Price:     listing.Price,
		Total:     listing.Price * int64(randCheckOut-randCheckIn),
	}

	err := testQueries.Bookings.Insert(booking)
	if err != nil {
		t.Fatal(err)
	}
	require.NotEmpty(t, booking.ID)
	require.NotEmpty(t, booking.CreatedAt)
}

func TestBookingModel_Delete(t *testing.T) {
	user := CreateRandomUser(t)
	listing := CreateRandomListing(t, user)
	randCheckIn := int(random.RandInt(1, 100))
	randCheckOut := int(random.RandInt(1, 100))
	booking := &Booking{
		ListingID: listing.ID,
		GuestID:   user.ID,
		CheckIn:   time.Now().AddDate(0, 0, randCheckIn),
		CheckOut:  time.Now().AddDate(0, 0, randCheckOut),
		Price:     listing.Price,
		Total:     listing.Price * int64(randCheckOut-randCheckIn),
	}

	err := testQueries.Bookings.Insert(booking)
	require.NoError(t, err)

	err = testQueries.Bookings.Delete(booking.ID, user.ID)
	require.NoError(t, err)
}

func TestBookingModel_Get(t *testing.T) {
	user := CreateRandomUser(t)
	listing := CreateRandomListing(t, user)
	randCheckIn := int(random.RandInt(1, 100))
	randCheckOut := int(random.RandInt(1, 100))
	booking := &Booking{
		ListingID: listing.ID,
		GuestID:   user.ID,
		CheckIn:   time.Now().AddDate(0, 0, randCheckIn),
		CheckOut:  time.Now().AddDate(0, 0, randCheckOut),
		Price:     listing.Price,
		Total:     listing.Price * int64(randCheckOut-randCheckIn),
	}

	err := testQueries.Bookings.Insert(booking)
	require.NoError(t, err)

	booking, err = testQueries.Bookings.Get(booking.ID)
	require.NoError(t, err)
	require.NotEmpty(t, booking)
}

func TestBookingModel_GetForListing(t *testing.T) {
	user := CreateRandomUser(t)
	listing := CreateRandomListing(t, user)
	for i := 0; i < 10; i++ {
		randCheckIn := int(random.RandInt(1, 100))
		randCheckOut := int(random.RandInt(1, 100))
		booking := &Booking{
			ListingID: listing.ID,
			GuestID:   user.ID,
			CheckIn:   time.Now().AddDate(0, 0, randCheckIn),
			CheckOut:  time.Now().AddDate(0, 0, randCheckOut),
			Price:     listing.Price,
			Total:     listing.Price * int64(randCheckOut-randCheckIn),
		}

		err := testQueries.Bookings.Insert(booking)
		require.NoError(t, err)
		require.NotEmpty(t, booking.ID)
		require.NotEmpty(t, booking.CreatedAt)

	}

	bookings, err := testQueries.Bookings.GetForListing(listing.ID)
	require.NoError(t, err)
	require.NotEmpty(t, bookings)
	require.Len(t, bookings, 10)

}

func TestBookingModel_GetForUser(t *testing.T) {
	user := CreateRandomUser(t)
	listing := CreateRandomListing(t, user)
	for i := 0; i < 10; i++ {
		randCheckIn := int(random.RandInt(1, 100))
		randCheckOut := int(random.RandInt(1, 100))
		booking := &Booking{
			ListingID: listing.ID,
			GuestID:   user.ID,
			CheckIn:   time.Now().AddDate(0, 0, randCheckIn),
			CheckOut:  time.Now().AddDate(0, 0, randCheckOut),
			Price:     listing.Price,
			Total:     listing.Price * int64(randCheckOut-randCheckIn),
		}

		err := testQueries.Bookings.Insert(booking)
		require.NoError(t, err)
		require.NotEmpty(t, booking.ID)
		require.NotEmpty(t, booking.CreatedAt)

	}

	bookings, err := testQueries.Bookings.GetForUser(user.ID)
	require.NoError(t, err)
	require.NotEmpty(t, bookings)
	require.Len(t, bookings, 10)
}
