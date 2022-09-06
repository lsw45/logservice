package domain

import (
	"encoding/json"
	"log-ext/domain/dependency"
	"log-ext/domain/entity"
)

type SearchService interface {
	Histogram()
	SearchLogsByFilter(filter *entity.LogsFilter) ([]entity.LogsResult, error)
}

func NewSearchLogService(repo dependency.OpensearchRepo) SearchService {
	return &SearchLogService{repo: repo}
}

type SearchLogService struct {
	repo dependency.OpensearchRepo
}

func (srv *SearchLogService) Histogram() {

}

func (srv *SearchLogService) SearchLogsByFilter(filter *entity.LogsFilter) ([]entity.LogsResult, error) {
	// struct to string
	var content string = ""

	data, err := srv.repo.SearchRequest(content)
	if err != nil {
		return nil, err
	}

	var result []entity.LogsResult
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
