package common

import (
	"errors"
	"os"

	"github.com/goccy/go-yaml"

	"github.com/loggie-io/loggie/pkg/control"
)

var piplineTemplate = `
pipelines:
  - name: demo
    sources:
      - type: file
        name: mylog
        addonMeta: true
        paths:
          - "/var/log/kdump.log"
    interceptors:
    - type: transformer
      actions:
        - if: hasPrefix(body, {)
          then:
            - action: jsonDecode(body)
            - action: del(stream)
            - action: add(topic, json)
          else:
            - action: add(topic, plain)
            
    sink:
      type: kafka
      brokers: ["10.0.3.116:9093"]
      topic: "test"
      balance: "roundRobin"
      compression: "gzip"
`

func WriteLoggiePipeline(index, ip, filePath string) error {
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
