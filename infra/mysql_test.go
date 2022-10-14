package infra

import (
	"fmt"
	"log-ext/common"
	"testing"
)

func TestSaveNotifyMessage(t *testing.T) {
	conf := common.Mysql{
		Port:     3306,
		Host:     "121.37.173.234",
		Username: "root",
		Password: "Cocos@2021",
		DataBase: "paas_logservicev2_dev",
	}

	db, _ := NewMysql(conf)
	exit, err := db.ExitsNotifyByUUId("8989898898")
	fmt.Println(err)
	fmt.Print(exit)
}
