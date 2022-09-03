package infra

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
	"log-ext/common"
	"net/http"
	"strings"
)

type OpensearchInfra interface {
	SearchRequest(content string) (*opensearchapi.Response, error)
	IndicesDeleteRequest(indexNames []string) (*opensearchapi.Response, error)
}

type Opensearch struct {
	Client *opensearch.Client
}

func NewOpensearch(conf common.Opensearch) (*opensearch.Client, error) {
	client, err := opensearch.NewClient(opensearch.Config{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Addresses: conf.Address,
		Username:  conf.Username,
		Password:  conf.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("数据源配置不正确: " + err.Error())
	}

	return client, err
}

func (open *Opensearch) SearchRequest(content string) (*opensearchapi.Response, error) {
	search := &opensearchapi.SearchRequest{
		Body: strings.NewReader(content),
	}

	return search.Do(context.Background(), open.Client)
}

func (open *Opensearch) IndicesDeleteRequest(indexNames []string) (*opensearchapi.Response, error) {
	del := &opensearchapi.IndicesDeleteRequest{
		Index: indexNames,
	}

	return del.Do(context.Background(), open.Client)
}
