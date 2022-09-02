package repository

import (
	"log-ext/domain/dependency"
	"log-ext/domain/entity"
)

var _ dependency.MysqlRepo = (*MysqlRepo)(nil)

type MysqlRepo struct {
	//infra.MysqlInfra
}

func (m *MysqlRepo) GetUser(id int) (*entity.User, error) {
	return nil, nil
}

func (m *MysqlRepo) GetUserConfigName(ingestID, version string) (string, error){
	return "nil", nil
}



