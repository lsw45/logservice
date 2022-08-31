package repository

type search interface {
	Search()
	List()
}

type Opensearch struct {
}

func (o *Opensearch) Search() {

}

func (o *Opensearch) List() {

}
