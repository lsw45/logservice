package infra

import (
	"context"
	"errors"
	"log-ext/common"
	"log-ext/domain/entity"

	"github.com/olivere/elastic/v7"
)

/*
{
    "_index" : "game_static-2022.09.29",
    "_id" : "-f5FiIMBP53dBzVBjHbr",
    "_score" : null,
    "_ignored" : [
      "event.original.keyword"
    ],
	"_source" : {
      "@version" : "1",
      "uuid" : "YzE+3wPadjA7QohW",
      "role_id" : 12334,
      "IMEI" : "xxxx",
      "role_name" : "cocos",
      "time" : 1664269483,
      "mac_address" : "00-15-5D-0C-55-55",
      "@timestamp" : "2022-09-29T08:03:22.822548117Z",
      "os_name" : "ios 16",
      "ip" : "127.0.0.1",
      "game_id" : "cc88",
      "os_ver" : "1.1.1",
      "type" : "game_static",
      "operation" : "LogoutRole",
      "server_id" : "10001",
      "app_channel" : "AppStore",
      "logout_time" : 1664269483,
      "account_id" : "cocos",
      "network" : "WIFI",
      "service" : "/usr/local/go",
      "index" : "55-1-3",
      "country_code" : "",
      "state" : {
        "filename" : "/opt/carey/modify_es_index/source_loggie/test.log",
        "hostname" : "paas-dev",
        "bytes" : 431,
        "source" : "kdump",
        "timestamp" : "2022-09-29T16:03:19.790Z",
        "offset" : 0,
        "pipeline" : "demo"
      }
    },
    "sort" : [
      1664438602822
    ]
}
*/
type ElasticsearchInfra interface {
	SearchRequest(indexNames []string, quer *entity.QueryDocs) (*elastic.SearchResult, error)
	IndicesDeleteRequest(indexNames []string) (*elastic.Response, error)
}

var _ ElasticsearchInfra = new(elasticsearch)

type elasticsearch struct {
	Client *elastic.Client
}

func NewElasticsearch(conf common.Elasticsearch) (*elasticsearch, error) {
	client, err := elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetURL(conf.Address...),
		elastic.SetBasicAuth(conf.Username, conf.Password),
	)

	if err != nil {
		return nil, err
	}

	_, _, err = client.Ping(conf.Address[0]).Do(context.Background())

	if err != nil {
		return nil, err
	}
	return &elasticsearch{client}, nil
}

func (es *elasticsearch) SearchRequest(indexNames []string, search *entity.QueryDocs) (*elastic.SearchResult, error) {
	query := elastic.RawStringQuery(`{"match_all":{}}`)

	res, err := es.Client.Search().Index(indexNames...).Query(query).From(search.From).Size(search.Size).SortBy(search.Sort...).Do(context.Background())

	if err != nil {
		return nil, err
	}
	if res == nil {
		err = errors.New("got results is nil")
		return nil, err
	}
	if res.Hits == nil {
		err = errors.New("got SearchResult.Hits is nil")
		return nil, err
	}

	return res, nil
}

func (es *elasticsearch) IndicesDeleteRequest(indexNames []string) (*elastic.Response, error) {
	return nil, nil
}
