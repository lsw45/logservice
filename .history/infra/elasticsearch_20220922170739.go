package infra

import (
	"context"
	"errors"
	"log-ext/common"
	"log-ext/domain/entity"

	"github.com/olivere/elastic/v7"
)

type ElasticsearchInfra interface {
	SearchRequest(indexNames []string, quer *entity.QueryDocs) (*elastic.Response, error)
	IndicesDeleteRequest(indexNames []string) (*elastic.Response, error)
}

type elasticsearch struct {
	Client *elastic.Client
}

func NewElasticsearch(conf common.AppConfig) (*elasticsearch, error) {
	client, err := elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetURL(conf.Elasticsearch.Address...),
		elastic.SetBasicAuth(conf.Elasticsearch.Username, conf.Elasticsearch.Password),
	)

	if err != nil {
		return nil, err
	}
	return &elasticsearch{client}, nil
}

func (es *elasticsearch) SearchRequest(indexNames []string, quer *entity.QueryDocs) (*elastic.SearchResult, error) {
	query := elastic.NewMatchAllQuery()
	res, err := es.Client.Search().Index(indexNames...).From(quer.From).Size(quer.Size).SortBy(quer.Sort...).Query(query).Do(context.Background())

	if err != nil {
		return nil, err
	}
	if res == nil {
		err = errors.New("expected results != nil; got nil")
		return nil, err
	}
	if res.Hits == nil {
		err = errors.New("expected SearchResult.Hits != nil; got nil")
		return nil, err
	}
	if got, want := res.TotalHits(), int64(1); got < want {
		err = errors.New("got SearchResult.TotalHits() = 0")
		return nil, err
	}
	if got, want := len(res.Hits.Hits), 1; got < want {
		err = errors.New("got len(SearchResult.Hits.Hits) = 0")
		return nil, err
	}
	return res, nil
}
