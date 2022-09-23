package domain

import "log-ext/adapter/repository"

type ealsticsearchService struct {
	elasticRepo repository.Elasticsearch
}

