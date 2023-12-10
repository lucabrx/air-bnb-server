package data

import (
	"github.com/air-bnb/internal/random"
	"github.com/air-bnb/internal/validator"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPassword_Set(t *testing.T) {
	p := random.RandString(10)
	password := &password{}
	err := password.Set(p)
	require.NoError(t, err)
}
func TestPassword_Matches(t *testing.T) {
	p := random.RandString(10)
	password := &password{}
	err := password.Set(p)
	require.NoError(t, err)

	match, err := password.Matches(p)
	require.NoError(t, err)
	require.True(t, match)
}

func TestValidateEmail(t *testing.T) {
	email := random.RandString(10) + "@gmail.com"
	v := validator.New()
	ValidateEmail(v, email)
	require.Equal(t, 0, len(v.Errors))
}

func TestValidateEmail_Invalid(t *testing.T) {
	email := random.RandString(10)
	v := validator.New()
	ValidateEmail(v, email)
	require.Equal(t, 1, len(v.Errors))
}

func TestValidatePasswordPlaintext(t *testing.T) {
	password := random.RandString(10)
	v := validator.New()
	ValidatePasswordPlaintext(v, password)
	require.Equal(t, 0, len(v.Errors))
}

func TestValidateUser(t *testing.T) {
	user := &User{
		Email: random.RandString(10) + "@gmail.com",
		Name:  random.RandString(10),
	}
	err := user.Password.Set(random.RandString(10))
	require.NoError(t, err)

	v := validator.New()
	ValidateUser(v, user)
	require.Equal(t, 0, len(v.Errors))
}

func TestUserModel_Insert_Valid(t *testing.T) {
	CreateRandomUser(t)
}

func TestUserModel_Insert_EmailInUse(t *testing.T) {
	email := random.RandString(10) + "@gmail.com"
	user1 := &User{
		Email: email,
		Name:  random.RandString(10),
	}
	user2 := &User{
		Email: email,
		Name:  random.RandString(10),
	}

	err := testQueries.Users.Insert(user1)
	require.NoError(t, err)

	require.Regexp(t, validator.EmailRX, user1.Email)

	require.NotZero(t, user1.ID)
	require.NotZero(t, user1.CreatedAt)

	err = testQueries.Users.Insert(user2)
	require.Error(t, err)
}

func TestUserModel_Insert_HashCode(t *testing.T) {
	password := random.RandString(10)
	user := &User{
		Email: random.RandString(10) + "@gmail.com",
		Name:  random.RandString(10),
	}
	err := user.Password.Set(password)
	require.NoError(t, err)

	err = testQueries.Users.Insert(user)
	require.NoError(t, err)

	require.Regexp(t, validator.EmailRX, user.Email)

	require.NotZero(t, user.ID)
	require.NotZero(t, user.CreatedAt)
}

func TestUserModel_Get_ID(t *testing.T) {
	user := CreateRandomUser(t)
	user2, err := testQueries.Users.Get(user.ID, "")
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user.ID, user2.ID)
	require.Equal(t, user.CreatedAt, user2.CreatedAt)
	require.Equal(t, user.Name, user2.Name)
	require.Equal(t, user.Email, user2.Email)
}

func TestUserModel_Get_Email(t *testing.T) {
	user := CreateRandomUser(t)
	user2, err := testQueries.Users.Get(0, user.Email)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user.ID, user2.ID)
	require.Equal(t, user.CreatedAt, user2.CreatedAt)
	require.Equal(t, user.Name, user2.Name)
	require.Equal(t, user.Email, user2.Email)
}

func TestUserModel_Get_InvalidID(t *testing.T) {
	user, err := testQueries.Users.Get(0, "")
	require.Error(t, err)
	require.Empty(t, user)
	require.EqualError(t, err, ErrRecordNotFound.Error())
}

func TestUserModel_Update(t *testing.T) {
	user := CreateRandomUser(t)

	dbUser, err := testQueries.Users.Get(user.ID, "")
	require.NoError(t, err)
	require.NotEmpty(t, dbUser)

	dbUser.Name = random.RandString(10)
	dbUser.Email = random.RandString(10) + "@gmail.com"

	err = testQueries.Users.Update(dbUser)
	require.NoError(t, err)

}

func TestUserModel_Delete(t *testing.T) {
	user := CreateRandomUser(t)

	err := testQueries.Users.Delete(user.ID)
	require.NoError(t, err)

	user2, err := testQueries.Users.Get(user.ID, "")
	require.Error(t, err)
	require.Empty(t, user2)
}
