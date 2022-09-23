package repository

import "log-ext/infra"

type Elasticsearch struct {
	infra infra.ElasticsearchInfra
}

func ()SearchRequest(indexNames []string, query *entity.QueryDocs) (map[string]interface{}, error)
func ()IndicesDeleteRequest(indexNames []string) ([]byte, error)
