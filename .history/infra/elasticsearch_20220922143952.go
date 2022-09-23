package infra

type ElasticsearchInfra interface {
SearchRequest(indexNames []string, content string) (*opensearchapi.Response, error)
	IndicesDeleteRequest(indexNames []string) (*opensearchapi.Response, error)
}
}

type Elasticsearch struct {

}