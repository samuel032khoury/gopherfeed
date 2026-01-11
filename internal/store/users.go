package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lib/pq"
	"github.com/samuel032khoury/gopherfeed/internal/auth"
	"github.com/samuel032khoury/gopherfeed/internal/utils"
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
	IsActive  bool   `json:"is_active" example:"false"`
	RoleID    int64  `json:"role_id" example:"1"`
}

type UserStore struct {
	db *sql.DB
}

var (
	ErrDuplicateEmail     = errors.New("user with that email already exists")
	ErrDuplicateUsername  = errors.New("user with that username already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidToken       = errors.New("invalid or expired token")
)

func (s *UserStore) getUserFromInvitation(ctx context.Context, token string) (*User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.created_at, u.is_active
		FROM users u
		JOIN user_invitations ui ON u.id = ui.user_id
		WHERE ui.token = $1 AND ui.expires_at > NOW()
	`
	ctx, cancel := withTimeout(ctx)
	defer cancel()

	tokenHash := utils.Hash(token)
	user := &User{}
	err := s.db.QueryRowContext(ctx, query, tokenHash).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.IsActive,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrInvalidToken
		}
		return nil, err
	}
	return user, nil
}

func (s *UserStore) Create(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `
		INSERT INTO users (username, email, password_hash, role_id)
		VALUES ($1, $2, $3, $4) RETURNING id, created_at
	`
	ctx, cancel, execer := prepareContext(ctx, s.db, tx)
	defer cancel()

	err := execer.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Email,
		user.Password,
		user.RoleID,
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

func (s *UserStore) update(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `UPDATE users SET username = $1, email = $2, is_active = $3 WHERE id = $4`
	ctx, cancel, execer := prepareContext(ctx, s.db, tx)
	defer cancel()

	_, err := execer.ExecContext(ctx, query, user.Username, user.Email, user.IsActive, user.ID)
	return err
}

func (s *UserStore) createUserInvitation(ctx context.Context, tx *sql.Tx, token string, userID int64, exp time.Duration) error {
	query := `
		INSERT INTO user_invitations (token, user_id, expires_at)
		VALUES ($1, $2, $3)
	`
	ctx, cancel, execer := prepareContext(ctx, s.db, tx)
	defer cancel()

	_, err := execer.ExecContext(ctx, query, token, userID, time.Now().Add(exp))
	return err
}

func (s *UserStore) deleteUserInvitation(ctx context.Context, tx *sql.Tx, userID int64) error {
	query := `DELETE FROM user_invitations WHERE user_id = $1`
	ctx, cancel, execer := prepareContext(ctx, s.db, tx)
	defer cancel()

	_, err := execer.ExecContext(ctx, query, userID)
	return err
}

func (s *UserStore) GetByID(ctx context.Context, id int64) (*User, error) {
	query := `
		SELECT id, username, email, password_hash, created_at, is_active, role_id
		FROM users
		WHERE id = $1 AND is_active = TRUE
	`
	ctx, cancel := withTimeout(ctx)
	defer cancel()

	user := &User{}
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.IsActive,
		&user.RoleID,
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

func (s *UserStore) Authenticate(ctx context.Context, email, password string, authenticator auth.Authenticator) (string, error) {
	query := `
		SELECT id, username, password_hash
		FROM users
		WHERE email = $1 AND is_active = TRUE
	`
	ctx, cancel := withTimeout(ctx)
	defer cancel()

	var id int64
	var username, passwordHash string
	err := s.db.QueryRowContext(ctx, query, email).Scan(&id, &username, &passwordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", ErrInvalidCredentials
		}
		return "", err
	}
	if !utils.CheckPasswordHash(password, passwordHash) {
		return "", ErrInvalidCredentials
	}
	exp, iss, aud := authenticator.GetMetadata()
	claims := jwt.MapClaims{
		"sub": id,
		"exp": time.Now().Add(exp).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": iss,
		"aud": aud,
	}
	token, err := authenticator.GenerateToken(claims)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *UserStore) Activate(ctx context.Context, token string) error {
	user, err := s.getUserFromInvitation(ctx, token)
	if err != nil {
		return err
	}
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		user.IsActive = true
		if err := s.update(ctx, tx, user); err != nil {
			return err
		}
		if err := s.deleteUserInvitation(ctx, tx, user.ID); err != nil {
			return err
		}
		return nil
	})
}

func (s *UserStore) Delete(ctx context.Context, id int64) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.deleteUser(ctx, tx, id); err != nil {
			return err
		}
		if err := s.deleteUserInvitation(ctx, tx, id); err != nil {
			return err
		}
		return nil
	})
}

func (s *UserStore) deleteUser(ctx context.Context, tx *sql.Tx, id int64) error {
	query := `DELETE FROM users WHERE id = $1`
	ctx, cancel, execer := prepareContext(ctx, s.db, tx)
	defer cancel()

	_, err := execer.ExecContext(ctx, query, id)
	return err
}
