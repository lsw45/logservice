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
	return &OpensearchRepo{defaultRepo.Opensearch}
}

func (open *OpensearchRepo) SearchRequest(indexNames []string, content string) (map[string]interface{}, error) {
	resp, err := open.OpensearchInfra.SearchRequest(indexNames, content)
	defer resp.Body.Close()
	if err != nil {
		common.Logger.Errorf("repository opensearch search error:%s", err.Error())
		return nil, err
	}

	// 将查询结果转化为map
	objStr, err := common.ReadAsString(resp.Body)
	if err != nil {
		common.Logger.Errorf("infra read byte error: %+v", err)
		return nil, err
	}

	obj := common.ParseJSON(objStr)
	resB := common.AssertAsMap(obj)

	result := make(map[string]interface{})

	err = common.Convert(map[string]interface{}{
		"body":   resB,
		"header": resp.Header,
	}, &result)

	if err != nil {
		common.Logger.Errorf("infra convert map error: %+v", err)
		return nil, err
	}

	return result, nil
}

func (open *OpensearchRepo) IndicesDeleteRequest(indexNames []string) ([]byte, error) {
	resp, err := open.OpensearchInfra.IndicesDeleteRequest(indexNames)
	if err != nil {
		common.Logger.Errorf("repository opensearch delete index error:%s", err.Error())
		return nil, err
	}
	buff := make([]byte, 100)
	_, err = resp.Body.Read(buff)

	return buff, err
}
