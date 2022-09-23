package infra

import (
	"context"
	"log-ext/common"

	"github.com/olivere/elastic/v7"
)

type ElasticsearchInfra interface {
	SearchRequest(indexNames []string, content string) (*elastic.Response, error)
	IndicesDeleteRequest(indexNames []string) (*elastic.Response, error)
}

type elasticsearch struct {
	Client *elastic.Client
}

func NewElasticsearch(conf common.AppConfig) (*elasticsearch, error) {
	client, err := elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetURL(conf.Elasticsearch.Address...),
		elastic.SetBasicAuth(conf.Elasticsearch.Username, conf.Elasticsearch.Password),
	)

	if err != nil {
		return nil, err
	}
	return &elasticsearch{client}, nil
}

type QueryDocs struct {
	From      int
	Size      int
	enableDSL bool
	sort      map[string]bool
}

func (es *elasticsearch) SearchRequest(indexNames []string, QueryDocs) (*elastic.Response, error) {
	query := elastic.NewMatchAllQuery()
	res, err := es.Client.Search().Index(indexNames...).From().Size().SortBy().Query(query).Do(context.Background())

	if err != nil {
		return nil, err
	}

	return nil, nil
}
