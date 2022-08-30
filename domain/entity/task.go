package entity

import "time"

type TaskConfig struct {
	Id       string
	Input    Input
	Output   Output
	Interval time.Duration
	Type     string
	Tag      []string
	Status   int // 任务状态
}

type Input struct {
	Paths []string
}

type Output struct {
	Hosts string
	Port  string
}
