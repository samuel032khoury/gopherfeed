package cache

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/samuel032khoury/gopherfeed/internal/store"
)

type MockCacheStorage struct {
	Users interface {
		Get(context.Context, int64) (*store.User, error)
		Set(context.Context, *store.User) error
	}
}

func NewMockCacheStorage(client *redis.Client) *CacheStorage {
	return &CacheStorage{
		Users: &MockUserCache{},
	}
}

type MockUserCache struct {
}

func (m *MockUserCache) Get(ctx context.Context, id int64) (*store.User, error) {
	return nil, nil
}

func (m *MockUserCache) Set(ctx context.Context, user *store.User) error {
	return nil
}
