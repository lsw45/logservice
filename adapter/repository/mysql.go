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
	exit, err := m.MysqlInfra.ExitsNotifyByUUId(uuid)
	return exit, err
}

func (m *MysqlRepo) SaveNotifyMessage(msg *entity.NotifyDeployMessage) error {
	return nil
}

func (m *MysqlRepo) UpdateDeployeIngestTask(id int, status int) error {
	return nil
}

func (m *MysqlRepo) SaveDeployeIngestTask(tasks []*entity.DeployIngestTable) error {
	err := m.MysqlInfra.SaveDeployeIngestTask(tasks)
	return err
}
