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
	Elastic    infra.ElasticsearchInfra
}

func SetRepoInfra(conf *common.AppConfig) {
	once.Do(func() {
		defaultInfra = &RepoInfra{
			Mysql:      mysql,
			Redis:      redis,
			Opensearch: openDB,
			Tunnel:     tunnel,
			Elastic:    elastic,
		}
	})
}
