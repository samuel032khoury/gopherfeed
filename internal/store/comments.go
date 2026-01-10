package store

import (
	"context"
	"database/sql"
)

// Comment represents a comment on a post
//
//	@Description	Comment information
type Comment struct {
	ID        int64  `json:"id" example:"1"`
	PostID    int64  `json:"post_id" example:"1"`
	UserID    int64  `json:"user_id" example:"2"`
	Content   string `json:"content" example:"Great post!"`
	CreatedAt string `json:"created_at" example:"2026-01-06T07:22:18Z"`
}

type CommentStore struct {
	db *sql.DB
}

func (s *CommentStore) GetByPostID(ctx context.Context, postId int64) ([]*Comment, error) {
	query := `
		SELECT c.id, c.post_id, c.user_id, c.content, c.created_at
		FROM comments c
		JOIN users ON users.id = c.user_id
		WHERE c.post_id = $1
		ORDER BY c.created_at DESC
	`
	ctx, cancel := withTimeout(ctx)
	defer cancel()
	rows, err := s.db.QueryContext(ctx, query, postId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []*Comment{}
	for rows.Next() {
		comment := &Comment{}

		err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.UserID,
			&comment.Content,
			&comment.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func (s *CommentStore) Create(ctx context.Context, comment *Comment) error {
	query := `
		INSERT INTO comments (post_id, user_id, content)
		VALUES ($1, $2, $3) RETURNING id, created_at
	`
	ctx, cancel := withTimeout(ctx)
	defer cancel()
	return s.db.QueryRowContext(
		ctx,
		query,
		comment.PostID,
		comment.UserID,
		comment.Content,
	).Scan(&comment.ID, &comment.CreatedAt)
}
