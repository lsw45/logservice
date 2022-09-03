package infra

import (
	"github.com/go-redis/redis"
	"log-ext/common"
)

type RedisInfra interface {
	Get()
}

type Redis struct {
	redis.Client
}

func NewRedis(conf common.Redis) (*redis.Client, error) {

	return nil, nil
}

func (r *Redis) Get() {

}
