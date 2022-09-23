package infra

import "github.com/elastic/elasticsearch"
type ElasticsearchInfra interface {
	SearchRequest(indexNames []string, content string) (*elasticsearch.Response, error)
	IndicesDeleteRequest(indexNames []string) (*opensearchapi.Response, error)
}

type Elasticsearch struct {
}
