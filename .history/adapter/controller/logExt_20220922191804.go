package controller

import (
	"context"
	"encoding/json"
	"log-ext/adapter/repository"
	"log-ext/common"
	"log-ext/common/errorx"
	"log-ext/domain/entity"
	"log-ext/infra"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

var authRedis *repository.Redis
var authKey = "auth"

type AppServer interface {
	RegisterRouter(e *gin.Engine)
}

func NewLogExtServer(conf *common.AppConfig) AppServer {
	// authRedis = &repository.Redis{RedisInfra: &infra.Redis{Client: conf.RedisCli}}
	authRedis = repository.NewRedis()
	return &logExtServer{
		searchCtl: *NewSearchController(repository.NewElasticsearchRepo()),
		deployCtl: *NewNotifyController(repository.NewMysqlRepo(), repository.NewTunnelRepo()),
	}
}

type logExtServer struct {
	searchCtl SearchController
	deployCtl DeployController
}

func (ctl *logExtServer) RegisterRouter(e *gin.Engine) {
	// logsrv := e.Group("/logservice2").Use(AuthCheck())
	logsrv := e.Group("/logservice2")
	logsrv.GET("/logs", ctl.searchCtl.SearchLogsByFilter)
	logsrv.GET("/histogram", ctl.searchCtl.Histogram)

	notify := logsrv.Use(timeoutMiddleware(2 * time.Second))
	notify.POST("/notify", ctl.deployCtl.Notify)
}

func AuthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		authInfo, ok := common.ParseBasicAuth(c.GetHeader("Authorization"))
		if !ok {
			c.Abort()
			c.JSON(http.StatusUnauthorized, common.ThrowErr(errorx.NewErrCode(errorx.AUTH_ERROR)))
		}

		if len(authInfo) == 0 {
			c.Abort()
			c.JSON(http.StatusUnauthorized, common.ThrowErr(errorx.NewErrMsg("Authorization is null")))
		}

		data, errs := authRedis.Get(c.Request.Context(), authInfo)
		if errs != nil {
			c.Abort()
			common.Logger.Errorf("【API-SRV-ERR】 %+v", errors.Wrap(errs, "获取鉴权信息失败"))
			c.JSON(http.StatusInternalServerError, common.ThrowErr(errs))
		}

		var userInfo entity.UserInfo
		err := json.Unmarshal([]byte(data), &userInfo)
		if err != nil {
			c.Abort()
			common.Logger.Errorf("【API-SRV-ERR】 %+v", errors.Wrap(err, "解析用户鉴权信息失败"))
			c.JSON(http.StatusInternalServerError, common.ThrowErr(errorx.NewErrMsg(err.Error())))
		}

		company, err := strconv.Atoi(userInfo.CorporationId)
		if err != nil {
			c.Abort()
			common.Logger.Errorf("【API-SRV-ERR】 %+v", err)
			c.JSON(http.StatusInternalServerError, common.ThrowErr(errorx.NewErrMsg(err.Error())))
		}

		userInfo.Company = int64(company)
		c.Set(authKey, userInfo)
		c.Next()
	}
}

func timeoutMiddleware(timeout time.Duration) func(c *gin.Context) {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)

		defer func() {
			if ctx.Err() == context.DeadlineExceeded {
				c.Writer.WriteHeader(http.StatusGatewayTimeout)
				c.Abort()
			}

			cancel() // clean the resource
		}()

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
