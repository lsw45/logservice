package controller

import (
	"encoding/json"
	"log-ext/adapter/repository"
	"log-ext/common"
	"log-ext/domain"
	"log-ext/domain/entity"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SearchController struct {
	searchSrv domain.SearchService
}

func NewSearchController(openRepo *repository.OpensearchRepo) *SearchController {
	search := domain.NewSearchLogService(openRepo)

	return &SearchController{
		searchSrv: search,
	}
}

func (sctl *SearchController) Histogram(c *gin.Context) {

}

func (sctl *SearchController) SearchLogsByFilter(c *gin.Context) {
	var filter *entity.LogsFilterReq
	err := c.ShouldBindJSON(&filter)
	if err != nil {
		common.Logger.Errorf("params error: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "validation failed!",
			"error":   err.Error(),
		})
		return
	}

	list, total, err := sctl.searchSrv.SearchLogsByFilter(&entity.LogsFilter{LogsFilterReq: *filter})
	if err != nil {
		common.Logger.Errorf("controller search error: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "search log failed!",
			"error":   err.Error(),
		})
		return
	}

	result, err := json.Marshal(list)
	if err != nil {
		common.Logger.Errorf("json marshal error: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "system error",
			"error":   err.Error(),
		})
		return
	}

	var resp entity.LogsFilterResp
	resp.Code = 0
	resp.Msg = "success"
	resp.Data.Results = string(result)
	resp.Data.Count = len(list)
	resp.Data.Total = total

	c.JSON(http.StatusOK, resp)
}
