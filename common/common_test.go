package common

import (
	"testing"
)

var piptest = `
pipelines:
- name: demo
  sources:
  - name: mylog
    type: file
    addonMeta: true
    paths:
    - /home/logservice2/log/GameStatistic.log
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
`

func TestWriteYaml(t *testing.T) {
	WriteLoggiePipeline("x", "x", "./pip.yml")
}
