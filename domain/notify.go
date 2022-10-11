package domain

import (
	"errors"
	"fmt"
	"log-ext/common"
	"log-ext/domain/dependency"
	"log-ext/domain/entity"
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
	if message.Title != "broadcast.region.update" {
		common.Logger.Error("domain error: notify message title error: " + message.Title)
		return errors.New("message title error: " + message.Title)
	}

	// 消息幂等性*？*
	exist, err := dsvc.depRepo.ExitsNotifyByUUId(message.UUID)
	if err != nil {
		common.Logger.Errorf("domain error: %s", err)
		return err
	}

	if exist {
		return nil
	}

	var tasks []*entity.DeployIngestModel
	for _, content := range message.Content {
		for _, server := range content.Servers {
			// 新建任务——上传文件
			tasks = append(tasks, &entity.DeployIngestModel{
				GameIp:        server.IP,
				Index:         fmt.Sprintf("operator-%v-%v-%v", content.RegionID, server.Project, server.EnvID), // 拼接规则：区服ID-项目ID-环境ID
				Status:        1,
				NotifyId:      message.UUID,
				Env:           server.Env,
				EnvId:         server.EnvID,
				Project:       server.Project,
				CorporationId: server.CorporationID,
				RegionID:      content.RegionID,
				KafkaBroker:   common.GetKafka().Broker,
			})
		}
	}

	// 持久化保存回调消息和部署采集器任务
	err = dsvc.depRepo.SaveNotifyMessage(message)
	if err != nil {
		common.Logger.Errorf("domain error: SaveNotifyMessage %s", err)
		return err
	}

	_, err = dsvc.depRepo.SaveDeployeIngestTask(tasks)
	if err != nil {
		common.Logger.Errorf("domain error: SaveDeployeIngestTask %s", err)
		return err
	}

	// 异步上传pipeline文件
	for _, task := range tasks {
		go dsvc.TunnelUploadIngest(task)
	}

	return nil
}

// 上传pipeline并启动采集器
func (dsvc *depolyService) TunnelUploadIngest(task *entity.DeployIngestModel) {
	err := common.LoggieOperatorPipeline(task.Index, task.GameIp, "../doc/pipelines.yml", task.KafkaBroker)
	if err != nil {
		return
	}

	err = dsvc.depTunnel.UploadFile("../doc/pipelines.yml", task.GameIp, task.Env)
	if err != nil {
		common.Logger.Errorf("domain error: upload file: %s", err)

		err = dsvc.depRepo.UpdateDeployeIngestTask([]int{task.Id}, 3)
		if err != nil {
			common.Logger.Errorf("domain error: UpdateDeployeIngestTask 3: %s", err)
			return
		}

		return
	}

	// 更新任务状态为上传文件成功
	err = dsvc.depRepo.UpdateDeployeIngestTask([]int{task.Id}, 2)
	if err != nil {
		common.Logger.Errorf("domain error: UpdateDeployeIngestTask 2: %s", err)
		return
	}

	// 启动采集器：开启rsyslog，日志滚动等
	err = dsvc.TunnelDeployIngestTask(task)
	if err != nil {
		common.Logger.Errorf("domain error: deploy ingest: %s", err)
		return
	}
	return
}

// 启动采集器
func (dsvc *depolyService) TunnelDeployIngestTask(task *entity.DeployIngestModel) error {
	var err error

	sucess, err := dsvc.depTunnel.ShellTask(task.EnvId, task.Project, task.CorporationId, task.GameIp, true)
	if err != nil {
		common.Logger.Errorf("domain error: shell task: %s", err)
		return err
	}
	if !sucess {
		common.Logger.Errorf("domain error: shell task failed: " + task.GameIp)
		return err
	}

	// 更新任务状态
	err = dsvc.depRepo.UpdateDeployeIngestTask([]int{task.Id}, 6)
	if err != nil {
		common.Logger.Errorf("domain error: UpdateDeployeIngestTask 6: %s", err)
		return err
	}

	return nil
}

// 检查部署任务
func (dsvc *depolyService) TunnelCheckTask() {

}
