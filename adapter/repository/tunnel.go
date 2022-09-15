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

func (t *TunnelRepo) UploadFile(fileData []byte, ip string) error {
	buf := &bytes.Buffer{}
	writer := multipart.NewWriter(buf)

	field := &entity.UpdateFileReq{
		Remote:   "/usr/local",
		Server:   ip,
		Preserve: true,
		File:     fileData,
	}

	writer.WriteField("remote", field.Remote)
	writer.WriteField("preserve", strconv.FormatBool(field.Preserve))
	writer.WriteField("server", field.Server)
	boundary := writer.FormDataContentType()

	common.Logger.Infof("上传文件开始: %s", field.Server)
	upwriter, _ := writer.CreateFormFile("file", "loggie")

	// formFile, err := writer.CreateFormField("file")
	// if err != nil {
	// 	common.Logger.Errorf("创建form文件失败,field: %+v", field)
	// 	return err
	// }

	_, err := io.Copy(upwriter, bytes.NewReader(fileData))
	if err != nil {
		common.Logger.Errorf("复制文件字段失败,field: %s", field.Server)
		return err
	}

	err = writer.Close()
	if err != nil {
		common.Logger.Errorf("关闭writer缓冲失败,field: %s", field.Server)
		return err
	}

	res, err := t.TunnelInfra.UploadFile(buf, boundary)
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

// func (t *TunnelRepo) UploadFile(fileData []byte, ip string) error {
// 	data := &bytes.Buffer{}
// 	writer := multipart.NewWriter(data)

// 	field := &entity.UpdateFileReq{
// 		Remote:   "/usr/local",
// 		Server:   ip,
// 		Preserve: true,
// 		File:     fileData,
// 	}
// 	common.Logger.Infof("上传文件开始: %s", field.Server)

// 	formFile, err := writer.CreateFormField("file")
// 	if err != nil {
// 		common.Logger.Errorf("创建form文件失败,field: %+v", field)
// 		return err
// 	}

// 	_, err = io.Copy(formFile, bytes.NewReader(fileData))
// 	if err != nil {
// 		common.Logger.Errorf("复制文件字段失败,field: %+v", field.File)
// 		return err
// 	}

// 	err = writer.WriteField("remote", field.Remote)
// 	if err != nil {
// 		common.Logger.Errorf("写入form字段失败,field: %+v", field.Remote)
// 		return err
// 	}

// 	err = writer.WriteField("preserve", strconv.FormatBool(field.Preserve))
// 	if err != nil {
// 		common.Logger.Errorf("写入form字段失败,field: %+v", field.Preserve)
// 		return err
// 	}

// 	err = writer.WriteField("server", field.Server)
// 	if err != nil {
// 		common.Logger.Errorf("写入form字段失败,field: %+v", field.Server)
// 		return err
// 	}

// 	err = writer.Close()
// 	if err != nil {
// 		common.Logger.Errorf("关闭writer缓冲失败,field: %+v", field)
// 		return err
// 	}

// 	res, err := t.TunnelInfra.UploadFile(data)
// 	if err != nil {
// 		common.Logger.Errorf("shell上传文件失败field: %+v", field)
// 		return err
// 	}

// 	if res.Code != 0 {
// 		common.Logger.Errorf("shell上传文件失败: %+v", field)
// 		err = errors.Errorf("shell上传文件失败field: %+v", res)
// 		return err
// 	}

// 	common.Logger.Infof("上传文件结束: %s", field.Server)

// 	return nil
// }

func (t *TunnelRepo) ShellTask(env, project int, corporationId, server string, async bool) (bool, error) {
	command := []string{}
	params := entity.ShellParams{
		Shell:            "/bin/bash",
		Server:           server,
		Command:          command,
		ShellEnvironment: map[string]interface{}{"TUNNEL_TEST": "02"},
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
		common.Logger.Errorf("下发shell任务失败: %+v", err)
		return false, err
	}
	if resp.Code != 0 {
		common.Logger.Errorf("下发shell任务失败: %s", server)
		err = errors.Errorf("下发shell任务失败: %+v", resp)
		return false, err
	}
	return true, nil
}

func (t *TunnelRepo) CheckTask(id string) (*entity.ShellTaskStateResp, error) {

	return nil, nil
}
