package util

import (
	"context"
	"fmt"
	"goAccounting/global/constant"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/patrickmn/go-cache"
)

// provide the interface to operate for RedisCache and LocalCache
type Cache interface {
	GetKey(tab constant.CacheTab, Unique string) string
	Init() error
	Get(key string) (any, bool)
	GetInt(key string) (int, bool)
	Set(key string, val any, duration time.Duration)
	Increment(key string, num int64) error
	Delete(key string) error
	Close() error
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
	val, err := rc.client.Get(ctx, key).Result()
	if err == redis.Nil {
		log.Println("Redis doesn't exist")
		return nil, false
	} else if err != nil {
		log.Printf("Error while getting key %s: %v\n", key, err)
		return nil, false
	}
	return val, true
}

func (rc *RedisCache) GetInt(key string) (result int, isSuccess bool) {
	ctx := context.Background()
	val, err := rc.client.Get(ctx, key).Result()
	if err == redis.Nil {
		log.Println("Redis doesn't exist")
		return
	} else if err != nil {
		log.Printf("Error while getting key %s: %v\n", key, err)
		return
	}
	result, err = convertToInt(val)
	if err != nil {
		log.Println("Error while convert val into INT")
	}
	return result, true
}

func (rc *RedisCache) Set(key string, val any, duration time.Duration) {
	ctx := context.Background()
	err := rc.client.Set(ctx, key, val, duration).Err()
	if err != nil {
		log.Printf("Error while setting key %s: %v\n", key, err)
	}
}

func (rc *RedisCache) Increment(key string, number int64) error {
	ctx := context.Background()
	_, err := rc.client.IncrBy(ctx, key, number).Result()
	if err != nil {
		log.Printf("Error while incrementing key %s: %v\n", key, err)
		return err
	}
	return nil
}

func (rc *RedisCache) Delete(key string) error {
	ctx := context.Background()
	_, err := rc.client.Del(ctx, key).Result()
	if err != nil {
		log.Printf("Error while deleting key %s: %v\n", key, err)
		return err
	}
	return nil
}
func (rc *RedisCache) Close() error {
	return rc.client.Close()
}

type LocalCache struct {
	cacheBase
	Cache *cache.Cache
}

func (lc *LocalCache) Init() error {
	lc.Cache = cache.New(2*time.Hour, 10*time.Minute)
	if lc.Cache == nil {
		return fmt.Errorf("failed to initialize local cache")
	}
	return nil
}

func (lc *LocalCache) Get(key string) (any, bool) {
	val, found := lc.Cache.Get(key)
	if !found {
		log.Printf("Key '%s' not found in local cache\n", key)
		return nil, false
	}
	return val, true
}

func (lc *LocalCache) GetInt(key string) (result int, isSuccess bool) {
	val, found := lc.Get(key)
	if !found {
		log.Printf("Key '%s' not found in local cache\n", key)
		return
	}
	result, err := convertToInt(val)
	if err != nil {
		log.Printf("Error while converting val into int")
		return
	}
	return result, isSuccess
}

func (lc *LocalCache) Set(key string, value interface{}, duration time.Duration) {
	lc.Cache.Set(key, value, duration)
	log.Printf("[INFO] Key '%s' set in local cache; Value: %v; Expiration: %v\n", key, value, duration)
}

func (lc *LocalCache) Increment(key string, number int64) error {
	log.Printf("[WARN] Increment is not supported for LocalCache")
	return fmt.Errorf("Increment is not supported for LocalCache")
}

func (lc *LocalCache) Delete(key string) error {
	_, found := lc.Cache.Get(key)
	if !found {
		log.Printf("[WARN] Attempted to delete key '%s', but it was not found in the local cache\n", key)
		return nil
	}
	lc.Cache.Delete(key)
	return nil
}
func (lc *LocalCache) Close() error {
	lc.Cache.Flush()
	log.Println("[INFO] Local cache flushed successfully")
	return nil
}

func convertToInt(i interface{}) (int, error) {
	switch v := i.(type) {
	case int:
		return v, nil
	case int32:
		return int(v), nil
	case int64:
		return int(v), nil
	case uint:
		return int(v), nil
	case uint32:
		return int(v), nil
	case uint64:
		if v > uint64(int(^uint(0)>>1)) { // Check if it exceeds the range of int
			return 0, fmt.Errorf("uint64 value exceeds the range of int: %d", v)
		}
		return int(v), nil
	case float32:
		return int(v), nil
	case float64:
		return int(v), nil
	case string:
		if value, err := strconv.Atoi(v); err == nil {
			return value, nil
		}
		return 0, fmt.Errorf("unable to convert string to int: %s", v)
	case []byte:
		if value, err := strconv.Atoi(string(v)); err == nil {
			return value, nil
		}
		return 0, fmt.Errorf("unable to convert []byte to int: %s", v)
	default:
		return 0, fmt.Errorf("unsupported type: %T", v)
	}
}
