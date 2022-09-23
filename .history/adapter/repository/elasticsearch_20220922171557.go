package repository

import (
	"errors"
	"log-ext/common"
	"log-ext/domain/entity"
	"log-ext/infra"
)

type Elasticsearch struct {
	infra infra.ElasticsearchInfra
}

func (this *Elasticsearch) SearchRequest(indexNames []string, query *entity.QueryDocs) (map[string]interface{}, error) {
	res, err := this.infra.SearchRequest(indexNames, query)
	if err != nil {
		common.Logger.Error(err)
		return nil, err
	}
	if got, want := res.TotalHits(), int64(1); got < want {
		err = errors.New("got SearchResult.TotalHits() = 0")
		return nil, err
	}
	
	if got, want := len(res.Hits.Hits), 1; got < want {
		err = errors.New("got len(SearchResult.Hits.Hits) = 0")
		return nil, err
	}

}
func (this *Elasticsearch) IndicesDeleteRequest(indexNames []string) ([]byte, error) {

}
