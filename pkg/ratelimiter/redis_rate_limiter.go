package ratelimiter

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type IRateLimiter interface {
	Allow(ctx context.Context, key string) (bool, error)
}

type RedisRateLimiter struct {
	rdb        *redis.Client
	limit      int           // max requests
	windowSize time.Duration // vd: 1 phút
}

func NewRedisRateLimiter(rdb *redis.Client, limit int, window time.Duration) *RedisRateLimiter {
	return &RedisRateLimiter{
		rdb:        rdb,
		limit:      limit,
		windowSize: window,
	}
}

func (r *RedisRateLimiter) Allow(ctx context.Context, key string) (bool, error) {
	// window hiện tại
	now := time.Now()
	window := now.Unix() / int64(r.windowSize.Seconds())

	redisKey := fmt.Sprintf("rate_limit:%s:%d", key, window)

	// tăng counter
	count, err := r.rdb.Incr(ctx, redisKey).Result()
	if err != nil {
		return false, err
	}

	// set TTL nếu lần đầu
	if count == 1 {
		r.rdb.Expire(ctx, redisKey, r.windowSize)
	}

	if count > int64(r.limit) {
		return false, nil
	}

	return true, nil
}
