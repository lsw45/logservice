package domain

import (
	"encoding/json"
	"log-ext/common"
	"log-ext/domain/dependency"
	"log-ext/domain/entity"

	"github.com/olivere/elastic/v7"
)

type ealsticsearchService struct {
	elasticDep dependency.ElasticsearchDependency
}

func NewElasticsearchService(dep dependency.ElasticsearchDependency) SearchService {
	return &ealsticsearchService{elasticDep: dep}
}

func (this *ealsticsearchService) Histogram() {

}

func (this *ealsticsearchService) SearchLogsByFilter(filter *entity.LogsFilter) ([]*elastic.SearchHit, int, error) {
	query := &entity.QueryDocs{
		From: filter.Page * filter.PageSize,
		Size: filter.PageSize,
	}

	// elastic:true为升序序，false为降序
	for key, sor := range filter.Sort {
		if sor {
			query.Sort = append(query.Sort, elastic.NewFieldSort(key).Asc())
		} else {
			query.Sort = append(query.Sort, elastic.NewFieldSort(key).Desc())
		}
	}

	hits, err := this.elasticDep.SearchRequest(filter.Indexs, query)
	if err != nil {
		common.Logger.Errorf("search log error: %v", err)
		return nil, 0, err
	}

	json.Marshal(hits)

	return hits, len(hits), nil
}
