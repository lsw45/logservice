package dependency

import (
	"context"
	"log-ext/common/errorx"
	"log-ext/domain/entity"
)

//go:generate mockgen -source ../dependency/dependency.go -destination ../../mock/mock_dependency.go -package mock
type OpensearchRepo interface {
	SearchRequest(content string) ([]byte, error)
	IndicesDeleteRequest(indexNames []string) ([]byte, error)
}

type MysqlRepo interface {
	GetUser(id int) (*entity.User, error)
	GetUserConfigName(ingestID, version string) (string, error)
}

type RedisRepo interface {
	Get(ctx context.Context, key string) (string, *errorx.CodeError)
}
