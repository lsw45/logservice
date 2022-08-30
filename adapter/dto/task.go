package adapter

import (
	"../../domain/entity"
)

type TaskDto struct {
	Task   entity.TaskConfig
	Status int
}
