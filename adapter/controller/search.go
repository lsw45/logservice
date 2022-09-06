package controller

import (
	"log-ext/adapter/repository"
	"log-ext/common"
	"log-ext/domain"
	"log-ext/domain/entity"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SearchController struct {
	searchSrv domain.SearchService
	openRepo  *repository.OpensearchRepo
}

func NewSearchController(openRepo *repository.OpensearchRepo) *SearchController {
	search := domain.NewSearchLogService(openRepo)

	return &SearchController{
		searchSrv: search,
		openRepo:  openRepo,
	}
}

func (sctl *SearchController) Histogram(c *gin.Context) {

}

func (sctl *SearchController) SearchLogsByFilter(c *gin.Context) {
	var filter *entity.LogsFilterReq
	err := c.ShouldBindJSON(&filter)
	if err != nil {
		common.Logger.Errorf("%s", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "validation failed!",
			"error":   err.Error(),
		})
	}

	result, err := sctl.searchSrv.SearchLogsByFilter(&entity.LogsFilter{LogsFilterReq: *filter})
	if err != nil {
		common.Logger.Errorf("%s", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "search log failed!",
			"error":   err.Error(),
		})
	}

	var resp entity.LogsFilterResp
	resp.Code = 0
	resp.Msg = "success"
	resp.Data.Results = result
	resp.Data.Count = len(result)

	c.JSON(http.StatusOK, resp)
}
