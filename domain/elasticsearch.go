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

func (svc *ealsticsearchService) NearbyDoc() {

}

func (svc *ealsticsearchService) Histogram(query *entity.DateHistogramReq) ([]entity.Buckets, int64, error) {
	histogram, total, err := svc.elasticDep.Histogram(query)
	if err != nil {
		common.Logger.Errorf("histogram error: %v", err)
		return nil, 0, err
	}
	return histogram, total, nil
}

func (svc *ealsticsearchService) SearchLogsByFilter(filter *entity.LogsFilter) ([]byte, int, error) {
	query, err := transQuerydoc(filter)
	if err != nil {
		return nil, 0, err
	}

	hits, err := svc.elasticDep.SearchRequest(filter.Indexs, query)
	if err != nil {
		common.Logger.Errorf("search log error: %v", err)
		return nil, 0, err
	}

	re, _ := json.Marshal(hits.Hits)

	return re, int(hits.TotalHits.Value), nil
}

func transQuerydoc(filter *entity.LogsFilter) (*entity.QueryDocs, error) {
	query := &entity.QueryDocs{
		From: (filter.Page - 1) * filter.PageSize,
		Size: filter.PageSize,
	}

	// elastic:true为升序，false为降序
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
		return nil, err
	}
	query.Query = filter.Keywords

	if len(filter.Date) > 1 {
		query.StartTime = filter.Date[0]
		query.EndTime = filter.Date[1]
	}
	return query, nil
}
