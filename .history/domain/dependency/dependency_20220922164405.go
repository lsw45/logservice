package dependency

import (
	"context"
	"log-ext/common/errorx"
	"log-ext/domain/entity"
)

//go:generate mockgen -source ../dependency/dependency.go -destination ../../mock/mock_dependency.go -package mock
type OpensearchRepo interface {
	SearchRequest(indexNames []string, content string) (map[string]interface{}, error)
	IndicesDeleteRequest(indexNames []string) ([]byte, error)
}

type MysqlRepo interface {
	// 用户模块
	GetUser(id int) (*entity.User, error)
	GetUserConfigName(ingestID, version string) (string, error)

	// 消息回调模块
	ExitsNotifyByUUId(uuid string) (bool, error)
	SaveNotifyMessage(msg *entity.NotifyDeployMessage) error
	SaveDeployeIngestTask(tasks []*entity.DeployIngestModel) (map[string]int, error)
	UpdateDeployeIngestTask(id []int, status int) error
}

type TunnelRepo interface {
	UploadFile(file_path, ip string) error
	ShellTask(env, project int, corporationId, server string, async bool) (bool, error)
}

type RedisRepo interface {
	Get(ctx context.Context, key string) (string, *errorx.CodeError)
}

type ElasticsearchDependency interface {
	SearchRequest(indexNames []string, query *entity.QueryDocs) (map[string]interface{}, error)
	IndicesDeleteRequest(indexNames []string) ([]byte, error)
}
