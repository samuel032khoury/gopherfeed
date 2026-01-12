package ratelimiter

import "time"

type Limiter interface {
	Allow(key string) (bool, time.Duration)
}
