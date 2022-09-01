package common

import (
	"encoding/json"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"sync"
	"time"
)

var (
	configPath = "./server/conf"
	configName = "app"
	configExt  = "toml"

	appConf     *AppConfig
	settingOnce sync.Once
)

type AppConfig struct {
	Logging    Logging    `toml:"Log"`
	Server     Web        `toml:"Web"`
	Kafka      Kafka      `toml:"Kafka"`
	Mysql      Mysql      `toml:"Mysql"`
	Opensearch Opensearch `toml:"Opensearch"`

	DB *gorm.DB `json:"-"`
}
type Web struct {
	RunMode                  string        `toml:"run_mode"`
	HTTPPort                 int           `toml:"http_port"`
	ServiceName              string        `toml:"service_name"`
	Language                 string        `toml:"language"`
	ReadTimeOut              time.Duration `toml:"read_timeOut"`
	WriteTimeOut             time.Duration `toml:"write_timeOut"`
	ConnectCheckTimeout      time.Duration `toml:"connect_check_timeout"`
	RepositoryRequestTimeout time.Duration `toml:"repository_request_timeout"`
}

type Logging struct {
	LogFilePath string `toml:"logFilePath"`
	LogLevel    string `toml:"log_level"`
	DevelopMode bool   `toml:"develop_mode"`
	MaxAge      int    `toml:"max_age"`
	MaxBackups  int    `toml:"max_backups"`
	MaxSize     int    `toml:"max_size"`

	atomicLevel zap.AtomicLevel
}
type Kafka struct {
}

type Mysql struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Username string
	Password string `json:"-"`
}
type Opensearch struct {
}

// NewAppConfig 读取服务配置
func NewAppConfig() *AppConfig {
	settingOnce.Do(func() {
		appConf = &AppConfig{}
		vp := viper.New()
		initSetting(vp)
	})

	return appConf
}

// 初始化配置
func initSetting(vp *viper.Viper) {
	Logger.Infof("Init Setting From File %s%s.%s", configPath, configName, configExt)

	vp.AddConfigPath(configPath)
	vp.SetConfigName(configName)
	vp.SetConfigType(configExt)

	loadSetting(vp)

	vp.WatchConfig()
	vp.OnConfigChange(func(e fsnotify.Event) {
		Logger.Infof("Config file changed:%s", e)
		loadSetting(vp)
	})
}

// 读取配置文件
func loadSetting(vp *viper.Viper) {
	Logger.Infof("Load Setting File %s%s.%s", configPath, configName, configExt)

	if err := vp.ReadInConfig(); err != nil {
		Logger.Fatalf("err:%s\n", err)
	}

	if err := vp.Unmarshal(appConf); err != nil {
		Logger.Fatalf("err:%s\n", err)
	}

	appConf.SetMySQLSetting()
	appConf.Logging.SetLogLevel()

	s, _ := json.Marshal(appConf)

	Logger.Info(string(s))
}

// SetMySQLSetting 初始化mysql
func (appConf *AppConfig) SetMySQLSetting() {
	appConf.Mysql.Host = "10.4.106.129"
	appConf.Mysql.Port = 30006
	appConf.Mysql.Username = "root"
	appConf.Mysql.Password = "root"
}
