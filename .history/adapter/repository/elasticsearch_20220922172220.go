package repository

import (
	"errors"
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
		common.Logger.Error(err)
		return nil, err
	}

	if res.TotalHits() == 0 {
		err = errors.Warnf("got SearchResult.TotalHits() = 0")
		return nil, err
	}

	if len(res.Hits.Hits) == 0 {
		err = errors.New("got len(SearchResult.Hits.Hits) = 0")
		return nil, err
	}
	return res.Hits.Hits, nil
}

func (this *Elasticsearch) IndicesDeleteRequest(indexNames []string) ([]byte, error) {
	return nil, nil
}
