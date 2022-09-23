package repository

import (
	"log-ext/domain/entity"
	"log-ext/infra"
)

type Elasticsearch struct {
	infra infra.ElasticsearchInfra
}

func (this *Elasticsearch) SearchRequest(indexNames []string, query *entity.QueryDocs) (map[string]interface{}, error) {
this.infra.SearchRequest()
}
func (this *Elasticsearch) IndicesDeleteRequest(indexNames []string) ([]byte, error) {

}
