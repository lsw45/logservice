package domain

import (
	"log-ext/domain/entity"

	"github.com/olivere/elastic/v7"
)

type SearchService interface {
	Aggregation(req entity.AggregationReq) (*elastic.SearchResult, error)
	NearbyDoc(indexName string, times int64, num int) ([]*elastic.SearchHit, error)
	Histogram(query *entity.DateHistogramReq) ([]entity.Buckets, int64, error)
	SearchLogsByFilter(filter *entity.LogsFilter) (*elastic.SearchHits, int, error)
}
