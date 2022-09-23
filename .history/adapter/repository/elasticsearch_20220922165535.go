package repository

import "log-ext/infra"

type Elasticsearch struct {
	infra infra.ElasticsearchInfra
}

	SearchRequest(indexNames []string, query *entity.QueryDocs) (map[string]interface{}, error)
	IndicesDeleteRequest(indexNames []string) ([]byte, error)
