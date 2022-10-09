package controller

import (
	"fmt"
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

func NewSearchController(esRepo *repository.ElasticsearchRepo) *SearchController {
	search := domain.NewElasticsearchService(esRepo)

	return &SearchController{
		searchSrv: search,
	}
}

func (sctl *SearchController) NearbyDoc(c *gin.Context) {
	docid := c.Query("docid")
	if len(docid) == 0 {
		common.Logger.Errorf("docid empty")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "validation failed!",
			"error":   "docid empty",
		})
		return
	}

	sctl.searchSrv.NearbyDoc()

}

func (sctl *SearchController) Histogram(c *gin.Context) {
	var req *entity.DateHistogramReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		common.Logger.Errorf("params error: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "validation failed!",
			"error":   err.Error(),
		})
		return
	}

	diff := req.EndTime - req.StartTime

	if diff <= 60 {
		req.Interval = fmt.Sprintf("%vs", diff)
	} else if diff > 60 && diff < 3600 {
		req.Interval = fmt.Sprintf("%vm", diff/60)
	} else if diff > 3600 && diff < 24*3600 {
		req.Interval = fmt.Sprintf("%vh", diff/3600)
	} else if diff > 24*3600 {
		req.Interval = fmt.Sprintf("%vd", diff/(24*3600))
	}

	list, total, err := sctl.searchSrv.Histogram(req)
	if err != nil {
		common.Logger.Errorf("params error: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "histograme failed",
			"error":   err.Error(),
		})
		return
	}

	var resp entity.HistogramResp
	resp.CommonResp.Code = 0
	resp.CommonResp.Msg = "success"
	resp.Data = list
	resp.Count = total

	c.JSON(http.StatusOK, resp)
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

	// result, err := json.Marshal(list)
	// if err != nil {
	// 	common.Logger.Errorf("json marshal error: %s", err)
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"message": "system error",
	// 		"error":   err.Error(),
	// 	})
	// 	return
	// }

	var resp entity.LogsFilterResp
	resp.Code = 0
	resp.Msg = "success"
	resp.Data.Results = string(list)
	resp.Data.Count = total

	c.JSON(http.StatusOK, resp)
}
