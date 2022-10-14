package dependency

import (
	"context"
	"log-ext/common/errorx"
	"log-ext/domain/entity"

	"github.com/olivere/elastic/v7"
)

//go:generate mockgen -source ../dependency/dependency.go -destination ../../mock/mock_dependency.go -package mock
type MysqlRepo interface {
	// 用户模块
	GetUser(id int) (*entity.User, error)
	GetUserConfigName(ingestID, version string) (string, error)

	// 消息回调模块
	ExitsNotifyByUUId(uuid string) (string, error)
	SaveNotifyMessage(msg *entity.NotifyMsgModel) error
	SaveDeployeIngestTask(tasks []*entity.DeployIngestModel) (map[string]int, error)
	UpdateDeployeIngestTask(id []int, status int) error
	ReleaseRegion(regionId int) error
}

type TunnelRepo interface {
	UploadFile(file_path, ip, env string) error
	ShellTask(env, project int, corporationId, server string, async bool) (bool, error)
}

type RedisRepo interface {
	Get(ctx context.Context, key string) (string, *errorx.CodeError)
}

type ElasticsearchDependency interface {
	Aggregation(req entity.AggregationReq) (*elastic.SearchResult, error)
	SearchRequest(indexNames []string, query *entity.QueryDocs) (*elastic.SearchHits, error)
	IndicesDeleteRequest(indexNames []string) ([]byte, error)
	Histogram(query *entity.DateHistogramReq) ([]entity.Buckets, int64, error)
	NearbyDoc(indexName string, times int64, num int) ([]*elastic.SearchHit, error)
}
