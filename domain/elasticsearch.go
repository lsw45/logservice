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

func (svc *ealsticsearchService) Histogram() {

}

func (svc *ealsticsearchService) SearchLogsByFilter(filter *entity.LogsFilter) ([]byte, int, error) {
	query := &entity.QueryDocs{
		From: (filter.Page - 1) * filter.PageSize,
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

	// 检验查询语句的json格式是否正确
	var f interface{}
	err := json.Unmarshal([]byte(filter.Keywords), &f)
	if err != nil {
		common.Logger.Errorf("unmarshal json error: %v", err)
		return nil, 0, err
	}
	query.Query = filter.Keywords
	
	if len(filter.Date) > 0 {
		query.StartTime = filter.Date[0]
		query.EndTime = filter.Date[1]
	}

	hits, err := svc.elasticDep.SearchRequest(filter.Indexs, query)
	if err != nil {
		common.Logger.Errorf("search log error: %v", err)
		return nil, 0, err
	}

	re, _ := json.Marshal(hits.Hits)

	return re, int(hits.TotalHits.Value), nil
}
