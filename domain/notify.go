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

	// 消息有重试机制，所以要防止重复处理
	exist, err := dsvc.depRepo.ExitsNotifyByUUId(message.UUID)
	if err != nil {
		common.Logger.Errorf("domain error: %s", err)
		return err
	}

	if exist == "exist" || exist == entity.DeployIngestTableName {
		return nil
	}

	// 持久化保存回调消息
	if exist == entity.NotifyMsgTableName {
		mo := &entity.NotifyMsgModel{
			UUID:  message.UUID,
			Title: message.Title,
		}

		err = dsvc.depRepo.SaveNotifyMessage(mo)
		if err != nil {
			common.Logger.Errorf("domain error: SaveNotifyMessage %s", err)
			return err
		}
	}

	// servers为空，是服务器释放阶段发来的
	var tasks []entity.DeployIngestModel
	for _, content := range message.Content {
		for _, server := range content.Servers {
			// 新建任务——上传文件
			tasks = append(tasks, entity.DeployIngestModel{
				GameIp:        server.IP,
				Index:         fmt.Sprintf("%v-%v-%v", server.Project, server.EnvID, content.RegionID), // 拼接规则：项目ID-环境ID-区服ID
				Status:        1,
				NotifyId:      message.UUID,
				Env:           server.Env,
				EnvId:         server.EnvID,
				Project:       server.Project,
				CorporationID: server.CorporationID,
				RegionID:      content.RegionID,
				KafkaBroker:   common.GetKafka().Broker,
			})
		}
	}

	// 释放采集任务
	if len(tasks) == 0 {
		err = dsvc.depRepo.ReleaseRegion(message.Content[0].RegionID)
		if err != nil {
			common.Logger.Errorf("domain error: ReleaseRegion %s", err)
			return err
		}
		return nil
	}

	// 保存部署采集器任务
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
func (dsvc *depolyService) TunnelUploadIngest(task entity.DeployIngestModel) {
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
	}
}

// 启动采集器
func (dsvc *depolyService) TunnelDeployIngestTask(task entity.DeployIngestModel) error {
	success, err := dsvc.depTunnel.ShellTask(task.EnvId, task.Project, task.CorporationID, task.GameIp, true)
	if err != nil {
		common.Logger.Errorf("domain error: shell task: %s", err)
		return err
	}
	if !success {
		common.Logger.Errorf("domain error: shell task failed: " + task.GameIp)
		return err
	}

	// 更新任务状态
	err = dsvc.depRepo.UpdateDeployeIngestTask([]int{task.Id}, 6)
	if err != nil {
		common.Logger.Errorf("domain error: UpdateDeployeIngestTask 6: %s", err)
		return err
	}

	dsvc.depTunnel.CheckTask(task.Id)

	return nil
}

// TODO:检查部署任务
func (dsvc *depolyService) TunnelCheckTask(id int) error {
	resp, err := dsvc.depTunnel.CheckTask(id)
	if err != nil {
		common.Logger.Errorf("domain error: TunnelCheckTask 6: %s", err)
		return err
	}

	if resp.Data.Status == "ERROR" {
		err = dsvc.depRepo.UpdateDeployeIngestTask([]int{id}, 8)
		if err != nil {
			common.Logger.Errorf("domain error: UpdateDeployeIngestTask 8: %s", err)
			return err
		}
	}

	for _, detail := range resp.Data.Result {
		if detail.Status != "SUCCESS" {
			err = dsvc.depRepo.UpdateDeployeIngestTask([]int{id}, 8)
			if err != nil {
				common.Logger.Errorf("domain error: UpdateDeployeIngestTask 8: %s", err)
				return err
			}
			return nil
		}
		for _, shell := range detail.Detail {
			if shell.Exited != 0 {
				err = dsvc.depRepo.UpdateDeployeIngestTask([]int{id}, 8)
				if err != nil {
					common.Logger.Errorf("domain error: UpdateDeployeIngestTask 8: %s", err)
					return err
				}
				return nil
			}
		}
	}

	err = dsvc.depRepo.UpdateDeployeIngestTask([]int{id}, 7)
	if err != nil {
		common.Logger.Errorf("domain error: UpdateDeployeIngestTask 7: %s", err)
		return err
	}

	return nil
}
