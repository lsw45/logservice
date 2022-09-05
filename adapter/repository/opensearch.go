package repository

import (
	"log-ext/common"
	"log-ext/domain/dependency"
	"log-ext/infra"
)

var _ dependency.OpensearchRepo = (*OpensearchRepo)(nil)

type OpensearchRepo struct {
	infra.OpensearchInfra
}

func NewOpensearchRepo() *OpensearchRepo {
	return &OpensearchRepo{defaultInfra.Opensearch}
}

func (open *OpensearchRepo) SearchRequest(content string) ([]byte, error) {
	resp, err := open.OpensearchInfra.SearchRequest(content)
	if err != nil {
		common.Logger.Errorf("opensearch search error:%s", err.Error())
		return nil, err
	}
	buff := make([]byte, 1024)
	_, err = resp.Body.Read(buff)

	return buff, err
}

func (open *OpensearchRepo) IndicesDeleteRequest(indexNames []string) ([]byte, error) {
	resp, err := open.OpensearchInfra.IndicesDeleteRequest(indexNames)
	if err != nil {
		common.Logger.Errorf("opensearch delete index error:%s", err.Error())
		return nil, err
	}
	buff := make([]byte, 100)
	_, err = resp.Body.Read(buff)

	return buff, err
}
