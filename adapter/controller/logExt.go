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
	// logsrv := e.Group("/logservice2").Use(AuthCheck())
	logsrv := e.Group("/logservice2")
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
