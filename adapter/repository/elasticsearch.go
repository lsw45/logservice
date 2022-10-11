package repository

import (
	"log-ext/common"
	"log-ext/domain/dependency"
	"log-ext/domain/entity"
	"log-ext/infra"

	"github.com/olivere/elastic/v7"
)

var _ dependency.ElasticsearchDependency = new(ElasticsearchRepo)

type ElasticsearchRepo struct {
	infra.ElasticsearchInfra
}

func NewElasticsearchRepo() *ElasticsearchRepo {
	return &ElasticsearchRepo{defaultRepo.Elastic}
}

func (elastic *ElasticsearchRepo) SearchRequest(indexNames []string, query *entity.QueryDocs) (*elastic.SearchHits, error) {
	res, err := elastic.ElasticsearchInfra.SearchRequest(indexNames, query)
	if err != nil {
		return nil, err
	}

	if res.TotalHits() == 0 {
		common.Logger.Warn("got SearchResult.Hits = 0")
		return res.Hits, nil
	}

	if len(res.Hits.Hits) == 0 {
		common.Logger.Warn("got SearchResult.Hits.Hits = 0")
	}

	return res.Hits, nil
}

func (elastic *ElasticsearchRepo) IndicesDeleteRequest(indexNames []string) ([]byte, error) {
	return nil, nil
}

func (elastic *ElasticsearchRepo) Histogram(query *entity.DateHistogramReq) ([]entity.Buckets, int64, error) {
	res, total, err := elastic.ElasticsearchInfra.Histogram(query)
	if err != nil {
		return nil, 0, err
	}

	return res, total, nil
}

func (elastic *ElasticsearchRepo) NearbyDoc(indexName string, times int64, num int) ([]*elastic.SearchHit, error) {
	return elastic.ElasticsearchInfra.NearbyDoc(indexName, times, num)
}

func (elastic *ElasticsearchRepo) Aggregation(req entity.AggregationReq) (*elastic.SearchResult, error) {
	return elastic.ElasticsearchInfra.Aggregation(req)
}
