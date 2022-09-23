package infra

import "https://github.com/olivere/elastic"

type ElasticsearchInfra interface {
	SearchRequest(indexNames []string, content string) (*elasticsearch.Response, error)
	IndicesDeleteRequest(indexNames []string) (*elasticsearch.Response, error)
}

type Elasticsearch struct {
}
