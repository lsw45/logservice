package domain

import (
	"log-ext/domain/dependency"
)

type ealsticsearchService struct {
	elasticDep dependency.ElasticsearchDependency
}

func NewElasticsearchService(dep dependency.ElasticsearchDependency) SearchService {
	return &ealsticsearchService{elasticDep: dep}
}


func Histogram()
func SearchLogsByFilter(filter *entity.LogsFilter) ([]interface{}, int, error)
