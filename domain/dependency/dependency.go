package dependency

//go:generate mockgen -source ../dependency/dependency.go -destination ../../mock/mock_dependency.go -package mock
type OpensearchRepo interface {
	Count()
	ListLog()
	Filter()
}
