package domain

import (
	"errors"
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

func (svc *ealsticsearchService) Aggregation(req entity.AggregationReq) (*elastic.SearchResult, error) {
	return svc.elasticDep.Aggregation(req)
}

func (svc *ealsticsearchService) NearbyDoc(indexName string, times int64, num int) ([]*elastic.SearchHit, error) {
	return svc.elasticDep.NearbyDoc(indexName, times, num)
}

func (svc *ealsticsearchService) Histogram(query *entity.DateHistogramReq) ([]entity.Buckets, int64, error) {
	if query.StartTime > query.EndTime || query.StartTime == 0 {
		common.Logger.Error("time error")
		return nil, 0, errors.New("param error")
	}

	if query.EnvID == 0 || query.ProjectId == 0 || query.RegionID == 0 {
		common.Logger.Error("index error")
		return nil, 0, errors.New("param error")
	}

	histogram, total, err := svc.elasticDep.Histogram(query)
	if err != nil {
		common.Logger.Errorf("histogram error: %v", err)
		return nil, 0, err
	}
	return histogram, total, nil
}

func (svc *ealsticsearchService) SearchLogsByFilter(filter *entity.LogsFilter) (*elastic.SearchHits, int, error) {
	if filter.StartTime > filter.EndTime || filter.StartTime == 0 {
		common.Logger.Error("time error")
		return nil, 0, errors.New("param error")
	}

	if filter.EnvID == 0 || filter.ProjectId == 0 || filter.RegionID == 0 {
		common.Logger.Error("index error")
		return nil, 0, errors.New("param error")
	}

	query, err := transQuerydoc(filter)
	if err != nil {
		return nil, 0, err
	}

	hits, err := svc.elasticDep.SearchRequest(filter.Indexs, query)
	if err != nil {
		common.Logger.Errorf("search log error: %v", err)
		return nil, 0, err
	}

	return hits, int(hits.TotalHits.Value), nil
}

func transQuerydoc(filter *entity.LogsFilter) (*entity.QueryDocs, error) {
	query := &entity.QueryDocs{
		From:      (filter.Page - 1) * filter.PageSize,
		Size:      filter.PageSize,
		StartTime: filter.StartTime,
		EndTime:   filter.EndTime,
		Query:     filter.Keywords,
	}

	// elastic:true为升序，false为降序
	for key, sor := range filter.Sort {
		if sor {
			query.Sort = append(query.Sort, elastic.NewFieldSort(key).Asc())
		} else {
			query.Sort = append(query.Sort, elastic.NewFieldSort(key).Desc())
		}
	}

	return query, nil
}
