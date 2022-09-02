package main

import (
	"context"
	"log-ext/adapter/controller"
	"log-ext/common"
	"log-ext/infra"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	common.Logger.Info("Service Starting")
	conf := common.NewAppConfig()

	common.Logger.Info("Start Http")
	server := NewServer(conf)
	server.Start()

	common.Logger.Info("wait signal")
	server.AwaitSignal()
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

func (ls *LogExternalServer) AwaitSignal() {
	c := make(chan os.Signal, 1)
	signal.Reset(syscall.SIGTERM, syscall.SIGINT)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

	<-c
	ls.GracefulShutDown()
	close(c)
}

func (ls *LogExternalServer) GracefulShutDown() {
	// 通知其他协程退出
	controller.ExitChan <- nil
	defer close(controller.ExitChan)
	// 关闭服务
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if err := ls.server.Shutdown(ctx); err != nil {
		common.Logger.Fatalf("Server Shutdown:%v", err)
	}

	// DB 关闭
	db, _ := ls.conf.DB.DB()
	defer db.Close()

	common.Logger.Info("Server Exiting")
}
