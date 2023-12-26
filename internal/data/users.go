package data

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"time"

	"github.com/air-bnb/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

var ErrDuplicateEmail = errors.New("duplicate email")
var AnonymousUser = &User{}

type UserModel struct {
	DB *sql.DB
}

type password struct {
	plaintext *string
	hash      []byte
}

type User struct {
	ID                int64     `json:"id"`
	CreatedAt         time.Time `json:"createdAt"`
	Name              string    `json:"name,omitempty"`
	Email             string    `json:"email"`
	Image             string    `json:"image,omitempty"`
	Password          password  `json:"-"`
	Activated         bool      `json:"activated"`
	VerificationToken string    `json:"verificationToken,omitempty"`
	ResetToken        string    `json:"resetToken,omitempty"`
	ResetEmailToken   string    `json:"resetEmailToken,omitempty"`
}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	p.plaintext = &plaintextPassword
	p.hash = hash

	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Name != "", "name", "must be provided")
	v.Check(len(user.Name) <= 500, "name", "must not be more than 500 bytes long")

	ValidateEmail(v, user.Email)

	if user.Password.plaintext != nil {
		ValidatePasswordPlaintext(v, *user.Password.plaintext)
	}

	if user.Password.hash == nil {
		panic("missing password hash for user")
	}
}

func (m UserModel) Insert(user *User) error {
	query := `
        INSERT INTO users (name, email, password_hash, activated, image, verification_token, reset_token, update_email_token) 
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING id, created_at`

	args := []interface{}{
		NewNullString(user.Name),
		user.Email,
		user.Password.hash,
		user.Activated,
		NewNullString(user.Image),
		NewNullString(user.VerificationToken),
		NewNullString(user.ResetToken),
		NewNullString(user.ResetEmailToken),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		switch {
		case err.Error() == `ERROR: duplicate key value violates unique constraint "users_email_key" (SQLSTATE 23505)`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}

	return nil
}

func (m UserModel) Get(id int64, email string) (*User, error) {
	query := `SELECT id, created_at, COALESCE(name, ''), email, COALESCE(image, ''),
       		  COALESCE(password_hash, ''), activated, COALESCE( verification_token, ''),
       		  COALESCE(reset_token, ''), COALESCE(update_email_token, '')
			  FROM users
			  WHERE id = $1 OR email = $2`

	var user User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id, email).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Email,
		&user.Image,
		&user.Password.hash,
		&user.Activated,
		&user.VerificationToken,
		&user.ResetToken,
		&user.ResetEmailToken,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (m UserModel) Update(user *User) error {
	query := `
        UPDATE users 
        SET name = $1, email = $2, password_hash = $3, activated = $4, 
        image = $5, verification_token = $6, reset_token = $7, update_email_token = $8
        WHERE id = $9`

	args := []interface{}{
		NewNullString(user.Name),
		user.Email,
		NewNullByteSlice(user.Password.hash),
		user.Activated,
		NewNullString(user.Image),
		NewNullString(user.VerificationToken),
		NewNullString(user.ResetToken),
		NewNullString(user.ResetEmailToken),
		user.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		switch {
		case err.Error() == `ERROR: duplicate key value violates unique constraint "users_email_key" (SQLSTATE 23505)`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}

	return nil
}

func (m UserModel) Delete(id int64) error {
	query := `
		DELETE FROM users 
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (m UserModel) GetForToken(tokenScope, tokenPlaintext string) (*User, error) {
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))
	query := `
        SELECT u.id, u.activated, u.created_at, COALESCE(u.name, ''), u.email ,COALESCE(u.image,''), COALESCE(u.password_hash, ''), COALESCE(u.update_email_token, '')
        FROM users u
        INNER JOIN tokens t
        ON u.id = t.user_id
        WHERE t.hash = $1
        AND t.scope = $2 
        AND t.expiry > $3`

	args := []any{tokenHash[:], tokenScope, time.Now()}

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.Activated,
		&user.CreatedAt,
		&user.Name,
		&user.Email,
		&user.Image,
		&user.Password.hash,
		&user.ResetEmailToken,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}
