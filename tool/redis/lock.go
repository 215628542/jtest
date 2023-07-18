package redis

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"time"
)

// Redis分布式写锁
type RedisMutex struct {
	Key      string
	ExpireIn time.Duration
	Ctx      context.Context
	Rdb      *redis.Client
}

var ErrGetLockFail = errors.New("该事件已被占用，请稍后重试")

func NewRedisLocker(ctx context.Context, rdb *redis.Client, key string, expire time.Duration) RedisMutex {
	return RedisMutex{key, expire, ctx, rdb}
}

func (r RedisMutex) Lock() bool {
	return r.Rdb.SetNX(r.Ctx, r.Key, "success", r.ExpireIn).Val()
}

func (r RedisMutex) Unlock() int64 {
	return r.Rdb.Del(r.Ctx, r.Key).Val()
}
