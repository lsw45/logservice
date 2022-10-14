package infra

import (
	"fmt"
	"log-ext/common"
	"testing"
)

var conf = common.Mysql{
	Port:     3306,
	Host:     "121.37.173.234",
	Username: "root",
	Password: "Cocos@2021",
	DataBase: "paas_logservicev2_dev",
}

func TestSaveNotifyMessage(t *testing.T) {
	db, _ := NewMysql(conf)
	exit, err := db.ExitsNotifyByUUId("8989898898")
	fmt.Println(err)
	fmt.Print(exit)
}

func TestUpdate(t *testing.T) {
	db, _ := NewMysql(conf)

	err := db.UpdateDeployeIngestTask([]int{1, 2, 3}, 2)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestRelease(t *testing.T) {
	db, _ := NewMysql(conf)

	err := db.ReleaseRegion(12)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(1)
}
