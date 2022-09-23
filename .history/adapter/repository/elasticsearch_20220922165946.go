package repository

import (
	"log-ext/domain/entity"
	"log-ext/infra"
)

type Elasticsearch struct {
	infra infra.ElasticsearchInfra
}

func (this *Elasticsearch) SearchRequest(indexNames []string, query *entity.QueryDocs) (map[string]interface{}, error) {
this.infra.SearchRequest()

	hit := res.Hits.Hits[0]
	if hit.Index != testQueryIndex {
		t.Fatalf("expected SearchResult.Hits.Hit.Index = %q; got %q", testQueryIndex, hit.Index)
	}
	got := string(hit.Source)
}
func (this *Elasticsearch) IndicesDeleteRequest(indexNames []string) ([]byte, error) {

}
