package redis

import (
	"fmt"
	"github.com/cy18cn/zlog"
	"github.com/go-redis/redis"
)

type RedisDB struct {
	db *redis.Client
}

func NewRedisDB(opts *redis.Options) (*RedisDB, error) {
	redisAddr := opts.Addr
	if redisAddr == "" {
		zlog.Error("redis addr is not set")
		return nil, fmt.Errorf("redis addr is not set")
	}

	db := redis.NewClient(opts)
	if pong, err := db.Ping().Result(); err != nil || pong != "PONG" {
		zlog.Errorf("could not connect to the redis service, err: %v", err)
		return nil, err
	}

	return &RedisDB{db}, nil
}

func (self *RedisDB) DB() *redis.Client {
	return self.db
}

func (self *RedisDB) Close() error {
	if self.db != nil {
		return self.db.Close()
	}

	return nil
}
