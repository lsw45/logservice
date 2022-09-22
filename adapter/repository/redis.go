package repository

import (
	"context"
	"log-ext/common/errorx"
	"log-ext/domain/dependency"
	"log-ext/infra"
)

var _ dependency.RedisRepo = (*Redis)(nil)

type Redis struct {
	infra.RedisInfra
}

func NewRedis() *Redis {
	return &Redis{defaultRepo.Redis}
}

func (r *Redis) Get(ctx context.Context, key string) (string, *errorx.CodeError) {
	val, err := r.RedisInfra.Get(ctx, key)

	if err != errorx.ErrNil {
		return "", errorx.NewErrCode(err)
	}

	return val, nil
}
