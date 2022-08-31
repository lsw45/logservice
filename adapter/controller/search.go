package controller

import (
	"log-ext/domain"

	"github.com/gin-gonic/gin"
)

type SearchController struct {
	srv *domain.SearchInterface
}

func NewSearchController(srv domain.SearchInterface) *SearchController {
	return &SearchController{
		srv: &srv,
	}
}

func (s *SearchController) Search() {

}

func (s *SearchController) List(c *gin.Context) {
}
