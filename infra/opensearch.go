package infra

type OpensearchInfra interface {
	List()
	Count()
}

type Opensearch struct {

}