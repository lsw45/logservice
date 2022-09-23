package infra

import "github.com/olivere/elastic/v7"

type ElasticsearchInfra interface {
	SearchRequest(indexNames []string, content string) (*elastic.Response, error)
	IndicesDeleteRequest(indexNames []string) (*elasticsearch.Response, error)
}

type Elasticsearch struct {
}
