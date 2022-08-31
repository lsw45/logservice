package adapter

import (
	"log-ext/domain/entity"
)

type TaskDto struct {
	Task   entity.Job
	Status int
}
