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
	Tunnel     infra.TunnelInfra
	Elastic    infra.TunnelInfra
}

func SetRepoInfra(conf *common.AppConfig) {
	once.Do(func() {
		elastic,err := infra.NewElasticsearch(conf.Elasticsearch)
		if err != nil {
			
		}
		defaultInfra = &RepoInfra{
			Mysql:      &infra.Mysql{DB: conf.DB},
			Redis:      &infra.Redis{Client: conf.RedisCli},
			Opensearch: &infra.Opensearch{Client: conf.OpenDB},
			Tunnel:     &infra.Tunnel{Client: conf.TunnelCli},
			Elastic:     infra.NewElasticsearch(conf.ElasticCli),
		}
	})
}
