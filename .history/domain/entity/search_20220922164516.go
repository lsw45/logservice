package entity

import (
	"time"

	"github.com/olivere/elastic/v7"
)

type QueryDocs struct {
	From      int
	Size      int
	EnableDSL bool
	Sort      []elastic.Sorter
}

type CommonResp struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
}
type LogsFilterReq struct {
	Indexs          []string    `json:"indexs"`
	Env             string      `json:"env"`
	Project         int         `json:"project"`
	Limit           int         `json:"limit"`
	Page            int         `json:"page"`
	PageSize        int         `json:"page_size"`
	LineNum         interface{} `json:"line_num"`
	Keywords        string      `json:"keywords"`
	StartTime       int64       `json:"start_time"`
	EndTime         int64       `json:"end_time"`
	IsDesc          bool        `json:"is_desc"`
	RegionVal       int         `json:"region_val"`
	RegionServerVal string      `json:"region_server_val"`
	Date            []time.Time `json:"date"`
	Sort 
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

type LogsFilterResp struct {
	CommonResp
	Data struct {
		Results string `json:"results"`
		Count   int    `json:"count"`
		Total   int    `json:"total"`
	} `json:"data"`
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
	Data struct {
		Results []HistogramResult `json:"results"`
		Count   int               `json:"count"`
	} `json:"data"`
}
