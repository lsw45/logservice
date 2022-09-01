package domain

import (
	"log-ext/domain/dependency"
)

type SearchService interface {
	Search()
	List()
}

func NewSearchLogService(repo dependency.OpensearchRepo) SearchService {
	return &searchLogService{repo: repo}
}

type searchLogService struct {
	repo dependency.OpensearchRepo
}

func (srv *searchLogService) Search() {

}

func (srv *searchLogService) List() {
	//s.repo.List()
}
