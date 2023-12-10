package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	Users    UserModel
	Tokens   TokenModel
	Listings ListingsModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Users:    UserModel{DB: db},
		Tokens:   TokenModel{DB: db},
		Listings: ListingsModel{DB: db},
	}
}

func NewNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func NewNullByteSlice(b []byte) sql.NullString {
	if len(b) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: string(b),
		Valid:  true,
	}
}
