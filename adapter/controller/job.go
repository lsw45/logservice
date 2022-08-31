package controller

import "log-ext/domain"

type JobController struct {
	svc *domain.JobService
}

func NewJobController(svc *domain.JobService) *JobController {
	return &JobController{}
}

func (job *JobController) Get() {

}
