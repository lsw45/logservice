package domain

import (
	"log-ext/adapter/repository"
	"log-ext/domain/dependency"

	"github.com/olivere/elastic/v7"
)

type ealsticsearchService struct {
	elasticDep dependency.ElasticsearchDependency
}

func NewElasticsearchService(dep dependency.ElasticsearchDependency) elasticsearchservice{

}