package entity

import (
	"time"

	"github.com/olivere/elastic/v7"
)

type DateHistogram struct {
	Field     string
	Interval  string
	GroupName string
	StartTime time.Time
	EndTime   time.Time
}

type QueryDocs struct {
	From      int
	Size      int
	EnableDSL bool
	Query     string
	Sort      []elastic.Sorter
}

type CommonResp struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
}

type LogsFilterReq struct {
	Page      int             `json:"page"`
	Limit     int             `json:"limit"`
	PageSize  int             `json:"page_size"`
	Keywords  string          `json:"keywords"`
	Indexs    []string        `json:"indexs"`
	StartTime int64           `json:"start_time"`
	EndTime   int64           `json:"end_time"`
	Sort      map[string]bool `json:"sort"`
}

type AggregationReq struct {
	Indexs []string `json:"indexs"`
	Aggs   string   `json:"aggs"`
}

type DateHistogramReq struct {
	Indexs    []string `json:"indexs"`
	Interval  int64    `json:"interval"`
	GroupName string   `json:"group_name"`
	StartTime int64    `json:"start_time"`
	EndTime   int64    `json:"end_time"` //second
}

type LogsFilter struct {
	LogsFilterReq
}

type OpenResp struct {
	Took     int  `json:"took"`
	TimedOut bool `json:"timed_out"`
	Hits     Hits `json:"hits"`
	Shards   struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Skipped    int `json:"skipped"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
}

type Hits struct {
	MaxScore int     `json:"max_score"`
	Hits     []Hits2 `json:"hits"`
	Total    struct {
		Value    int    `json:"value"`
		Relation string `json:"relation"`
	} `json:"total"`
}

type Hits2 struct {
	Index  string                 `json:"_index"`
	ID     string                 `json:"_id"`
	Score  int                    `json:"_score"`
	Source map[string]interface{} `json:"_source"`
}

type LogsResult struct {
	Content string `json:"content"`
	LineNum string `json:"line_num"`
	Labels  Labels `json:"labels"`
}

type AggregationResp struct {
	CommonResp
	Data *elastic.SearchResult `json:"data"`
}

type LogsFilterResp struct {
	CommonResp
	Data struct {
		Results *elastic.SearchHits `json:"results"`
		Count   int64               `json:"count"`
	} `json:"data"`
}

type NearbyDocResp struct {
	CommonResp
	Data []*elastic.SearchHit `json:"data"`
}

type Labels struct {
	HostName      string `json:"hostName"`
	HostIP        string `json:"hostIP"`
	AppName       string `json:"appName"`
	ContainerName string `json:"containerName"`
	ClusterName   string `json:"clusterName"`
	HostID        string `json:"hostId"`
	PodName       string `json:"podName"`
	ClusterID     string `json:"clusterId"`
	NameSpace     string `json:"nameSpace"`
	Time          string `json:"time"`
	PathFile      string `json:"pathFile"`
	Category      string `json:"category"`
}

type HistogramResult struct {
	Num       int   `json:"num"`
	StartTime int64 `json:"startTime"`
	EndTime   int64 `json:"endTime"`
}

type HistogramResp struct {
	CommonResp
	Data  []BucketsList `json:"data"`
	Count int64     `json:"count"`
}

type BucketsList struct {
	DocCount  int   `json:"doc_count"`
	StartTime float64 `json:"start_time"`
	EndTime   float64 `json:"end_time"`
}

type Buckets struct {
	DocCount    int         `json:"doc_count"`
	KeyAsString string      `json:"key_as_string"`
	Key         interface{} `json:"key"` //s
}

type DateHistAggre struct {
	Buckets []Buckets `json:"buckets"`
}
