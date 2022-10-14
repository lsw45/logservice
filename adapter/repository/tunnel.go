package repository

import (
	"bytes"
	"io"
	"log-ext/common"
	"log-ext/domain/entity"
	"log-ext/infra"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"

	"github.com/pkg/errors"
)

type TunnelRepo struct {
	infra.TunnelInfra
	infra.MysqlInfra
}

func NewTunnelRepo() *TunnelRepo {
	return &TunnelRepo{defaultRepo.Tunnel, defaultRepo.Mysql}
}

func (t *TunnelRepo) UploadFile(file_path string, ip, env string) error {
	// 上传pipeline

	buf := &bytes.Buffer{}
	writer := multipart.NewWriter(buf)

	field := &entity.UpdateFileReq{
		Remote:   common.RemoteFilepath,
		Server:   ip,
		Preserve: true,
	}

	writer.WriteField("remote", field.Remote)
	writer.WriteField("preserve", strconv.FormatBool(field.Preserve))
	writer.WriteField("server", field.Server)

	common.Logger.Infof("上传文件开始: %s", field.Server)

	file, err := os.Open(file_path)
	if err != nil {
		common.Logger.Errorf("domain error: open file: %s", err)
		return err
	}
	defer file.Close()

	upwriter, err := writer.CreateFormFile("file", filepath.Base(file_path))
	if err != nil {
		common.Logger.Errorf("domain error: CreateFormFile: %s", err)
		return err
	}

	_, err = io.Copy(upwriter, file)
	if err != nil {
		common.Logger.Errorf("复制文件字段失败,field: %s", field.Server)
		return err
	}

	err = writer.Close()
	if err != nil {
		common.Logger.Errorf("关闭writer缓冲失败,field: %s", field.Server)
		return err
	}

	res, err := t.TunnelInfra.UploadFile(buf, writer.FormDataContentType())
	if err != nil {
		common.Logger.Errorf("shell上传文件失败field: %s", field.Server)
		return err
	}

	if res.Code != 0 {
		common.Logger.Errorf("shell上传文件失败: %s", field.Server)
		err = errors.Errorf("shell上传文件失败field: %+v", res)
		return err
	}

	common.Logger.Infof("上传文件结束: %s", field.Server)

	return nil
}

func (t *TunnelRepo) ShellTask(env, project int, corporationId, server string, async bool) (bool, error) {
	common.Logger.Info("开发启动采集器任务")
	command := []string{
		"cd $workspace;wget https://for-frontend-imagefile.obs.cn-east-3.myhuaweicloud.com/logservice2/loggie.yml",
		"cd $workspace;wget https://for-frontend-imagefile.obs.cn-east-3.myhuaweicloud.com/logservice2/loggie;chmod +x loggie",
		"cd $workspace;wget https://for-frontend-imagefile.obs.cn-east-3.myhuaweicloud.com/logservice2/start.sh;chmod +x start.sh",
		"cd $workspace;wget https://for-frontend-imagefile.obs.cn-east-3.myhuaweicloud.com/logservice2/openlog.sh;chmod +x openlog.sh",
		"source $workspace/openlog.sh",
		"source $workspace/start.sh",
	}
	params := entity.ShellParams{
		Shell:            "/bin/bash",
		Server:           server,
		Command:          command,
		ShellEnvironment: map[string]interface{}{"workspace": common.RemoteFilepath},
	}
	reqData := &entity.ShellTaskReq{
		Env:           env,
		Params:        []entity.ShellParams{params},
		Project:       project,
		Asynchronous:  async,
		CorporationID: corporationId,
	}

	resp, err := t.TunnelInfra.ShellTask(reqData)
	if err != nil {
		common.Logger.Error(err)
		return false, err
	}
	if resp.Code != 0 {
		common.Logger.Errorf("下发shell任务失败: %s", server)
		err = errors.Errorf("下发shell任务失败: %+v", resp)
		return false, err
	}

	common.Logger.Infof("采集任务执行id: %v", resp.Data.ID)

	// Todo 保存shell task的下发记录，并与文件下发任务记录做关联
	// t.MysqlInfra.SaveShellTask()

	common.Logger.Info("下发启动任务成功")
	return true, nil
}

func (t *TunnelRepo) CheckTask(id string) (*entity.ShellTaskStateResp, error) {

	return nil, nil
}
