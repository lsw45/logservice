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
		mysql, err := infra.NewMysql(conf.Mysql)
		if err != nil {
			common.Logger.Fatal(err.Error())
		}
		common.Logger.Infof("mysql setting: %+v", mysql.DB.Statement)

		// opensearch
		openDB, err := infra.NewOpensearch(conf.Opensearch)
		if err != nil {
			common.Logger.Fatal(err.Error())
		}
		common.Logger.Infof("opensearch url: %+v", openDB.Transport.(*opensearchtransport.Client).URLs())

		// redis
		redisCli, err := infra.NewRedis(conf.Redis)
		if err != nil {
			common.Logger.Fatal(err.Error())
		}
		common.Logger.Infof("redis client: %+v", redisCli)

		defaultInfra = &RepoInfra{
			Mysql:      mysql,
			Redis:      &infra.Redis{Client: redisCli},
			Opensearch: OP,
			Tunnel:     &infra.Tunnel{Client: conf.TunnelCli},
			Elastic:    elastic,
		}
	})
}
