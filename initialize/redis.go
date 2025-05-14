package initialize

import (
	"context"
	"goAccounting/global/constant"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type _redis struct {
	Addr     string `yaml:"Addr"`
	Password string `yaml:"Password"`
	Db       int    `yaml:"Db"`
}

type RedisHook struct {
	name string
}

func (r *_redis) getNewRedisClient(name string, dbNum int) (*redis.Client, error) {
	// connect to redis
	connect := func() (*redis.Client, error) {
		db := redis.NewClient(&redis.Options{Addr: r.Addr, Password: r.Password, DB: r.Db})
		ctx, cancel := context.WithTimeout(context.TODO(), time.Second*3)
		defer cancel()
		return db, db.Ping(ctx).Err()
	}

	db, err := reconnection[*redis.Client](connect, 3)
	if err != nil {
		return db, err
	}
	if Config.Mode == constant.Debug {
		db.AddHook(&RedisHook{name: name})
	}
	return db, err
}

func (r *_redis) initializeRedis() error {
	if len(r.Addr) == 0 {
		log.Println("Redis initialization skipped: no address provided")
		return nil
	}
	var err error
	Rdb, err = r.getNewRedisClient("", r.Db)
	if err != nil {
		return err
	}

}
