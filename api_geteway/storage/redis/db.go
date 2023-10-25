package redis

import (
	"context"
	"fmt"
	"main/config"
	"main/storage"

	"github.com/go-redis/cache/v9"
	goRedis "github.com/redis/go-redis/v9"
)

type redis struct {
	db    *cache.Cache
	cache *cacheRepo
}

func NewCache(ctx context.Context, cfg config.Config) (storage.CacheI, error) {

	redisClient := goRedis.NewClient(&goRedis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       0,
	})
	redisCache := cache.New(&cache.Options{
		Redis: redisClient,
	})

	return &redis{
		db: redisCache,
	}, nil

}

func (r *redis) Cache() storage.RedisI {
	if r.cache == nil {
		r.cache = NewCacheRepo(r.db)
	}
	return r.cache
}
