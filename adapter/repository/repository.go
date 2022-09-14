package repository

import (
	"log-ext/common"
	"log-ext/infra"
	"sync"
)

var (
	defaultInfra = &RepoInfra{}
	once         = sync.Once{}
)

type RepoInfra struct {
	Redis      infra.RedisInfra
	Mysql      infra.MysqlInfra
	Opensearch infra.OpensearchInfra
	TunnelRepo infra.TunnelInfra
}

func SetRepoInfra(conf *common.AppConfig) {
	once.Do(func() {
		defaultInfra = &RepoInfra{
			Redis:      &infra.Redis{Client: conf.RedisCli},
			Mysql:      &infra.Mysql{DB: conf.DB},
			Opensearch: &infra.Opensearch{Client: conf.OpenDB},
		}
	})
}
