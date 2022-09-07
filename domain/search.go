package domain

import (
	"fmt"
	"log-ext/common"
	"log-ext/domain/dependency"
	"log-ext/domain/entity"
)

type SearchService interface {
	Histogram()
	SearchLogsByFilter(filter *entity.LogsFilter) ([]entity.LogsResult, error)
}

func NewSearchLogService(depOpen dependency.OpensearchRepo) SearchService {
	return &SearchLogService{depOpen: depOpen}
}

type SearchLogService struct {
	depOpen dependency.OpensearchRepo
}

func (srv *SearchLogService) Histogram() {

}

func (srv *SearchLogService) SearchLogsByFilter(filter *entity.LogsFilter) ([]entity.LogsResult, error) {
	// struct to string

	content := `{
		"query": { "match_all":{} }
    }`

	data, err := srv.depOpen.SearchRequest(filter.Indexs, content)
	if err != nil {
		common.Logger.Errorf("domain log search error: %+v", err)
		return nil, err
	}

	hits := data["body"].(map[string]interface{})["hits"].(map[string]interface{})
	fmt.Printf("%+v\n", hits)

	// for _, v := range data["body"].(map[string]interface{})["hits"].(map[string]interface{})["hits"].([]map[string]interface{}) {
	// 	for k1, v1 := range v {
	// 		fmt.Printf("%+v\n", k1)
	// 		fmt.Printf("%+v\n", v1)
	// 	}
	// }
	// var result []entity.LogsResult
	// err = json.Unmarshal(data, &result)
	// if err != nil {
	// 	common.Logger.Errorf("domain error: %+v", err)
	// 	return nil, err
	// }

	return nil, nil
}
