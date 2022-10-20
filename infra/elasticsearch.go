package infra

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log-ext/common"
	"log-ext/domain/entity"
	"time"

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
	Aggregation(req entity.AggregationReq) (*elastic.SearchResult, error)
	SearchRequest(indexNames []string, quer *entity.QueryDocs) (*elastic.SearchResult, error)
	Histogram(search *entity.DateHistogramReq) ([]entity.Buckets, int64, error)
	IndicesDeleteRequest(indexNames []string) (*elastic.Response, error)
	NearbyDoc(indexName string, times int64, num int) ([]*elastic.SearchHit, error)
}

var _ ElasticsearchInfra = new(elasticsearch)

type elasticsearch struct {
	Client *elastic.Client
}

var eslog bytes.Buffer

func NewElasticsearch(conf common.Elasticsearch) (*elasticsearch, error) {
	tout := log.New(&eslog, "TRACER ", log.LstdFlags)
	client, err := elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetURL(conf.Address...),
		elastic.SetBasicAuth(conf.Username, conf.Password),
		elastic.SetTraceLog(tout),
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
	// timeRange := elastic.NewRangeQuery("time").Gte(search.StartTime).Lte(search.EndTime)
	if len(search.Query) == 0 {
		search.Query = `{"query":{"match_all":{}},"track_total_hits":true}`
	}

	timeRange := elastic.NewRangeQuery("time").Gt(search.StartTime).Lt(search.EndTime)

	source := elastic.NewSearchSource().Query(timeRange)

	res, err := es.Client.Search().Index(indexNames...).SearchSource(source).Source(search.Query).
		From(search.From).Size(search.Size).SortBy(search.Sort...).
		Do(context.Background())

	common.Logger.Error(eslog.String())
	if err != nil {
		return nil, err
	}
	if res == nil {
		err = errors.New("got results is nil")
		return nil, err
	}
	if res.Hits == nil {
		common.Logger.Info(eslog.String())
		err = errors.New("got SearchResult.Hits is nil")
		return nil, err
	}

	return res, nil
}

func (es *elasticsearch) IndicesDeleteRequest(indexNames []string) (*elastic.Response, error) {
	return nil, nil
}

func (es *elasticsearch) Histogram(search *entity.DateHistogramReq) ([]entity.Buckets, int64, error) {
	search.GroupName = "dateGroup"

	// 如果interval大于所有日志的时间，则查询到
	// "buckets" : [
	// {
	//   "key" : 0,
	//   "doc_count" : 10832
	// }
	// 第一个doc作为起始时间
	sort := []elastic.Sorter{elastic.NewFieldSort("time").Asc()}

	h := elastic.NewHistogramAggregation().Field("time").Interval(float64(search.Interval))

	timeRange := elastic.NewRangeQuery("time").Gte(search.StartTime).Lte(search.EndTime)
	builder := es.Client.Search().Index(search.Indexs...).Query(timeRange).TrackTotalHits(true).
		Size(1).SortBy(sort...). // 拿到第一个doc
		Pretty(true)

	res, err := builder.Aggregation(search.GroupName, h).Do(context.TODO())

	common.Logger.Infof(eslog.String())
	if err != nil {
		return nil, 0, err
	}

	aggs := res.Aggregations
	if aggs == nil {
		err = errors.New("got Aggregations is nil")
		return nil, 0, err
	}
	if len(aggs[search.GroupName]) == 0 {
		return nil, 0, nil
	}

	histogra := &entity.DateHistAggre{}
	err = json.Unmarshal(aggs[search.GroupName], histogra)
	if err != nil {
		common.Logger.Error(err.Error())
		return nil, 0, err
	}

	if len(histogra.Buckets) == 1 {
		b := histogra.Buckets[0]

		if (b.Key == 0 || len(b.KeyAsString) == 0) && b.DocCount > 0 {
			hit := make(map[string]interface{})

			if res.Hits.TotalHits.Value > 0 {
				xx, _ := res.Hits.Hits[0].Source.MarshalJSON()
				_ = json.Unmarshal(xx, &hit)
			}

			b.Key = int64(hit["time"].(float64))

			tm := time.Unix(b.Key.(int64), 0)
			b.KeyAsString = tm.Format("2006-01-02 15:04:05")
		}
		histogra.Buckets[0] = b
	}

	return histogra.Buckets, res.Hits.TotalHits.Value, nil
}

func (es *elasticsearch) NearbyDoc(indexName string, times int64, num int) ([]*elastic.SearchHit, error) {
	sortField := "time"

	timeRange := elastic.NewRangeQuery(sortField).Gte(times)
	afterRes, err := es.Client.Search().Index(indexName).Query(timeRange).Sort(sortField, true).Size(num).Do(context.Background())
	if err != nil {
		return nil, err
	}
	if afterRes == nil {
		err = errors.New("search results is nil")
		return nil, err
	}
	if afterRes.Hits == nil {
		err = errors.New("got aggregation.Hits is nil")
		return nil, err
	}

	timeRange = elastic.NewRangeQuery(sortField).Lt(times)
	preRes, err := es.Client.Search().Index(indexName).Query(timeRange).Sort(sortField, true).Size(num).Do(context.Background())
	if err != nil {
		return nil, err
	}

	if afterRes == nil && preRes == nil {
		err = errors.New("aggregation results is nil")
		return nil, err
	}

	if afterRes.Hits == nil && preRes.Hits == nil {
		err = errors.New("got aggregation.Hits is nil")
		return nil, err
	}

	searchHit := make([]*elastic.SearchHit, 0)
	searchHit = append(searchHit, preRes.Hits.Hits...)
	searchHit = append(searchHit, afterRes.Hits.Hits...)

	return searchHit, nil
}

func (es *elasticsearch) Aggregation(req entity.AggregationReq) (*elastic.SearchResult, error) {
	res, err := es.Client.Search().Index(req.Indexs...).Source(req.Aggs).Size(0).Do(context.Background())

	if err != nil {
		return nil, err
	}
	if res == nil {
		err = errors.New("aggregation results is nil")
		return nil, err
	}
	if res.Hits == nil {
		err = errors.New("got aggregation.Hits is nil")
		return nil, err
	}

	return res, nil
}

func (es *elasticsearch) Search(index string, must string, body string) (*elastic.SearchResult, error) {
	query := `{
"query": {
    "bool": {
		"%s": [
			{
				"match_phrase": {
					"body": "%s"
				}
			}
		]
	}
  }
}`
	query = fmt.Sprintf(query, must, body)
	res, err := es.Client.Search().Index(index).Source(query).Size(0).Do(context.Background())

	if err != nil {
		return nil, err
	}
	if res == nil {
		err = errors.New("aggregation results is nil")
		return nil, err
	}
	if res.Hits == nil {
		err = errors.New("got aggregation.Hits is nil")
		return nil, err
	}

	return res, nil
}
