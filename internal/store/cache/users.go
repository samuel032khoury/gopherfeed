package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/samuel032khoury/gopherfeed/internal/store"
)

type UserCache struct {
	client *redis.Client
	ttl    time.Duration
}

func (c *UserCache) Get(ctx context.Context, userID int64) (*store.User, error) {
	userCacheKey := fmt.Sprintf("user:%d", userID)
	data, err := c.client.Get(ctx, userCacheKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	if data == "" {
		return nil, nil
	}
	var user store.User
	if err := json.Unmarshal([]byte(data), &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (c *UserCache) Set(ctx context.Context, user *store.User) error {
	userCacheKey := fmt.Sprintf("user:%d", user.ID)
	json, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return c.client.SetEX(ctx, userCacheKey, json, c.ttl).Err()
}
