package store

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type Post struct {
	ID        int64      `json:"id"`
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	UserID    int64      `json:"user_id"`
	Tags      []string   `json:"tags"`
	CreatedAt string     `json:"created_at"`
	UpdatedAt string     `json:"updated_at"`
	Version   int        `json:"version"`
	Comments  []*Comment `json:"comments"`
}

type FeedablePost struct {
	Post
	CommentsCount int    `json:"comments_count"`
	Username      string `json:"username"`
}

type PostStore struct {
	db *sql.DB
}

func (s *PostStore) Create(ctx context.Context, post *Post) error {
	query := `
		INSERT INTO posts (title, content, user_id, tags)
		VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	return s.db.QueryRowContext(
		ctx,
		query,
		post.Title,
		post.Content,
		post.UserID,
		pq.Array(post.Tags),
	).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)
}

func (s *PostStore) GetByID(ctx context.Context, id int64) (*Post, error) {
	query := `
		SELECT id, title, content, user_id, tags, created_at, updated_at, version
		FROM posts
		WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	post := &Post{}
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.UserID,
		pq.Array(&post.Tags),
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.Version,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return post, nil
}

func (s *PostStore) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM posts WHERE id = $1`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	_, err := s.db.ExecContext(ctx, query, id)
	return err
}

func (s *PostStore) Update(ctx context.Context, post *Post) error {
	query := `
		UPDATE posts
		SET title = $1, content = $2, tags = $3, updated_at = NOW(), version = version + 1
		WHERE id = $4 AND version = $5
		RETURNING updated_at, version
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	return s.db.QueryRowContext(
		ctx,
		query,
		post.Title,
		post.Content,
		pq.Array(post.Tags),
		post.ID,
		post.Version,
	).Scan(&post.UpdatedAt, &post.Version)
}

func (s *PostStore) GetFeed(ctx context.Context, userID int64) ([]*FeedablePost, error) {
	query := `
		SELECT p.id, p.title, p.content, p.user_id, p.tags, p.created_at, p.updated_at, p.version, u.username,
		       COUNT(c.id) AS comments_count
		FROM posts p
		LEFT JOIN comments c ON p.id = c.post_id
		LEFT JOIN users u ON p.user_id = u.id
		WHERE p.user_id = $1 OR p.user_id IN (
			SELECT followee_id FROM followers WHERE user_id = $1
		)
		GROUP BY p.id, u.username
		ORDER BY p.created_at DESC
		LIMIT 100
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feed []*FeedablePost
	for rows.Next() {
		post := &FeedablePost{}
		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.UserID,
			pq.Array(&post.Tags),
			&post.CreatedAt,
			&post.UpdatedAt,
			&post.Version,
			&post.Username,
			&post.CommentsCount,
		)
		if err != nil {
			return nil, err
		}
		feed = append(feed, post)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return feed, nil
}
