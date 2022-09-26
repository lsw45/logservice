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
// 上传文件：POST http://ops-dev.cocos.org/paas/tunnel/file/
type TunnelInfra interface {
	UploadFile(data *bytes.Buffer, boundary string) (*entity.TunnelUploadFileRes, error)
	ShellTask(data *entity.ShellTaskReq) (*entity.ShellTaskDeployResp, error)
	CheckTask(id string) (*entity.ShellTaskStateResp, error)
}

var (
	token          = "14e58ac5e45f4fefa924a040c581698d"
	shell_task     = []string{"POST", "http://ops-dev.cocos.org/paas/tunnel/task/?x-token=" + token}
	check_task     = []string{"GET", "http://ops-dev.cocos.org/paas/tunnel/task/1/"}
	upload_file    = []string{"POST", "http://ops-dev.cocos.org/paas/tunnel/file/?x-token=" + token}
	RemoteFilepath = "/opt/loggie"
)

type Tunnel struct {
	Client http.Client
}

func NewTunnelClient(conf common.Tunnel) TunnelInfra {
	tr := &http.Transport{
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: conf.InsecureSkipVerify},
		DisableKeepAlives: conf.DisableKeepAlives,
		IdleConnTimeout:   time.Second * time.Duration(conf.IdleConnTimeout),
	}

	return &Tunnel{
		Client: http.Client{Transport: tr, Timeout: time.Duration(conf.Timeout) * time.Second},
	}
}

// 上传文件
func (tu *Tunnel) UploadFile(data *bytes.Buffer, boundary string) (*entity.TunnelUploadFileRes, error) {
	req, err := http.NewRequest(upload_file[0], upload_file[1], data)
	if err != nil {
		common.Logger.Errorf("request error: %s", err)
		return nil, err
	}

	req.Header.Set("Content-Type", boundary)

	resp, err := tu.Client.Do(req)
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
		common.Logger.Errorf("infra error: shell上传文件失败", err)
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		common.Logger.Error("infra error: shell上传文件读取结果失败")
		return nil, err
	}

	if resp.StatusCode > 300 {
		common.Logger.Error("infra error: shell上传文件失败: StatusCode > 300")
		err = errors.Errorf("shell上传文件失败: %s。%+v", string(body), resp)
		return nil, err
	}

	var res *entity.TunnelUploadFileRes
	if err = json.Unmarshal(body, &res); err != nil {
		common.Logger.Errorf("infra error: 解析shell上传结果失败: %s", string(body))
		return nil, err
	}

	return res, nil
}

// 下发任务：采集器启动
func (tu *Tunnel) ShellTask(data *entity.ShellTaskReq) (*entity.ShellTaskDeployResp, error) {
	reqData, err := json.Marshal(data)
	if err != nil {
		common.Logger.Errorf("marshal error: %s", err)
		return nil, err
	}

	req, err := http.NewRequest(shell_task[0], shell_task[1], bytes.NewReader(reqData))
	if err != nil {
		common.Logger.Errorf("request error: %s", err)
		return nil, err
	}
	resp, err := tu.Client.Do(req)
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
		common.Logger.Error("infra error: 下发shell任务失败")
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		common.Logger.Error("infra error: 下发shell任务读取结果失败")
		return nil, err
	}

	if resp.StatusCode > 300 {
		common.Logger.Error("infra error: 下发shell任务失败: StatusCode > 300")
		err = errors.Errorf("下发shell任务失败: %s。%+v", string(body), resp)
		return nil, err
	}

	var result *entity.ShellTaskDeployResp
	err = json.Unmarshal(body, &result)
	if err != nil {
		common.Logger.Errorf("infra error: 解析下发shell任务结果失败: %s", string(body))
		return nil, err
	}

	return result, nil
}

// 检查任务：采集器启动情况
func (tu *Tunnel) CheckTask(id string) (*entity.ShellTaskStateResp, error) {
	resp, err := http.Get(check_task[1])
	if err != nil {
		common.Logger.Errorf("request error: %s", err)
		return nil, err
	}
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
		common.Logger.Error("infra error: 检查shell任务失败")
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		common.Logger.Error("infra error: 检查shell任务读取结果失败")
		return nil, err
	}

	var result *entity.ShellTaskStateResp
	err = json.Unmarshal(body, &result)
	if err != nil {
		common.Logger.Errorf("infra error: 解析检查shell任务结果失败: %s", string(body))
		return nil, err
	}

	return result, nil
}
