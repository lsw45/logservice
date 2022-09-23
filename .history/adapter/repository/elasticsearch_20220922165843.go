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
if res == nil {
		t.Fatal("expected results != nil; got nil")
	}
	if res.Hits == nil {
		t.Fatal("expected SearchResult.Hits != nil; got nil")
	}
	if got, want := res.TotalHits(), int64(1); got != want {
		t.Fatalf("expected SearchResult.TotalHits() = %d; got %d", want, got)
	}
	if got, want := len(res.Hits.Hits), 1; got != want {
		t.Fatalf("expected len(SearchResult.Hits.Hits) = %d; got %d", want, got)
	}
	hit := res.Hits.Hits[0]
	if hit.Index != testQueryIndex {
		t.Fatalf("expected SearchResult.Hits.Hit.Index = %q; got %q", testQueryIndex, hit.Index)
	}
	got := string(hit.Source)
}
func (this *Elasticsearch) IndicesDeleteRequest(indexNames []string) ([]byte, error) {

}
