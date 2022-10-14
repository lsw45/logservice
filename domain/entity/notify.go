package entity

import "time"

const (
	NotifyMsgTableName    = "notify_msg"
	DeployIngestTableName = "deploy_ingest_task"
)

type NotifyMsgModel struct {
	UUID      string `gorm:"primaryKey"`
	Title     string
	Msg       string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type DeployIngestModel struct {
	Id             int
	NotifyId       string
	Status         int
	GameIp         string
	FailedMsg      string
	Index          string
	Config         string
	RemoteFilePath string
	Env            string
	KafkaBroker    []string

	EnvId         int
	Project       int
	RegionID      int
	CorporationId string

	CreatedAt time.Time
	UpdatedAt time.Time
}

type NotifyDeployMessage struct {
	UUID    string `json:"uuid"`
	Title   string `json:"title"`
	Content []struct {
		Servers    []Servers `json:"servers"`
		RegionID   int       `json:"RegionID"`
		RegionName string    `json:"RegionName"`
	} `json:"content"`
}

type Servers struct {
	ID            int    `json:"Id"`
	IP            string `json:"Ip"`
	Env           string `json:"Env"`
	Name          string `json:"Name"`
	EnvID         int    `json:"EnvID"`
	Project       int    `json:"Project"`
	RemoteID      int    `json:"RemoteID"`
	CorporationID int    `json:"CorporationId"`
}
