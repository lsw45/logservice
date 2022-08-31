package controller

import (
	"log-ext/domain"
)

type SearchController struct {
	srv *domain.SearchService
}

func NewSearchController(srv *domain.SearchService) *SearchController {
	return &SearchController{
		srv: srv,
	}
}

func (s *SearchController) Search() {

}

func (s *SearchController) List() {
}
