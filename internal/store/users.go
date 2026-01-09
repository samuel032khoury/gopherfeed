package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
)

// User represents a user in the system
//
//	@Description	User account information
type User struct {
	ID        int64  `json:"id" example:"1"`
	Username  string `json:"username" example:"john_doe"`
	Email     string `json:"email" example:"john@example.com"`
	Password  string `json:"-"`
	CreatedAt string `json:"created_at" example:"2026-01-06T07:22:18Z"`
}

type UserStore struct {
	db *sql.DB
}

var (
	ErrDuplicateEmail    = errors.New("user with that email already exists")
	ErrDuplicateUsername = errors.New("user with that username already exists")
)

type execer interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func (s *UserStore) Create(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `
		INSERT INTO users (username, email, password_hash)
		VALUES ($1, $2, $3) RETURNING id, created_at
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var execer execer

	if tx != nil {
		execer = tx
	} else {
		execer = s.db
	}
	err := execer.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Email,
		user.Password,
	).Scan(&user.ID, &user.CreatedAt)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			// PostgreSQL unique constraint violation
			if pqErr.Code == "23505" {
				switch pqErr.Constraint {
				case "users_email_key":
					return ErrDuplicateEmail
				case "users_username_key":
					return ErrDuplicateUsername
				default:
					return errors.New("user already exists")
				}
			}
		}
		return err
	}
	return nil
}

func (s *UserStore) createUserInvitation(ctx context.Context, tx *sql.Tx, token string, userID int64, exp time.Duration) error {
	query := `
		INSERT INTO user_invitations (token, user_id, expires_at)
		VALUES ($1, $2, $3)
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	_, err := tx.ExecContext(ctx, query, token, userID, time.Now().Add(exp))
	return err
}

func (s *UserStore) GetByID(ctx context.Context, id int64) (*User, error) {
	query := `
		SELECT id, username, email, password_hash, created_at
		FROM users
		WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	user := &User{}
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (s *UserStore) Register(ctx context.Context, user *User, token string, exp time.Duration) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.Create(ctx, tx, user); err != nil {
			return err
		}
		if err := s.createUserInvitation(ctx, tx, token, user.ID, exp); err != nil {
			return err
		}
		return nil
	})
}
