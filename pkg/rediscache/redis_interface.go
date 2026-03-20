// user-management-api/pkg/redis/redis_interface.go
package rediscache

import (
	"context"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key string, dest any) error
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	DeleteByPattern(ctx context.Context, pattern string) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
}
