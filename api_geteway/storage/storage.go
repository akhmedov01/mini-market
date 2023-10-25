package storage

import (
	"context"
	"time"
)

type CacheI interface {
	Cache() RedisI
}

type RedisI interface {
	Create(context.Context, string, interface{}, time.Duration) error
	Get(context.Context, string, interface{}) (bool, error)
	Delete(context.Context, string) error
}
