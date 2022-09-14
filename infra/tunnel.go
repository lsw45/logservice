package infra

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"log-ext/common"
	"log-ext/domain/entity"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

// 提交执行任务：POST http://ops-dev.cocos.org/paas/tunnel/task/
// 获取任务详情信息：GET http://ops-dev.cocos.org/paas/tunnel/task/1/
// 上传文件：POST http://ops-dev.cocos.org/paas/tunnel/file/?x-token=8155214abb3206a0ef2e18d6bae586b0  Content-Type: multipart/form-data; boundary=<calculated when request is sent>
type TunnelInfra interface {
	UploadFile(data *bytes.Buffer) (*entity.TunnelUploadFileRes, error)
	DeployTask()
	CheckTask()
}

type tunnel struct {
	Url           string
	Method        string
	Boundary      string
	Authorization string
}

func NewTunnelClient(url, method, boundary, authorization string) TunnelInfra {
	return &tunnel{
		Url:           url,
		Method:        method,
		Boundary:      boundary,
		Authorization: authorization,
	}
}

func (tu *tunnel) UploadFile(data *bytes.Buffer) (*entity.TunnelUploadFileRes, error) {
	tr := &http.Transport{
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
		DisableKeepAlives: true,
		IdleConnTimeout:   time.Second * 60,
	}
	client := http.Client{Transport: tr, Timeout: 20 * time.Second}

	req, err := http.NewRequest(tu.Method, tu.Url, data)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if resp != nil {
		defer func() {
			err := resp.Body.Close()
			if err != nil {
				err = errors.Wrap(err, "连接关闭错误")
				common.Logger.Errorf("infra error: %v", err)
			}
		}()
	}
	if err != nil {
		common.Logger.Error("infra error: shell上传文件失败")
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		common.Logger.Error("infra error: shell上传文件读取结果失败")
		return nil, err
	}

	var res *entity.TunnelUploadFileRes
	if err = json.Unmarshal(body, &res); err != nil {
		common.Logger.Errorf("infra error: 解析shell上传结果失败: %s", string(body))
		return nil, err
	}

	return res, nil
}

func (tunnel *tunnel) DeployTask() {}

func (tunnel *tunnel) CheckTask() {}
