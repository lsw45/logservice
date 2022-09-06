package controller

import (
	"encoding/json"
	"log-ext/adapter/repository"
	"log-ext/common"
	"log-ext/common/errorx"
	"log-ext/domain/entity"
	"log-ext/infra"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

var ExitChan = make(chan interface{}, 1)
var authRedis *repository.Redis
var authKey = "auth"

type AppServer interface {
	RegisterRouter(e *gin.Engine)
}

func NewLogExtServer(conf *common.AppConfig) AppServer {
	repository.SetRepoInfra(conf)
	authRedis = &repository.Redis{RedisInfra: &infra.Redis{Client: conf.RedisCli}}

	return &logExtServer{
		searchCtl: *NewSearchController(repository.NewOpensearchRepo()),
	}
}

type logExtServer struct {
	searchCtl SearchController
}

func (ctl *logExtServer) RegisterRouter(e *gin.Engine) {
	if authRedis == nil {
		common.Logger.Fatal("auth server is nil")
	}
	e.Use(AuthCheck())
	logsrv := e.Group("logservice2")
	logsrv.GET("/logs", ctl.searchCtl.SearchLogsByFilter)
	logsrv.GET("/histogram", ctl.searchCtl.Histogram)
}

func AuthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		authInfo, ok := common.ParseBasicAuth(c.GetHeader("Authorization"))
		if !ok {
			c.Abort()
			c.JSON(http.StatusUnauthorized, common.ThrowErr(errorx.NewErrCode(errorx.AUTH_ERROR)))
		}

		data, errs := authRedis.Get(c.Request.Context(), authInfo)
		if errs != nil {
			common.Logger.Errorf("【API-SRV-ERR】 %+v", errors.Wrap(errs, "获取鉴权信息失败"))
			c.JSON(http.StatusInternalServerError, common.ThrowErr(errs))
		}

		var userInfo entity.UserInfo
		err := json.Unmarshal([]byte(data), &userInfo)
		if err != nil {
			common.Logger.Errorf("【API-SRV-ERR】 %+v", errors.Wrap(errs, "解析用户鉴权信息失败"))
			c.JSON(http.StatusInternalServerError, common.ThrowErr(errs))
		}

		company, err := strconv.Atoi(userInfo.CorporationId)
		if err != nil {
			err = errors.Wrap(err, "解析用户鉴权信息失败")
			common.Logger.Errorf("【API-SRV-ERR】 %+v", err)
			return
		}

		userInfo.Company = int64(company)
		c.Set(authKey, userInfo)
		c.Next()
	}
}
