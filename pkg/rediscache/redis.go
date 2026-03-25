// user-management-api/pkg/redis/redis.go
package rediscache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisCacheService struct {
	rdb *redis.Client
}

func New(rdb *redis.Client) Cache {
	return &redisCacheService{
		rdb: rdb,
	}
}

func (r *redisCacheService) DeleteByPattern(ctx context.Context, pattern string) error {
	var cursor uint64
	for {
		keys, nextCursor, err := r.rdb.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return err
		}

		if len(keys) > 0 {
			if err := r.rdb.Del(ctx, keys...).Err(); err != nil {
				return err
			}
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}
	return nil
}

func (r *redisCacheService) Get(ctx context.Context, key string, dest any) error {
	val, err := r.rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return err // cache miss
		}
		return err
	}

	return json.Unmarshal([]byte(val), dest)
}

func (r *redisCacheService) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.rdb.Set(ctx, key, data, ttl).Err()
}

func (r *redisCacheService) Delete(ctx context.Context, key string) error {
	return r.rdb.Del(ctx, key).Err()
}

func (r *redisCacheService) Exists(ctx context.Context, key string) (bool, error) {
	n, err := r.rdb.Exists(ctx, key).Result()
	return n > 0, err
}
