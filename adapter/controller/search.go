package controller

import (
	"log-ext/adapter/repository"
	"log-ext/common"
	"log-ext/domain"
	"log-ext/domain/entity"
	"net/http"
	"strconv"

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

func (sctl *SearchController) Aggregation(c *gin.Context) {
	var req *entity.AggregationReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		common.Logger.Errorf("param empty")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "validation failed!",
			"error":   "param empty",
		})
		return
	}

	result, err := sctl.searchSrv.Aggregation(*req)
	if err != nil {
		common.Logger.Errorf("controller search error: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "search log failed!",
			"error":   err.Error(),
		})
		return
	}

	var resp entity.AggregationResp
	resp.Code = 0
	resp.Msg = "success"
	resp.Data = result
	c.JSON(http.StatusOK, resp)
}

func (sctl *SearchController) NearbyDoc(c *gin.Context) {
	index := c.Param("index")
	time := c.Param("time")
	num := c.Param("num")

	if len(index) == 0 || len(num) == 0 {
		common.Logger.Errorf("params empty")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "validation failed!",
			"error":   "params empty",
		})
		return
	}

	times, err := strconv.ParseInt(time, 10, 0)
	if err != nil {
		common.Logger.Errorf("system error: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "search log failed!",
			"error":   err.Error(),
		})
		return
	}

	nums, err := strconv.Atoi(num)
	if err != nil {
		common.Logger.Errorf("system error: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "search log failed!",
			"error":   err.Error(),
		})
		return
	}

	list, err := sctl.searchSrv.NearbyDoc(index, times, nums)
	if err != nil {
		common.Logger.Errorf("controller search error: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "search log failed!",
			"error":   err.Error(),
		})
		return
	}

	var resp entity.NearbyDocResp
	resp.Code = 0
	resp.Msg = "success"
	resp.Data = list

	c.JSON(http.StatusOK, resp)
}

func (sctl *SearchController) Histogram(c *gin.Context) {
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

	if len(filter.Indexs) == 0 || filter.StartTime > filter.EndTime || filter.StartTime == 0 {
		common.Logger.Error("params error")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "validation failed!",
			"error":   "params error",
		})
		return
	}

	histoReq := &entity.DateHistogramReq{
		StartTime: filter.StartTime,
		EndTime:   filter.EndTime,
		Indexs:    filter.Indexs,
	}

	histoReq.Interval = (histoReq.EndTime - histoReq.StartTime) / 60 // 将时段60等分

	list, total, err := sctl.searchSrv.Histogram(histoReq)
	if err != nil {
		common.Logger.Errorf("params error: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "histograme failed",
			"error":   err.Error(),
		})
		return
	}

	data := []entity.BucketsList{}
	for _, v := range list {
		t := v.Key.(float64)
		data = append(data, entity.BucketsList{DocCount: v.DocCount, StartTime: t, EndTime: t + float64(histoReq.Interval)})
	}

	var resp entity.HistogramResp
	resp.CommonResp.Code = 0
	resp.CommonResp.Msg = "success"
	resp.Data = data
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

	var resp entity.LogsFilterResp
	resp.Code = 0
	resp.Msg = "success"
	resp.Data.Results = list
	resp.Data.Count = int64(total)

	c.JSON(http.StatusOK, resp)
}
