package repository

import (
	"log-ext/common"
	"log-ext/domain/entity"
	"log-ext/infra"

	"github.com/olivere/elastic/v7"
)

type Elasticsearch struct {
	infra infra.ElasticsearchInfra
}

func (this *Elasticsearch) SearchRequest(indexNames []string, query *entity.QueryDocs) ([]*elastic.SearchHit, error) {
	res, err := this.infra.SearchRequest(indexNames, query)
	if err != nil {
		return nil, err
	}

	if res.TotalHits() == 0 {
		common.Logger.Warn("got SearchResult.TotalHits() = 0")
	}

	if len(res.Hits.Hits) == 0 {
		common.Logger.Warn("got len(SearchResult.Hits.Hits) = 0")
	}

	return res.Hits.Hits, nil
}

func (this *Elasticsearch) IndicesDeleteRequest(indexNames []string) ([]byte, error) {
	return nil, nil
}
