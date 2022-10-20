package infra

import (
	"context"
	"fmt"
	"log"
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
	list, err := es.Search("server-y", "should", "(type:2)")

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, v := range list.Hits.Hits {
		fmt.Printf("%+v", v)
	}
}
