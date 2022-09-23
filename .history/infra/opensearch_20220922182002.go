package infra

import (
	"context"
	"crypto/tls"
	"fmt"
	"log-ext/common"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
)

var _ OpensearchInfra = &openDB{}

type OpensearchInfra interface {
	SearchRequest(indexNames []string, content string) (*opensearchapi.Response, error)
	IndicesDeleteRequest(indexNames []string) (*opensearchapi.Response, error)
}

type openDB struct {
	Client *opensearch.Client
}

func NewOpensearch(conf common.Opensearch) (*openDB, error) {
	client, err := opensearch.NewClient(opensearch.Config{
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Second,
			DialContext:           (&net.Dialer{Timeout: 5 * time.Second}).DialContext,
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: conf.InsecureSkipVerify},
		},
		Addresses: conf.Address,
		Username:  conf.Username,
		Password:  conf.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("数据源配置不正确: %v", err.Error())
	}

	res, err := client.Info()
	if err != nil {
		return nil, fmt.Errorf("获取信息失败: %v", err.Error())
	}
	defer res.Body.Close()

	return &openDB{client, err
}

func (open *openDB) SearchRequest(indexNames []string, content string) (*opensearchapi.Response, error) {
	search := &opensearchapi.SearchRequest{
		Index: indexNames,
		Body:  strings.NewReader(content),
	}

	resp, err := search.Do(context.Background(), open.Client)
	if err != nil {
		common.Logger.Errorf("infra opensearch search error: %+v", err)
		return nil, err
	}

	return resp, err
}

func (open *openDB) IndicesDeleteRequest(indexNames []string) (*opensearchapi.Response, error) {
	del := &opensearchapi.IndicesDeleteRequest{
		Index: indexNames,
	}

	resp, err := del.Do(context.Background(), open.Client)
	if err != nil {
		common.Logger.Errorf("infra opensearch delete error: %+v", err)
	}

	return resp, err
}
