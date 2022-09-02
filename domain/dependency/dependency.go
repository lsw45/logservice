package dependency

import "log-ext/domain/entity"

//go:generate mockgen -source ../dependency/dependency.go -destination ../../mock/mock_dependency.go -package mock
type OpensearchRepo interface {
	Count()
	ListLog()
	Filter()
}

type MysqlRepo interface {
	GetUser(id int) (*entity.User, error)
	GetUserConfigName(ingestID, version string) (string, error)
}
