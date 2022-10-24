package infra

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"log-ext/common"
	"log-ext/domain/entity"
	"strings"
	"time"

	"github.com/olivere/elastic/v7"
)

type ElasticsearchInfra interface {
	Aggregation(req entity.AggregationReq) (*elastic.SearchResult, error)
	SearchRequest(indexNames []string, quer *entity.QueryDocs) (*elastic.SearchResult, error)
	Histogram(search *entity.DateHistogramReq) ([]entity.Buckets, int64, error)
	IndicesDeleteRequest(indexNames []string) (*elastic.Response, error)
	NearbyDoc(indexName string, times int64, num int) ([]*elastic.SearchHit, error)
	IndexExists(index string) (bool, error)
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

func (es *elasticsearch) IndexExists(index string) (bool, error) {
	exist, err := es.Client.IndexExists(index).Do(context.TODO())
	if err != nil {
		common.Logger.Errorf("index exists: %v", err)
	}
	return exist, err
}

func (es *elasticsearch) IndicesDeleteRequest(indexNames []string) (*elastic.Response, error) {
	return nil, nil
}

func (es *elasticsearch) Histogram(search *entity.DateHistogramReq) ([]entity.Buckets, int64, error) {
	group_name := "dateGroup"

	// 如果interval大于所有日志的时间，则查询到
	// "buckets" : [
	// {
	//   "key" : 0,
	//   "doc_count" : 10832
	// }
	// 第一个doc作为起始时间
	sort := []elastic.Sorter{elastic.NewFieldSort(entity.LogSortField).Asc()}

	h := elastic.NewHistogramAggregation().Field(entity.LogSortField).Interval(float64(search.Interval))
	qb := newBoolQuery(search.Query, search.StartTime, search.EndTime)

	builder := es.Client.Search().Index(search.Indexs...).Query(qb).TrackTotalHits(true).
		Size(1).SortBy(sort...). // 拿到第一个doc
		Pretty(true)

	res, err := builder.Aggregation(group_name, h).Do(context.TODO())

	common.Logger.Infof(eslog.String())
	if err != nil {
		return nil, 0, err
	}

	aggs := res.Aggregations
	if aggs == nil {
		err = errors.New("got Aggregations is nil")
		return nil, 0, err
	}
	if len(aggs[group_name]) == 0 {
		return nil, 0, nil
	}

	histogra := &entity.DateHistAggre{}
	err = json.Unmarshal(aggs[group_name], histogra)
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

			b.Key = int64(hit[entity.LogSortField].(float64))

			tm := time.Unix(b.Key.(int64), 0)
			b.KeyAsString = tm.Format("2006-01-02 15:04:05")
		}
		histogra.Buckets[0] = b
	}

	return histogra.Buckets, res.Hits.TotalHits.Value, nil
}

func (es *elasticsearch) NearbyDoc(indexName string, times int64, num int) ([]*elastic.SearchHit, error) {

	timeRange := elastic.NewRangeQuery(entity.LogSortField).Gte(times)
	afterRes, err := es.Client.Search().Index(indexName).Query(timeRange).Sort(entity.LogSortField, true).Size(num).Do(context.Background())
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

	timeRange = elastic.NewRangeQuery(entity.LogSortField).Lt(times)
	preRes, err := es.Client.Search().Index(indexName).Query(timeRange).Sort(entity.LogSortField, true).Size(num).Do(context.Background())
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

func (es *elasticsearch) SearchRequest(index []string, search *entity.QueryDocs) (*elastic.SearchResult, error) {
	qb := newBoolQuery(search.Query, search.StartTime, search.EndTime)
	res, err := es.Client.Search().Index(index...).Query(qb).From(search.From).Size(search.Size).SortBy(search.Sort...).TrackTotalHits(true).Do(context.Background())

	if err != nil {
		// 不存在的索引，返回空
		if strings.Index(err.Error(), "index_not_found_exception") > 0 {
			return nil, nil
		}
		common.Logger.Infof(eslog.String())
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

func newBoolQuery(queryString string, start, end int64) *elastic.BoolQuery {
	qb := elastic.NewBoolQuery()
	if len(queryString) != 0 {
		qb = qb.Must(elastic.NewQueryStringQuery(queryString).TimeZone("Asia/Shanghai").AnalyzeWildcard(true))
	}
	qb = qb.Filter(elastic.NewRangeQuery(entity.LogSortField).Gte(start).Lte(end))
	return qb
}
