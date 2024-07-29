package repository

import "time"

type RateLimiterStorage interface {
	Set(key string, value any, expiration int32) (string, error)
	Increment(key string, expiration time.Duration) (int64, error)
	Get(key string) (string, error)
	Del(key string) (int64, error)
	Debug(val string) string
}
