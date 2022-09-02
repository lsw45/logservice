package repository

import (
	"log-ext/common"
	"log-ext/infra"
	"sync"
)

var (
	RepoInfra = &repoInfra{}
	once      = sync.Once{}
)

type repoInfra struct {
	Mysql infra.MysqlInfra
	Opensearch infra.OpensearchInfra
}

func SetRepoInfra(conf *common.AppConfig) {
	once.Do(func() {
		RepoInfra = &repoInfra{
			Mysql:      &infra.Mysql{DB: conf.DB},
			Opensearch: nil,
		}
	})
}
