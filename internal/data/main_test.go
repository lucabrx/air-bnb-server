package data

import (
	"database/sql"
	"github.com/air-bnb/internal/random"
	"github.com/air-bnb/internal/validator"
	"github.com/stretchr/testify/require"

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
