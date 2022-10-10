package domain

import (
	"log-ext/domain/entity"
)

type SearchService interface {
	NearbyDoc(docid string, num int) ([]byte, error)
	Histogram(query *entity.DateHistogramReq) ([]entity.Buckets, int64, error)
	SearchLogsByFilter(filter *entity.LogsFilter) ([]byte, int, error)
}
