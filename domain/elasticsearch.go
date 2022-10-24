package domain

import (
	"errors"
	"fmt"
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

func (svc *ealsticsearchService) Histogram(filter *entity.LogsFilter) ([]entity.BucketsList, int64, error) {
	query := &entity.DateHistogramReq{
		Query:     filter.Keywords,
		StartTime: filter.StartTime,
		EndTime:   filter.EndTime,
	}

	if filter.EnvID > 0 && filter.ProjectId > 0 && len(filter.RegionID) > 0 {
		for _, id := range filter.RegionID {
			query.Indexs = append(query.Indexs, fmt.Sprintf("server-%v-%v-%v", filter.ProjectId, filter.EnvID, id))
		}
	}

	var err error
	query.Indexs, err = svc.elasticDep.IndexExists(query.Indexs)
	if err != nil {
		common.Logger.Errorf("search index error: %v", err)
		return nil, 0, err
	}

	if len(query.Indexs) == 0 {
		return nil, 0, nil
	}

	query.Interval = (query.EndTime - query.StartTime) / 60 // 将时段60等分

	list, total, err := svc.elasticDep.Histogram(query)
	if err != nil {
		common.Logger.Errorf("histogram error: %v", err)
		return nil, 0, err
	}

	data := []entity.BucketsList{}
	for _, v := range list {
		t := v.Key.(float64)
		data = append(data, entity.BucketsList{DocCount: v.DocCount, StartTime: t, EndTime: t + float64(query.Interval)})
	}

	return data, total, nil
}

func (svc *ealsticsearchService) SearchLogsByFilter(filter *entity.LogsFilter) (*elastic.SearchHits, int, error) {
	if filter.StartTime > filter.EndTime || filter.StartTime == 0 {
		common.Logger.Error("time error")
		return nil, 0, errors.New("param error")
	}

	query, err := transQuerydoc(filter)
	if err != nil {
		return nil, 0, err
	}

	query.Indexs, err = svc.elasticDep.IndexExists(query.Indexs)
	if err != nil {
		common.Logger.Errorf("search index error: %v", err)
		return nil, 0, err
	}

	if len(query.Indexs) == 0 {
		return nil, 0, nil
	}

	hits, err := svc.elasticDep.SearchRequest(query.Indexs, query)
	if err != nil {
		common.Logger.Errorf("search log error: %v", err)
		return nil, 0, err
	}

	if hits == nil {
		return nil, 0, nil
	}
	return hits, int(hits.TotalHits.Value), nil
}

func transQuerydoc(filter *entity.LogsFilter) (*entity.QueryDocs, error) {
	query := &entity.QueryDocs{
		StartTime: filter.StartTime,
		EndTime:   filter.EndTime,
		Query:     filter.Keywords,
	}

	if filter.Page > 0 {
		query.From = (filter.Page - 1) * filter.PageSize
		query.Size = filter.PageSize
	}

	if filter.EnvID > 0 && filter.ProjectId > 0 && len(filter.RegionID) > 0 {
		for _, id := range filter.RegionID {
			query.Indexs = append(query.Indexs, fmt.Sprintf("server-%v-%v-%v", filter.ProjectId, filter.EnvID, id))
		}
	}

	// elastic:true为升序，false为降序
	for key, sor := range filter.Sort {
		if sor {
			query.Sort = append(query.Sort, elastic.NewFieldSort(key).Asc())
		} else {
			query.Sort = append(query.Sort, elastic.NewFieldSort(key).Desc())
		}
	}

	// 默认排序字段
	if len(query.Sort) == 0 {
		query.Sort = []elastic.Sorter{elastic.NewFieldSort(entity.LogSortField).Desc()}
	}

	return query, nil
}
