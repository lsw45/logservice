package domain

import (
	"log-ext/domain/dependency"
	"log-ext/domain/entity"
)

type ealsticsearchService struct {
	elasticDep dependency.ElasticsearchDependency
}

func NewElasticsearchService(dep dependency.ElasticsearchDependency) SearchService {
	return &ealsticsearchService{elasticDep: dep}
}

func (this *ealsticsearchService) Histogram() {

}
func (this *ealsticsearchService) SearchLogsByFilter(filter *entity.LogsFilter) ([]interface{}, int, error) {
	query := &entity.QueryDocs{
		From: filter.Page*filter.PageSize,
		Size: filter.PageSize,
		type SortBy []Type
		
		func (a SortBy) Len() int           { return len(a) }
		func (a SortBy) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
		func (a SortBy) Less(i, j int) bool { return a[i] < a[j] }
	}
	this.elasticDep.SearchRequest(filter.Indexs, query)
}
