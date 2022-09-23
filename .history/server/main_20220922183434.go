package main

import (
	"context"
	"log-ext/adapter/controller"
	"log-ext/adapter/repository"
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

	// 初始化基础设施
	initClient(conf)

	common.Logger.Info("Start Http")
	logServer := NewLogServer(conf)

	common.Logger.Info("wait signal")
	logServer.AwaitSignal()
}

type LogServer struct {
	conf *common.AppConfig
	web  *http.Server
}

func NewLogServer(conf *common.AppConfig) *LogServer {
	server := &LogServer{conf: conf}

	engine := gin.New()
	engine.Use(gin.Recovery())
	gin.SetMode(conf.Server.RunMode)

	// 创建控制器
	logExt := controller.NewLogExtServer(conf)
	logExt.RegisterRouter(engine)

	server.web = &http.Server{
		Addr:           ":" + strconv.Itoa(conf.Server.HTTPPort),
		Handler:        engine,
		TLSConfig:      nil,
		ReadTimeout:    conf.Server.ReadTimeOut * time.Second,
		WriteTimeout:   conf.Server.WriteTimeOut * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	err := server.web.ListenAndServe()
	if err != nil {
		common.Logger.Fatal(err.Error())
	}
	return server
}

// InitClient 初始化各客户端连接
func initClient(conf *common.AppConfig) {
	elastic, err := infra.NewElasticsearch(conf.Elasticsearch)
		if err != nil {
			common.Logger.Fatalf("new elasticsearch infra: %v", err)
		}

		// mysql
		mysql, err := infra.NewMysql(conf.Mysql)
		if err != nil {
			common.Logger.Fatal(err.Error())
		}
		common.Logger.Infof("mysql setting: %+v", mysql.DB.Statement)

		// opensearch
		openDB, err := infra.NewOpensearch(conf.Opensearch)
		if err != nil {
			common.Logger.Fatal(err.Error())
		}
		common.Logger.Infof("opensearch url: %+v", openDB.Client.Transport.(*opensearchtransport.Client).URLs())

		// redis
		redis, err := infra.NewRedis(conf.Redis)
		if err != nil {
			common.Logger.Fatal(err.Error())
		}
		common.Logger.Infof("redis client: %+v", redis)

		// tunnel
		tunnel := infra.NewTunnelClient(conf.Tunnel)
		
}

func (ls *LogServer) Start() {

}

func (ls *LogServer) AwaitSignal() {
	c := make(chan os.Signal, 1)
	signal.Reset(syscall.SIGTERM, syscall.SIGINT)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

	<-c
	ls.GracefulShutDown()
	close(c)
}

func (ls *LogServer) GracefulShutDown() {
	// 通知其他协程退出
	controller.ExitChan <- nil
	defer close(controller.ExitChan)
	// 关闭服务
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if err := ls.web.Shutdown(ctx); err != nil {
		common.Logger.Fatalf("Server Shutdown:%v", err)
	}

	// DB 关闭
	repository
	db, _ := ls.conf.DB.DB()
	defer db.Close()

	common.Logger.Info("Server Exiting")
}
