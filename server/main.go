package main

import (
	"log-ext/adapter/controller"
	"log-ext/common"
	"log-ext/infra"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	common.Logger.Info("Service Starting")
	conf := common.NewAppConfig()

	common.Logger.Info("Start Http")
	logExternal := NewServer(conf)
	logExternal.Start()
}

type LogExternalServer struct {
	conf   *common.AppConfig
	server *http.Server
}

func NewServer(conf *common.AppConfig) *LogExternalServer {
	engine := gin.New()
	gin.SetMode(conf.Server.RunMode)

	server := &LogExternalServer{
		conf: conf,
		server: &http.Server{
			Addr:           ":" + strconv.Itoa(conf.Server.HTTPPort),
			Handler:        engine,
			TLSConfig:      nil,
			ReadTimeout:    conf.Server.ReadTimeOut * time.Second,
			WriteTimeout:   conf.Server.WriteTimeOut * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
	}
	server.InitClient()

	// 创建控制器
	ctr := controller.NewLogExtServer(conf)
	// 注册路由
	ctr.RegisterRouter(engine)

	return server
}

// InitClient 初始化各客户端连接
func (ls *LogExternalServer) InitClient() {
	//MySQL
	db, err := infra.NewMysqlDB(ls.conf.Mysql)
	if err != nil {
		common.Logger.Fatal(err.Error())
	}
	common.Logger.Infof("mysql setting %+v", db.Statement)
	ls.conf.DB = db
}

func (ls *LogExternalServer) Start() {

}
