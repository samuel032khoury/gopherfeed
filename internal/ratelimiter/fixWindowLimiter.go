package ratelimiter

import (
	"sync"
	"time"
)

type FixedWindowLimiter struct {
	sync.RWMutex
	quota    int
	duration time.Duration
	clients  map[string]int
}

func NewFixedWindowLimiter(quota int, interval string) (*FixedWindowLimiter, error) {
	duration, err := time.ParseDuration(interval)
	if err != nil {
		return nil, err
	}

	limiter := &FixedWindowLimiter{
		quota:    quota,
		duration: duration,
		clients:  make(map[string]int),
	}

	return limiter, nil
}

func (f *FixedWindowLimiter) Allow(key string) (bool, time.Duration) {
	f.Lock()
	defer f.Unlock()
	count, exists := f.clients[key]
	if !exists {
		go func() {
			time.Sleep(f.duration)
			f.Lock()
			defer f.Unlock()
			delete(f.clients, key)
		}()
	}
	if count < f.quota {
		f.clients[key]++
		return true, 0
	}
	return false, f.duration
}
