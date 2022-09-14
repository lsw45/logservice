package repository

import (
	"log-ext/domain/dependency"
	"log-ext/domain/entity"
	"log-ext/infra"
)

var _ dependency.MysqlRepo = (*MysqlRepo)(nil)

type MysqlRepo struct {
	infra.MysqlInfra
}

func NewMysqlRepo() *MysqlRepo {
	return &MysqlRepo{defaultInfra.Mysql}
}

func (m *MysqlRepo) GetUser(id int) (*entity.User, error) {
	return nil, nil
}

func (m *MysqlRepo) GetUserConfigName(ingestID, version string) (string, error) {
	return "nil", nil
}

func (m *MysqlRepo) ExitsNotifyByUUId(uuid string) (bool, error) {
	return false, nil
}

func (m *MysqlRepo) SaveNotifyMessage(*entity.NotifyDeployMessage) (id int, err error) {
	return 0, nil
}

func (m *MysqlRepo) UpdateNotifyDeployed(status int) error {
	return nil
}
