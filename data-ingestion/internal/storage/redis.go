package storage

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Conn *redis.Client
}

func NewRedis(addr string) *RedisClient {
	return &RedisClient{
		Conn: redis.NewClient(&redis.Options{Addr: addr}),
	}
}

func (r *RedisClient) MarkCode(ctx context.Context, pipe redis.Pipeliner, code string, bitOffset int) {
	// u1: 1-bit unsigned int at the specific file's offset
	pipe.BitField(ctx, "promo:"+code, "SET", "u1", bitOffset, 1)
}