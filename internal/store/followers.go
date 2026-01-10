package store

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type Follower struct {
	UserID     int64  `json:"user_id"`
	FolloweeID int64  `json:"followee_id"`
	CreatedAt  string `json:"created_at"`
}

type FollowerStore struct {
	db *sql.DB
}

func (s *FollowerStore) Follow(ctx context.Context, userID, followeeID int64) error {
	query := `INSERT INTO followers (user_id, followee_id) VALUES ($1, $2)`
	ctx, cancel := withTimeout(ctx)
	defer cancel()
	_, err := s.db.ExecContext(ctx, query, userID, followeeID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code.Name() == "unique_violation" {
				// Already following; ignore
				return nil
			}
		}
	}
	return err
}

func (s *FollowerStore) Unfollow(ctx context.Context, userID, followeeID int64) error {
	query := `DELETE FROM followers WHERE user_id = $1 AND followee_id = $2`
	ctx, cancel := withTimeout(ctx)
	defer cancel()
	_, err := s.db.ExecContext(ctx, query, userID, followeeID)
	return err
}
