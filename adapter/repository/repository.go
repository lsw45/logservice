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
	Redis   infra.RedisInfra
	Mysql   infra.MysqlInfra
	Tunnel  infra.TunnelInfra
	Elastic infra.ElasticsearchInfra
}

func SetRepoInfra(redis infra.RedisInfra, mysql infra.MysqlInfra, tunnel infra.TunnelInfra, elastic infra.ElasticsearchInfra) {
	once.Do(func() {
		defaultRepo = &RepoInfra{
			Mysql:   mysql,
			Redis:   redis,
			Tunnel:  tunnel,
			Elastic: elastic,
		}
	})
}

func Close() {
	defaultRepo.Mysql.Close()
}
