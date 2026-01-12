package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/samuel032khoury/gopherfeed/internal/store"
)

type CacheStorage struct {
	Users interface {
		Get(context.Context, int64) (*store.User, error)
		Set(context.Context, *store.User) error
	}
}

func NewRedisStorage(client *redis.Client) *CacheStorage {
	return &CacheStorage{
		Users: &UserCache{client: client, ttl: time.Hour},
	}
}
