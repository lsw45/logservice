package controller

import (
	"log-ext/adapter/repository"
	"log-ext/common"
	"log-ext/domain"
	"log-ext/domain/entity"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 接收通知的回调接口需要确保幂等性，耗时不能超过2秒，如果超时会任务回调失败，会进行重试，最多5次重试
type NotifyController struct {
	notifySrv domain.NotifyService
	mysqlRepo *repository.MysqlRepo
}

func NewNotifyController(mysqlRepo *repository.MysqlRepo, tunnelRepo *repository.TunnelRepo) *NotifyController {
	srv := domain.NewDeployNotifyService(mysqlRepo, tunnelRepo)

	return &NotifyController{
		notifySrv: srv,
		mysqlRepo: mysqlRepo,
	}
}

func (nctl *NotifyController) Notify(c *gin.Context) {
	var message *entity.NotifyDeployMessage
	err := c.ShouldBindJSON(&message)
	if err != nil {
		common.Logger.Errorf("params error: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "validation failed!",
			"error":   err.Error(),
		})
		return
	}

	err = nctl.notifySrv.DeployNotify(message)
	if err != nil {
		common.Logger.Errorf("notify error: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "notify failed!",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"error":   "",
	})
	return
}
