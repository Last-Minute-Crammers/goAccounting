package util

import (
	"context"
	"goAccounting/global/constant"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/patrickmn/go-cache"
)

type Cache interface {
	GetKey(tab constant.CacheTab, Unique string) string
	Init() error
	Get(key string) (any, bool)
	GetInt(key string) (int, bool)
	Set(key string, val any, duration time.Duration)
	Increment(key string, num int64) error
	Close() error
	Delete(key string) error
}

type cacheBase struct {
}

func (ch *cacheBase) GetKey(tab constant.CacheTab, unique string) string {
	return string(tab) + "_" + unique
}

type RedisCache struct {
	cacheBase
	client   *redis.Client
	DB       int
	Addr     string
	Password string
}

func (rc *RedisCache) Init() error {
	client := redis.NewClient(
		&redis.Options{
			Addr:     rc.Addr,
			Password: rc.Password,
			DB:       rc.DB,
		},
	)
	_, err := client.Ping(context.Background()).Result()
	return err
}

func (rc *RedisCache) Get(key string) (any, bool) {
	ctx := context.Background()

}

type LocalCache struct {
	cacheBase
	Cache *cache.Cache
}
