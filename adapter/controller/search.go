package controller

import (
	"log-ext/adapter/repository"
)

type SearchService struct {
}

func (s *SearchService) Search() {
	repo := repository.SearchService{}
	repo.Search()
}

func (s *SearchService) Roll() {

}
