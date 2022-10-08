package common

import (
	"errors"
	"fmt"
	"os"

	"github.com/goccy/go-yaml"

	"github.com/loggie-io/loggie/pkg/control"
	// "github.com/loggie-io/loggie/pkg/core/interceptor"
)

var piplineTemplate = `
pipelines:
- name: demo
  sources:
  - name: mylog
    type: file
    addonMeta: true
    paths:
    - %s/log/GameOperate.log
    fieldsUnderRoot: true
    fields:
      index: x
      ip: x
  sink:
    type: kafka
    balance: roundRobin
    brokers:
    - 10.0.3.116:9093
    compression: gzip
    topic: test
    codec: {}
`

func WriteLoggiePipeline(index, ip, filePath string) error {
	conf := &control.PipelineConfig{}
	piplineTemplate = fmt.Sprintf(piplineTemplate, RemoteFilepath)
	err := yaml.Unmarshal([]byte(piplineTemplate), conf)
	if err != nil {
		Logger.Errorf("unmarshal yaml failed: %+v", err)
		return err
	}

	// 设置index和ip等自定义字段
	if len(conf.Pipelines) == 0 || len(conf.Pipelines[0].Sources) == 0 {
		Logger.Error("empty pipeline")
		return errors.New("empty pipeline")
	}
	conf.Pipelines[0].Sources[0].FieldsUnderRoot = true
	conf.Pipelines[0].Sources[0].Fields = map[string]interface{}{"index": index, "ip": ip}

	// conf.Pipelines[0].Interceptors
	// interceptors := make([]interceptor.Config,3)
	// interceptor := interceptor.Config{type:"transform"}

	pipe, err := yaml.Marshal(conf)
	if err != nil {
		Logger.Errorf("loggie pipeline config error: %+v", err)
		return err
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_RDWR|os.O_TRUNC, 0744)
	if err != nil {
		Logger.Errorf("open file failed: %+v", err)
		return err
	}
	_, err = file.Write(pipe)

	return err
}