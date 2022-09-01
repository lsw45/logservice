package controller

import (
	"log-ext/adapter/repository"
	"log-ext/domain"

	"github.com/gin-gonic/gin"
)

type SearchController struct {
	search domain.SearchService
}

func NewSearchController() *SearchController {
	search := domain.NewSearchLogService(&repository.OpensearchRepo{})

	return &SearchController{
		search: search,
	}
}

func (s *SearchController) Search() {

}

func (s *SearchController) List(c *gin.Context) {
	s.search.List()
}
