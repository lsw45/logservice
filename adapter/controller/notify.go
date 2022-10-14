package controller

import (
	"errors"
	"log-ext/adapter/repository"
	"log-ext/common"
	"log-ext/domain"
	"log-ext/domain/entity"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 接收通知的回调接口需要确保幂等性，耗时不能超过2秒，如果超时会任务回调失败，会进行重试，最多5次重试
type DeployController struct {
	notifySrv domain.NotifyService
	mysqlRepo *repository.MysqlRepo
}

func NewNotifyController(mysqlRepo *repository.MysqlRepo, tunnelRepo *repository.TunnelRepo) *DeployController {
	srv := domain.NewDeployNotifyService(mysqlRepo, tunnelRepo)

	return &DeployController{
		notifySrv: srv,
		mysqlRepo: mysqlRepo,
	}
}

func (dctl *DeployController) Notify(c *gin.Context) {
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

	common.Logger.Infof("notify message:%+v\n", message)

	if c.Request.URL.Query().Get("x-token") != "f68192e66ddc4d2a9fd4300bdd4a8f7e" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "validation failed!",
			"error":   errors.New("token invalid"),
		})
		return
	}

	// 处理回调消息
	err = dctl.notifySrv.DeployNotify(message)
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
	})
}
