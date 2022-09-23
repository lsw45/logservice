package infra

import "github.com/opensearch-project/opensearch-go/opensearchapi"
"github.com/elastic/elasticsearch"
type ElasticsearchInfra interface {
	SearchRequest(indexNames []string, content string) (*opensearchapi.Response, error)
	IndicesDeleteRequest(indexNames []string) (*opensearchapi.Response, error)
}

type Elasticsearch struct {
}
