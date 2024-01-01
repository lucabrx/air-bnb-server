package data

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	Users    UserModel
	Tokens   TokenModel
	Listings ListingsModel
	Images   ImageModel
	Bookings BookingModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Users:    UserModel{DB: db},
		Tokens:   TokenModel{DB: db},
		Listings: ListingsModel{DB: db},
		Images:   ImageModel{DB: db},
		Bookings: BookingModel{DB: db},
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

func SavePgFloatArray(arr []float64) string {
	var strArr []string
	for _, v := range arr {
		strArr = append(strArr, strconv.FormatFloat(v, 'f', -1, 64))
	}
	return "{" + strings.Join(strArr, ",") + "}"
}

func LoadPgFloatArray(arr string) *[]float64 {
	var floatArr []float64
	arr = strings.Trim(arr, "{}")
	strArr := strings.Split(arr, ",")
	for _, v := range strArr {
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		floatArr = append(floatArr, f)
	}
	return &floatArr
}
