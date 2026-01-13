package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/samuel032khoury/gopherfeed/internal/auth"
)

func NewMockStore() Storage {
	return Storage{
		Posts: &MockPostStore{},
		Users: &MockUserStore{},
	}
}

type MockPostStore struct{}

func (m *MockPostStore) Create(ctx context.Context, post *Post) error {
	return nil
}
func (m *MockPostStore) GetByID(ctx context.Context, id int64) (*Post, error) {
	return &Post{}, nil
}
func (m *MockPostStore) Delete(ctx context.Context, id int64) error {
	return nil
}
func (m *MockPostStore) Update(ctx context.Context, post *Post) error {
	return nil
}
func (m *MockPostStore) GetFeed(ctx context.Context, userID int64, params *PaginationParams) ([]*FeedablePost, error) {
	return []*FeedablePost{}, nil
}

type MockUserStore struct{}

func (m *MockUserStore) Create(ctx context.Context, tx *sql.Tx, user *User) error {
	return nil
}
func (m *MockUserStore) GetByID(ctx context.Context, id int64) (*User, error) {
	return &User{
		ID:       id,
		Username: "testuser",
		Email:    "test@example.com",
		RoleID:   1,
		IsActive: true,
	}, nil
}
func (m *MockUserStore) Register(ctx context.Context, user *User, token string, exp time.Duration) error {
	return nil
}
func (m *MockUserStore) Authenticate(ctx context.Context, email string, password string, authenticator auth.Authenticator) (string, error) {
	return "", nil
}
func (m *MockUserStore) Activate(ctx context.Context, token string) error {
	return nil
}
func (m *MockUserStore) Delete(ctx context.Context, id int64) error {
	return nil
}
