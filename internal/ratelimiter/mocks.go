package ratelimiter

import "time"

type MockRateLimiter struct{}

func NewMockRateLimiter() Limiter {
	return &MockRateLimiter{}
}

func (m *MockRateLimiter) Allow(key string) (bool, time.Duration) {
	return true, 0
}
