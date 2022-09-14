package domain

import (
	"context"
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
	// 消息幂等性
	exits, err := dsvc.depRepo.ExitsNotifyByUUId(message.UUID)
	if err != nil {
		common.Logger.Errorf("domain error: %s", err)
		return err
	}

	if exits {
		return nil
	}

	_, err = dsvc.depRepo.SaveNotifyMessage(message)
	if err != nil {
		common.Logger.Errorf("domain error: %s", err)
		return err
	}

	// 异步上传采集器文件
	for _, content := range message.Content {
		for _, server := range content.Servers {
			for _, ip := range server.IPObj {
				go dsvc.TunnelUploadIngest(ip.IP)
			}
		}
	}

	return nil
}

// 上传采集器文件
func (nctl *depolyService) TunnelUploadIngest(ip string) {
	file, err := os.Open("xxx")
	if err != nil {
		common.Logger.Errorf("domain error: %s", err)
		return
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		common.Logger.Errorf("domain error: ioutil: %s", err)
		return
	}

	err = nctl.depTunnel.UploadFile(data, ip)
	if err != nil {
		common.Logger.Errorf("domain error: upload file: %s", err)
	}
	return
}

// 部署采集器
func (nctl *depolyService) TunnelDeployTask(ctx context.Context, ip, username, password string) {

}

// 检查部署任务
func (nctl *depolyService) TunnelCheckTask() {

}
