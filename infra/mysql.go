package infra

import (
	"fmt"
	"gorm.io/driver/mysql"
	"log-ext/common"
	"log-ext/domain/entity"
	"time"

	"gorm.io/gorm"
)

type MysqlInfra interface {
	GetUser(id int) (*entity.User, error)
	GetUserConfigName(ingestID, version string) (string, error)
}

type Mysql struct {
	DB *gorm.DB
}

func NewMysqlDB(conf common.Mysql) (*gorm.DB, error) {
	dbDSN := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=True&loc=Local",
		conf.Username, conf.Password, conf.Host, conf.Port, conf.DataBase)

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		DSN: dbDSN,
	}), &gorm.Config{
		// PreparedStmt creates prepared statements when executing any SQL and caches them to speed up future calls
		PrepareStmt: true,
		//Logger:      logger.Default.LogMode(logger.Info), // 开启打印sql
	})

	if err != nil {
		return gormDB, fmt.Errorf("数据源配置不正确: " + err.Error())
	}

	db, err := gormDB.DB()
	if err != nil {
		return gormDB, fmt.Errorf("gorm 获取数据库失败: " + err.Error())
	}

	if err = db.Ping(); err != nil {
		return gormDB, fmt.Errorf("数据库连接失败: " + err.Error())
	}

	// 最大连接数
	db.SetMaxOpenConns(100)
	// 闲置连接数
	db.SetMaxIdleConns(20)
	// 最大连接周期
	db.SetConnMaxLifetime(100 * time.Second)

	return gormDB, nil
}

func (cli *Mysql) GetUser(id int) (*entity.User, error) {
	var result entity.User

	err := cli.DB.Table(entity.UserTableName).Where("id = ?", id).Find(&result).Error

	return &result, err
}

func (cli *Mysql) GetUserConfigName(ingestID, version string) (string, error) {
	var configName string
	err := cli.DB.Table(entity.UserTableName).Select("config_name").
		Where("user_id = ? and version = ?", ingestID, version).
		Find(&configName).Error
	return configName, err
}
