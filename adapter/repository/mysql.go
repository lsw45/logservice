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
	return &MysqlRepo{defaultRepo.Mysql}
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

func (m *MysqlRepo) SaveNotifyMessage(msg *entity.NotifyMsgModel) error {
	return m.MysqlInfra.SaveNotifyMessage(msg)
}

func (m *MysqlRepo) UpdateDeployeIngestTask(id []int, status int) error {
	return m.MysqlInfra.UpdateDeployeIngestTask(id, status)
}

func (m *MysqlRepo) SaveDeployeIngestTask(tasks []*entity.DeployIngestModel) (map[string]int, error) {
	return m.MysqlInfra.SaveDeployeIngestTask(tasks)
}

func (m *MysqlRepo) ReleaseRegion(regionId int) error {
	return m.MysqlInfra.ReleaseRegion(regionId)
}
