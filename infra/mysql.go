package infra

import (
	"fmt"
	"log-ext/common"
	"log-ext/domain/entity"
	"time"

	"gorm.io/driver/mysql"

	"gorm.io/gorm"
)

var _ MysqlInfra = new(mysqlDB)

type MysqlInfra interface {
	CommonInfra
	GetUser(id int) (*entity.User, error)
	GetUserConfigName(ingestID, version string) (string, error)
	ExitsNotifyByUUId(uuid string) (string, error)
	SaveNotifyMessage(msg *entity.NotifyMsgModel) error
	SaveDeployeIngestTask(tasks []entity.DeployIngestModel) (map[string]int, error)
	UpdateDeployeIngestTask(id []int, status int) error
	ReleaseRegion(regionId int) error
}

type mysqlDB struct {
	DB *gorm.DB
}

func NewMysql(conf common.Mysql) (*mysqlDB, error) {
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
		return nil, fmt.Errorf("数据源配置不正确: %v", err.Error())
	}

	db, err := gormDB.DB()
	if err != nil {
		return nil, fmt.Errorf("gorm 获取数据库失败: %v", err.Error())
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("数据库连接失败: %v", err.Error())
	}

	// 最大连接数
	db.SetMaxOpenConns(100)
	// 闲置连接数
	db.SetMaxIdleConns(20)
	// 最大连接周期
	db.SetConnMaxLifetime(100 * time.Second)

	return &mysqlDB{DB: gormDB}, nil
}

func (cli *mysqlDB) GetUser(id int) (*entity.User, error) {
	var result entity.User

	err := cli.DB.Table(entity.UserTableName).Where("id = ?", id).Find(&result).Error

	return &result, err
}

func (cli *mysqlDB) GetUserConfigName(ingestID, version string) (string, error) {
	var configName string
	err := cli.DB.Table(entity.UserTableName).Select("config_name").
		Where("user_id = ? and version = ?", ingestID, version).
		Find(&configName).Error
	return configName, err
}

func (cli *mysqlDB) ExitsNotifyByUUId(uuid string) (string, error) {
	tmp1 := &entity.NotifyMsgModel{}
	err := cli.DB.Table(entity.NotifyMsgTableName).Where("uuid = ?", uuid).First(&tmp1).Error
	if err == gorm.ErrRecordNotFound {
		return entity.NotifyMsgTableName, nil
	}
	if err != nil {
		common.Logger.Errorf("infra search NotifyMsgTableName error: %+v", err)
		return entity.NotifyMsgTableName, err
	}

	tmp2 := &entity.DeployIngestModel{}
	err = cli.DB.Table(entity.DeployIngestTableName).Where("notify_id = ?", uuid).First(&tmp2).Error
	if err == gorm.ErrRecordNotFound {
		return entity.DeployIngestTableName, nil
	}
	if err != nil {
		common.Logger.Errorf("infra search DeployIngestTableName error: %+v", err)
		return entity.DeployIngestTableName, err
	}
	return "exist", nil
}

func (cli *mysqlDB) SaveNotifyMessage(msg *entity.NotifyMsgModel) error {
	return cli.DB.Table(entity.NotifyMsgTableName).Create(msg).Error
}

func (cli *mysqlDB) SaveDeployeIngestTask(tasks []entity.DeployIngestModel) (map[string]int, error) {
	err := cli.DB.Table(entity.DeployIngestTableName).Create(tasks).Error
	if err != nil {
		common.Logger.Errorf("infra error: DeployIngestTableName create:%s", err)
		return nil, err
	}

	ids := make(map[string]int, 1)
	for _, task := range tasks {
		ids[task.GameIp] = task.Id
	}

	return ids, nil
}

func (cli *mysqlDB) UpdateDeployeIngestTask(id []int, status int) error {
	return cli.DB.Table(entity.DeployIngestTableName).Where("id in ?", id).Update("status", status).Error
}

func (cli *mysqlDB) ReleaseRegion(regionId int) error {
	return cli.DB.Table(entity.DeployIngestTableName).Where("region_id = ?", regionId).Update("status", 0).Error
}

func (cli *mysqlDB) Close() {
	db, _ := cli.DB.DB()
	err := db.Close()
	if err != nil {
		common.Logger.Errorf("mysql close error: %v", err)
	}
}
