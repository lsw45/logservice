package infra

import (
	"log-ext/common"

	"github.com/olivere/elastic/v7"
)

type ElasticsearchInfra interface {
	SearchRequest(indexNames []string, content string) (*elastic.Response, error)
	IndicesDeleteRequest(indexNames []string) (*elastic.Response, error)
}

type Elasticsearch struct {
	*elastic.Client
}

func NewElasticsearch(conf common.AppConfig) (*elastic.Client, error) {
	client, err := elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetURL(conf.Opensearch.Address...),
		elastic.SetBasicAuth(conf.Opensearch.Username, conf.Opensearch.Password),
	)

	if err != nil {
		return nil,err
	}
}
