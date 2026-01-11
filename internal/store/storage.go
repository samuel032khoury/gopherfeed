package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/samuel032khoury/gopherfeed/internal/auth"
)

const (
	QueryTimeoutDuration = 5 * time.Second
)

// Common database execution interface
type execer interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

// Common helper functions for all stores
func withTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, QueryTimeoutDuration)
}

func getExecer(db *sql.DB, tx *sql.Tx) execer {
	if tx != nil {
		return tx
	}
	return db
}

func prepareContext(ctx context.Context, db *sql.DB, tx *sql.Tx) (context.Context, context.CancelFunc, execer) {
	ctx, cancel := withTimeout(ctx)
	execer := getExecer(db, tx)
	return ctx, cancel, execer
}

type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetByID(context.Context, int64) (*Post, error)
		Delete(context.Context, int64) error
		Update(context.Context, *Post) error
		GetFeed(context.Context, int64, *PaginationParams) ([]*FeedablePost, error)
	}
	Users interface {
		Create(context.Context, *sql.Tx, *User) error
		GetByID(context.Context, int64) (*User, error)
		Register(context.Context, *User, string, time.Duration) error
		Authenticate(context.Context, string, string, auth.Authenticator) (string, error)
		Activate(context.Context, string) error
		Delete(context.Context, int64) error
	}
	Comments interface {
		GetByPostID(context.Context, int64) ([]*Comment, error)
		Create(context.Context, *Comment) error
	}
	Followers interface {
		Follow(context.Context, int64, int64) error
		Unfollow(context.Context, int64, int64) error
	}
}

func NewPostgresStorage(db *sql.DB) *Storage {
	return &Storage{
		Posts:     &PostStore{db: db},
		Users:     &UserStore{db: db},
		Comments:  &CommentStore{db: db},
		Followers: &FollowerStore{db: db},
	}
}

func withTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
