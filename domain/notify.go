package domain

import (
	"fmt"
	"io/ioutil"
	"log-ext/common"
	"log-ext/domain/dependency"
	"log-ext/domain/entity"
	"os"
)

type NotifyService interface {
	DeployNotify(message *entity.NotifyDeployMessage) error
}

func NewDeployNotifyService(depRepo dependency.MysqlRepo, depTunnel dependency.TunnelRepo) NotifyService {
	return &depolyService{
		depRepo:   depRepo,
		depTunnel: depTunnel,
	}
}

type depolyService struct {
	depRepo   dependency.MysqlRepo
	depTunnel dependency.TunnelRepo
}

func (dsvc *depolyService) DeployNotify(message *entity.NotifyDeployMessage) error {
	// 消息幂等性*？*
	exits, err := dsvc.depRepo.ExitsNotifyByUUId(message.UUID)
	if err != nil {
		common.Logger.Errorf("domain error: %s", err)
		return err
	}

	if exits {
		return nil
	}

	var tasks []*entity.DeployIngestTable
	var ips []string
	for _, content := range message.Content {
		for _, server := range content.Servers {
			// 新建任务——上传文件
			tasks = append(tasks, &entity.DeployIngestTable{
				Ip:       server.IP,
				Index:    fmt.Sprintf("%v-%v-%v", content.RegionID, server.Project, server.EnvID),
				Status:   1,
				NotifyId: message.UUID,
			})
			ips = append(ips, server.IP)
		}
	}

	// 异步上传采集器文件
	go dsvc.TunnelUploadIngest(ips)

	// 持久化保存回调消息和部署采集器任务
	err = dsvc.depRepo.SaveDeployeIngestTask(tasks)
	if err != nil {
		common.Logger.Errorf("domain error: SaveDeployeIngestTask %s", err)
		return err
	}

	err = dsvc.depRepo.SaveNotifyMessage(message)
	if err != nil {
		common.Logger.Errorf("domain error: SaveNotifyMessage %s", err)
		return err
	}

	return nil
}

// 上传采集器并启动采集器
func (dsvc *depolyService) TunnelUploadIngest(ip []string) {
	// 上传采集器
	file, err := os.Open("xxx")
	if err != nil {
		common.Logger.Errorf("domain error: open file: %s", err)
		return
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		common.Logger.Errorf("domain error: ioutil: %s", err)
		return
	}

	err = dsvc.depTunnel.UploadFile(data, ip)
	if err != nil {
		common.Logger.Errorf("domain error: upload file: %s", err)
		return
	}

	// 启动采集器
	err = dsvc.TunnelDeployIngestTask(ip)
	if err != nil {
		common.Logger.Errorf("domain error: deploy ingest: %s", err)
		return
	}
	return
}

// 启动采集器
func (dsvc *depolyService) TunnelDeployIngestTask(ip []string) error {

	return nil
}

// 检查部署任务
func (dsvc *depolyService) TunnelCheckTask() {

}
