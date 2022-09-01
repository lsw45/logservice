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
	infra.MysqlInfra
	infra.OpensearchInfra
}

func SetRepoInfra(conf *common.AppConfig) {
	once.Do(func() {
		RepoInfra = &repoInfra{
			MysqlInfra:      infra.MysqlInfra{conf.DB},
			OpensearchInfra: nil,
		}
	})
}
