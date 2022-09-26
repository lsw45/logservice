package common

import (
	"testing"
)

func TestWriteYaml(t *testing.T) {
	WriteLoggiePipeline("x", "x", "./pip.yml")
}
