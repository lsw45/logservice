package repository

import (
	"log-ext/infra"
	"sync"
)

var (
	defaultRepo = &RepoInfra{}
	once        = sync.Once{}
)

type RepoInfra struct {
	Redis      infra.RedisInfra
	Mysql      infra.MysqlInfra
	Opensearch infra.OpensearchInfra
	Tunnel     infra.TunnelInfra
	Elastic    infra.ElasticsearchInfra
}

func SetRepoInfra(redis infra.RedisInfra, mysql infra.MysqlInfra, tunnel infra.TunnelInfra, elastic infra.ElasticsearchInfra) {
	once.Do(func() {
		defaultRepo = &RepoInfra{
			Mysql:      mysql,
			Redis:      redis,
			Opensearch: openDB,
			Tunnel:     tunnel,
			Elastic:    elastic,
		}
	})
}

func Close() {
	defaultRepo.Mysql.Close()
}
