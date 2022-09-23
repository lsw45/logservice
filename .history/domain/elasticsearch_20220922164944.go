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
	}
for key,sor := range filter.Sort {
	query.Sort = append(query.Sort, elastic)
}
	this.elasticDep.SearchRequest(filter.Indexs, query)
}
