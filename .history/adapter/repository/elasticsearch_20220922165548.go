package repository

import "log-ext/infra"

type Elasticsearch struct {
	infra infra.ElasticsearchInfra
}

func (this *Elasticsearch)SearchRequest(indexNames []string, query *entity.QueryDocs) (map[string]interface{}, error)
func (this *Elasticsearch)IndicesDeleteRequest(indexNames []string) ([]byte, error)
