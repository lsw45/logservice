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


func (this *)Histogram()
func (this *)SearchLogsByFilter(filter *entity.LogsFilter) ([]interface{}, int, error)
