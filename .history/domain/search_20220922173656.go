package domain

import (
	"encoding/json"
	"log-ext/common"
	"log-ext/domain/dependency"
	"log-ext/domain/entity"
	"strconv"
)

type SearchService interface {
	Histogram()
	SearchLogsByFilter(filter *entity.LogsFilter) ([]byte, int, error)
}

func NewSearchLogService(depOpen dependency.OpensearchRepo) SearchService {
	return &SearchLogService{depOpen: depOpen}
}

type SearchLogService struct {
	depOpen dependency.OpensearchRepo
}

func (srv *SearchLogService) Histogram() {

}

func (srv *SearchLogService) SearchLogsByFilter(filter *entity.LogsFilter) ([]byte, int, error) {
	content := `{
		"query": { "match_all":{} }
    }`

	data, err := srv.depOpen.SearchRequest(filter.Indexs, content)
	if err != nil {
		common.Logger.Errorf("domain log search error: %+v", err)
		return nil, 0, err
	}

	var body map[string]interface{}
	var hits map[string]interface{}
	var result []interface{}
	var total int
	var ok bool

	if _, ok = data["body"]; !ok {
		common.Logger.Errorf("domain error: no body:%+v", err)
		return nil, 0, err
	}

	if body, ok = data["body"].(map[string]interface{}); !ok {
		common.Logger.Errorf("domain error: body transform map:%+v", err)
		return nil, 0, err
	}

	if _, ok = body["hits"]; !ok {
		common.Logger.Errorf("domain error: no hits:%+v", err)
		return nil, 0, err
	}

	if hits, ok = body["hits"].(map[string]interface{}); !ok {
		common.Logger.Errorf("domain error: hits transform map:%+v", err)
		return nil, 0, err
	}

	if _, ok = hits["hits"]; !ok {
		common.Logger.Errorf("domain error: no hits:%+v", err)
		return nil, 0, err
	}

	if result, ok = hits["hits"].([]interface{}); !ok {
		common.Logger.Errorf("domain error: hits transform slice:%+v", err)
		return nil, 0, err
	}

	if _, ok = hits["total"]; !ok {
		common.Logger.Errorf("domain error: no hits:%+v", err)
		return nil, 0, err
	}

	if num, ok := hits["total"].(map[string]interface{}); !ok {
		common.Logger.Errorf("domain error: hits transform slice:%+v", err)
		return nil, 0, err
	} else if _, ok := num["value"]; ok {
		var nums json.Number
		if nums, ok = num["value"].(json.Number); !ok {
			common.Logger.Errorf("domain error: hits transform slice:%+v", err)
			return nil, 0, err
		}
		total, err = strconv.Atoi(string(nums))
		if err != nil {
			common.Logger.Errorf("domain error: %+v", err)
			return nil, 0, err
		}
	}

	re,_ := json.Marshal(result)
	return result, total, nil
}
