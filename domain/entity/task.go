package entity

import "time"

type TaskConfig struct {
	Id       string
	Input    Input
	Output   Output
	Interval time.Duration
	Type     string
	Tag      []string
}

type Input struct {
	Paths []string
}

type Output struct {
	Hosts string
	Port  string
}
