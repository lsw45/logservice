package domain

type job interface {
	Get()
	Add()
	List()
	Update()
	Delete()
}

type JobService struct {
}
