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
	CorporationID string

	CreatedAt time.Time
	UpdatedAt time.Time
}

type NotifyDeployMessage struct {
	UUID    string   `json:"uuid"`
	Title   string   `json:"title"`
	Channel []string `json:"channel"`
	Content []struct {
		Servers    []Servers `json:"servers"`
		RegionName string    `json:"region_name"`
		RegionID   int       `json:"region_id"`
	} `json:"content"`
}

type Servers struct {
	ID            int    `json:"id"`
	IP            string `json:"ip"`
	Env           string `json:"env"`
	Name          string `json:"name"`
	EnvID         int    `json:"env_id"`
	Project       int    `json:"project"`
	RemoteID      string `json:"remote_id"`
	CorporationID string `json:"corporation_id"`
}
