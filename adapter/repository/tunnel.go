package repository

import (
	"bytes"
	"io"
	"log-ext/common"
	"log-ext/domain/entity"
	"log-ext/infra"
	"mime/multipart"
	"strconv"

	"github.com/pkg/errors"
)

type TunnelRepo struct {
	infra.TunnelInfra
}

func NewTunnelRepo() *TunnelRepo {
	return &TunnelRepo{defaultInfra.Tunnel}
}

func (t *TunnelRepo) UploadFile(fileData []byte, ips []string) error {
	data := &bytes.Buffer{}
	writer := multipart.NewWriter(data)

	var field entity.UpdateFileReq

	// 每个请求，都尽量成功，且超时重试
	for _, ip := range ips {
		field = entity.UpdateFileReq{
			Remote:   "/usr/local",
			Server:   ip,
			Preserve: true,
		}
	}

	formFile, err := writer.CreateFormField("file")
	if err != nil {
		common.Logger.Errorf("创建form文件失败,field: %+v", field)
		return err
	}

	_, err = io.Copy(formFile, bytes.NewReader(fileData))
	if err != nil {
		common.Logger.Errorf("复制文件字段失败,field: %+v", field)
		return err
	}

	err = writer.WriteField("remote", field.Remote)
	if err != nil {
		common.Logger.Errorf("写入form字段失败,field: %+v", field)
		return err
	}

	err = writer.WriteField("server", field.Server)
	if err != nil {
		common.Logger.Errorf("写入form字段失败,field: %+v", field)
		return err
	}

	err = writer.WriteField("preserve", strconv.FormatBool(field.Preserve))
	if err != nil {
		common.Logger.Errorf("写入form字段失败,field: %+v", field)
		return err
	}

	err = writer.Close()
	if err != nil {
		common.Logger.Errorf("关闭writer缓冲失败,field: %+v", field)
		return err
	}

	res, err := t.TunnelInfra.UploadFile(data)
	if err != nil {
		common.Logger.Errorf("shell上传文件失败field: %+v", field)
		return err
	}
	if res.Code != 0 {
		err = errors.Wrapf(err, "shell上传文件失败field: %+v", field)
		return err
	}
	return nil
}

func (t *TunnelRepo) ShellTask(env, project int, corporationId string, async bool) (bool, error) {
	command := []string{}
	params := entity.ShellParams{
		Shell:            "/bin/bash",
		Server:           "",
		ShellEnvironment: map[string]interface{}{"TUNNEL_TEST": "02"},
		Command:          command,
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
		common.Logger.Errorf("下发shell任务失败: %s", err)
		return false, err
	}
	if resp.Code != 0 {
		common.Logger.Errorf("下发shell任务失败: %s", err)
		return false, err
	}
	return true, nil
}

func (t *TunnelRepo) CheckTask(id string) (*entity.ShellTaskStateResp, error) {

	return nil, nil
}
