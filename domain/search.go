package domain

import "log-ext/adapter/repository"

type SearchInterface interface {
}

type SearchService struct {
	repo repository.SearchInterface
}
