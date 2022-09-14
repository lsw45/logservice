package repository

import (
	"bytes"
	"io"
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
	return &TunnelRepo{defaultInfra.TunnelRepo}
}

func (t *TunnelRepo) UploadFile(fileData []byte, ip string) error {
	data := &bytes.Buffer{}
	writer := multipart.NewWriter(data)

	field := entity.UpdateFileReq{
		Remote:   "xxx",
		Server:   ip,
		Preserve: true,
	}

	formFile, err := writer.CreateFormField("file")
	if err != nil {
		err = errors.Wrapf(err, "创建form文件失败,field: %+v", field)
		return err
	}

	_, err = io.Copy(formFile, bytes.NewReader(fileData))
	if err != nil {
		err = errors.Wrapf(err, "复制文件字段失败,field: %+v", field)
		return err
	}

	err = writer.WriteField("remote", field.Remote)
	if err != nil {
		err = errors.Wrapf(err, "写入form字段失败,field: %+v", field)
		return err
	}

	err = writer.WriteField("server", field.Server)
	if err != nil {
		err = errors.Wrapf(err, "写入form字段失败,field: %+v", field)
		return err
	}

	err = writer.WriteField("preserve", strconv.FormatBool(field.Preserve))
	if err != nil {
		err = errors.Wrapf(err, "写入form字段失败,field: %+v", field)
		return err
	}

	err = writer.Close()
	if err != nil {
		err = errors.Wrapf(err, "关闭writer缓冲失败,field: %+v", field)
		return err
	}

	res, err := t.TunnelInfra.UploadFile(data)
	if err != nil {
		err = errors.Wrapf(err, "shell上传文件失败field: %+v", field)
		return err
	}
	if res.Code != 0 {
		err = errors.Wrapf(err, "shell上传文件失败field: %+v", field)
		return err
	}
	return nil
}
