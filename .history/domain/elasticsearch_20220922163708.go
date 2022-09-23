package domain

import (
	"log-ext/adapter/repository"
	"log-ext/domain/dependency"
)

type ealsticsearchService struct {
	elasticDep dependency.ElasticsearchDependency
}

func NewElasticsearch