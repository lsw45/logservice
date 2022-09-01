package controller

import (
	"github.com/gin-gonic/gin"
	"log-ext/adapter/repository"
	"log-ext/common"
)

type AppServer interface {
	RegisterRouter(e *gin.Engine)
}

func NewLogExtServer(conf *common.AppConfig) AppServer {
	repository.SetRepoInfra(conf)

	return &logExtServer{}
}

type logExtServer struct {
	JobController
	SearchController
}

func (ctl *logExtServer) RegisterRouter(e *gin.Engine) {

}
