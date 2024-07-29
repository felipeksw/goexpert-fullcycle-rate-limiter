package repository

import "time"

type RateLimiterStorage interface {
	Set(key string, value any, duration int32) (string, error)
	Increment(key string, expiration time.Duration) (int64, error)
}
