package infra

import (
	"context"
	"log-ext/common"
	"log-ext/common/errorx"

	red "github.com/go-redis/redis/v8"
)

type RedisInfra interface {
	Get(ctx context.Context, key string) (string, errorx.ErrInt)
}

type Redis struct {
	Client *red.Client
}

func NewRedis(conf common.Redis) (*Redis, error) {
	rdb := red.NewClient(&red.Options{
		DB:           conf.DB,
		Addr:         conf.Addr,
		Password:     conf.Password,
		MaxRetries:   conf.MaxRetries,
		MinIdleConns: conf.MinIdleConns,
	})
	return &Redis{rdb}, nil
}

func (r *Redis) Get(ctx context.Context, key string) (string, errorx.ErrInt) {
	val, err := r.Client.Get(ctx, key).Result()

	if err == red.Nil {
		return "", errorx.AUTH_ERROR
	}

	if err != nil {
		common.Logger.Errorf("redis get error: %v", err)
		return "", errorx.SERVER_COMMON_ERROR
	}

	return val, errorx.ErrNil
}
