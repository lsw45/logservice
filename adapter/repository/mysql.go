package repository

import (
	"log-ext/common"
	"log-ext/domain/dependency"
	"log-ext/domain/entity"
	"log-ext/infra"

	"gorm.io/gorm"
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
	
	if err == gorm.ErrRecordNotFound {
		common.Logger.Infof("ExitsNotifyByUUId search error: %+v", err)
		return false, nil
	}

	if err != nil {
		common.Logger.Errorf("ExitsNotifyByUUId search error: %+v", err)
		return false, err
	}
	return exit, nil
}

func (m *MysqlRepo) SaveNotifyMessage(msg *entity.NotifyDeployMessage) error {
	return nil
}

func (m *MysqlRepo) UpdateDeployeIngestTask(id []int, status int) error {
	return nil
}

func (m *MysqlRepo) SaveDeployeIngestTask(tasks []*entity.DeployIngestModel) (map[string]int, error) {
	ids, err := m.MysqlInfra.SaveDeployeIngestTask(tasks)
	return ids, err
}
