package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

type RedisServer struct {
	client *redis.Client
}

func NewRedisCahe(client *redis.Client) *RedisServer {
	return &RedisServer{
		client: client,
	}
}

func (r *RedisServer) Set(key string, value any, expiration int32) (string, error) {
	err := r.client.Set(key, value, time.Duration(expiration)*time.Second).Err()
	if err != nil {
		return "", err
	}
	val, err := r.client.Get(key).Result()
	if err != nil {
		return "", err
	}
	if fmt.Sprintf("%v", value) != val {
		return "", errors.New("inconsistency to set the key")
	}
	return val, nil
}

func (r *RedisServer) Increment(key string, expiration time.Duration) (int64, error) {
	val, err := r.client.Incr(key).Result()
	if err != nil {
		return 0, err
	}
	if expiration > 0 {
		err := r.client.Expire(key, expiration*time.Second).Err()
		if err != nil {
			return 0, err
		}
	}
	return val, nil
}

func (r *RedisServer) Get(key string) (string, error) {
	val, err := r.client.Get(key).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			return "", nil
		}
		return "", err
	}
	return val, nil
}

func (r *RedisServer) Del(key string) (int64, error) {
	val, err := r.client.Del(key).Result()
	if err != nil {
		return 0, err
	}
	return val, nil
}
