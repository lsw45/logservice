package repository

import (
	"fmt"
	"log-ext/domain/dependency"
)

var _ dependency.OpensearchRepo = (*OpensearchRepo)(nil)

type OpensearchRepo struct {
	//Infra infra.OpensearchInfra
}

func (o *OpensearchRepo) Filter() {

}

func (o *OpensearchRepo) ListLog() {
	fmt.Println("xxx")
}

func (o *OpensearchRepo) Count() {
	fmt.Println("xxx")
}
