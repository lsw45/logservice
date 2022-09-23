package repository

import (
	"log-ext/common"
	"log-ext/domain/entity"
	"log-ext/infra"
)

type Elasticsearch struct {
	infra infra.ElasticsearchInfra
}

func (this *Elasticsearch) SearchRequest(indexNames []string, query *entity.QueryDocs) (map[string]interface{}, error) {
	hits, err := this.infra.SearchRequest(indexNames, query)
	if err != nil {
		common.Logger.Error(err)
		return nil, err
	}


	got := string(hit.Source)
}
func (this *Elasticsearch) IndicesDeleteRequest(indexNames []string) ([]byte, error) {

}
