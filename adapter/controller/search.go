package controller

import (
	"fmt"
	"log-ext/adapter/repository"
	"log-ext/common"
	"log-ext/domain"
	"log-ext/domain/entity"
	"net/http"
	"strconv"
	"time"

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

	buckets, err := sctl.searchSrv.Aggregation(req.Indexs, req.Aggs, req.AggsName)
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
	resp.Data.Results = string(buckets)

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

	s, err := time.Parse("2006-01-02 15:04:05", filter.Date[0])
	if err != nil {
		common.Logger.Errorf("time parse error: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "wrong start date format",
			"error":   err.Error(),
		})
		return
	}

	e, err := time.Parse("2006-01-02 15:04:05", filter.Date[1])
	if err != nil {
		common.Logger.Errorf("time parse error: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "wrong end date format",
			"error":   err.Error(),
		})
		return
	}

	histoReq := &entity.DateHistogramReq{
		StartTime: s.Unix(),
		EndTime:   e.Unix(),
	}

	interval := (histoReq.EndTime - histoReq.StartTime) / 60

	if interval <= 60 {
		histoReq.Interval = fmt.Sprintf("%vs", interval)
	} else if interval > 60 && interval < 3600 {
		histoReq.Interval = fmt.Sprintf("%vm", interval/60)
	} else if interval > 3600 && interval < 24*3600 {
		histoReq.Interval = fmt.Sprintf("%vh", interval/3600)
	} else if interval > 24*3600 {
		histoReq.Interval = fmt.Sprintf("%vd", interval/(24*3600))
	}

	list, total, err := sctl.searchSrv.Histogram(histoReq)
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
	resp.Data.Count = int64(total)

	c.JSON(http.StatusOK, resp)
}
