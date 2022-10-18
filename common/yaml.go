package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/loggie-io/loggie/pkg/control"
)

func LoggieOperatorPipeline(index, ip, filePath string, kafkaBroker []string) error {
	var piplineTemplate = `
pipelines:
- name: demo
  sources:
  - name: operator
    type: file
    addonMeta: true
    paths:
    - %slog/GameOperate.log
    fieldsUnderRoot: true
    fields:
      index: operator-%s
      ip: %s
  - name: DbManager
    type: file
    addonMeta: true
    paths:
    - /var/log/engine/DbManager.log
    fieldsUnderRoot: true
    fields:
      index: server-%s
      ip: %s
  - name: GameManager
    type: file
    addonMeta: true
    paths:
    - /var/log/engine/GameManager.log
    fieldsUnderRoot: true
    fields:
      index: server-%s
      ip: %s
  - name: GameServer
    type: file
    addonMeta: true
    paths:
    - /var/log/engine/GameServer.log
    fieldsUnderRoot: true
    fields:
      index: server-%s
      ip: %s
  - name: GateServer
    type: file
    addonMeta: true
    paths:
    - /var/log/engine/GateServer.log
    fieldsUnderRoot: true
    fields:
      index: %s
      ip: %s
  sink:
    type: kafka
    balance: roundRobin
    brokers: %s
    compression: gzip
    topic: logservice
    codec: {}
`

	if len(kafkaBroker) == 0 {
		Logger.Error("kafka broker is nil")
		return errors.New("kafka broker is nil")
	}

	broker, _ := json.Marshal(kafkaBroker)
	piplineTemplate = fmt.Sprintf(piplineTemplate, RemoteFilepath, index, ip, index, ip, index, ip, index, ip, index, ip, string(broker))

	conf := &control.PipelineConfig{}
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
	// conf.Pipelines[0].Sink

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
