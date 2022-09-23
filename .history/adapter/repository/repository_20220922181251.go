package repository

import (
	"log-ext/common"
	"log-ext/infra"
	"sync"

	"github.com/opensearch-project/opensearch-go/opensearchtransport"
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
		elastic, err := infra.NewElasticsearch(conf.Elasticsearch)
		if err != nil {
			common.Logger.Fatalf("new elasticsearch infra: %v", err)
		}

		// mysql
		DB, err = infra.NewMysqlDB(conf.Mysql)
		if err != nil {
			common.Logger.Fatal(err.Error())
		}
		common.Logger.Infof("mysql setting: %+v", conf.DB.Statement)

		// opensearch
		OpenDB, err = infra.NewOpensearch(conf.Opensearch)
		if err != nil {
			common.Logger.Fatal(err.Error())
		}
		common.Logger.Infof("opensearch url: %+v", conf.OpenDB.Transport.(*opensearchtransport.Client).URLs())

		// redis
		RedisCli, err = infra.NewRedis(conf.Redis)
		if err != nil {
			common.Logger.Fatal(err.Error())
		}
		common.Logger.Infof("redis client: %+v", conf.RedisCli)

		defaultInfra = &RepoInfra{
			Mysql:      &infra.Mysql{DB: mysql},
			Redis:      &infra.Redis{Client: conf.RedisCli},
			Opensearch: &infra.Opensearch{Client: conf.OpenDB},
			Tunnel:     &infra.Tunnel{Client: conf.TunnelCli},
			Elastic:    elastic,
		}
	})
}
