package infra

import (
	"context"
	"fmt"
	"log"
	"log-ext/domain/entity"
	"testing"

	"github.com/olivere/elastic/v7"
)

func TestNearby(t *testing.T) {
	tout := log.New(&eslog, "TRACER ", log.LstdFlags)
	client, err := elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetURL("http://10.0.0.73:9200"),
		elastic.SetTraceLog(tout),
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	timeRange := elastic.NewRangeQuery("time").Gt(1665244793)
	res, err := client.Search().Index("operator-55-1-3-2022.10.08").Query(timeRange).Sort("time", true).Size(3).Do(context.Background())

	timeRange = elastic.NewRangeQuery("time").Lt(1665244793)
	res, err = client.Search().Index("operator-55-1-3-2022.10.08").Query(timeRange).Sort("time", true).Size(3).Do(context.Background())

	fmt.Println(eslog.String())

	fmt.Println(err)
	fmt.Println(res)
}

func TestLuceneQuery(t *testing.T) {
	tout := log.New(&eslog, "TRACER ", log.LstdFlags)
	client, err := elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetURL("http://121.37.173.234:9200"),
		elastic.SetTraceLog(tout),
	)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	es := elasticsearch{Client: client}

	e := &entity.QueryDocs{
		Query:     "time: [1666260000 TO 1666261000]",
		Size:      1,
		From:      1,
		EndTime:   1666274687,
		StartTime: 1666260000,
		Sort:      []elastic.Sorter{elastic.NewFieldSort("time").Desc()},
	}

	list, err := es.SearchRequest([]string{"operator-55-1-5-2022.10"}, e)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if list != nil {
		for _, v := range list.Hits.Hits {
			fmt.Printf("%+v", v)
		}
	}
	fmt.Println(eslog.String())

	h := &entity.DateHistogramReq{
		Interval:  60,
		Query:     e.Query,
		EndTime:   e.EndTime,
		StartTime: e.StartTime,
		Indexs:    []string{"operator-55-1-4-2022.10"},
	}
	_, _, err = es.Histogram(h)

	fmt.Println(eslog.String())

}
