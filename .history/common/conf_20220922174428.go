package common

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	red "github.com/go-redis/redis/v8"
	"github.com/opensearch-project/opensearch-go"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	configPath = "./conf"
	configName = "logservice2"
	configExt  = "toml"

	appConf     *AppConfig
	settingOnce sync.Once
)

type AppConfig struct {
	Logging       Logging       `mapstructure:"Log"`
	Server        Web           `mapstructure:"Web"`
	Redis         Redis         `mapstructure:"Kafka"`
	Mysql         Mysql         `mapstructure:"Mysql"`
	Opensearch    Opensearch    `mapstructure:"Opensearch"`
	Elasticsearch Elasticsearch `mapstructure:"Elasticsearch"`

	Tunnel Tunnel `mapstructure:"tunnel"`

	DB        *gorm.DB           `json:"-"`
	RedisCli  *red.Client        `json:"-"`
	OpenDB    *opensearch.Client `json:"-"`
	TunnelCli http.Client        `json:"-"`
	El    *opensearch.Client `json:"-"`
}

type Web struct {
	RunMode                  string        `mapstructure:"run_mode"`
	HTTPPort                 int           `mapstructure:"http_port"`
	ServiceName              string        `mapstructure:"service_name"`
	Language                 string        `mapstructure:"language"`
	ReadTimeOut              time.Duration `mapstructure:"read_timeOut"`
	WriteTimeOut             time.Duration `mapstructure:"write_timeOut"`
	ConnectCheckTimeout      time.Duration `mapstructure:"connect_check_timeout"`
	RepositoryRequestTimeout time.Duration `mapstructure:"repository_request_timeout"`
}

type Logging struct {
	LogFilePath string `mapstructure:"log_file_path"`
	LogLevel    string `mapstructure:"log_level"`
	DevelopMode bool   `mapstructure:"develop_mode"`
	MaxAge      int    `mapstructure:"max_age"`
	MaxBackups  int    `mapstructure:"max_backups"`
	MaxSize     int    `mapstructure:"max_size"`

	atomicLevel zap.AtomicLevel
}

type Redis struct {
	DB           int    `mapstructure:"db"`
	Addr         string `mapstructure:"addr"`
	Password     string `mapstructure:"password"`
	MaxRetries   int    `mapstructure:"max_retries"`
	MinIdleConns int    `mapstructure:"min_idle_conns"`
}

type Mysql struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	DataBase string `mapstructure:"database"`
	Username string
	Password string `json:"-"`
}

type Opensearch struct {
	InsecureSkipVerify bool     `mapstructure:"insecure_skip_verify"`
	Address            []string `mapstructure:"address"`
	Username           string   `mapstructure:"username"`
	Password           string   `mapstructure:"password"`
}

type Elasticsearch struct {
	Address  []string `mapstructure:"address"`
	Username string   `mapstructure:"username"`
	Password string   `mapstructure:"password"`
}

type Tunnel struct {
	Timeout            int  `mapstructure:"timeout"`
	IdleConnTimeout    int  `mapstructure:"idle_conn_timeout"`
	DisableKeepAlives  bool `mapstructure:"disable_keep_alives"`
	InsecureSkipVerify bool `mapstructure:"insecure_skip_verify"`
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
	Logger.Infof("Init Setting From File %s/%s.%s", configPath, configName, configExt)

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
	Logger.Infof("Load Setting File %s/%s.%s", configPath, configName, configExt)

	if err := vp.ReadInConfig(); err != nil {
		Logger.Fatalf("err:%s\n", err)
	}

	if err := vp.Unmarshal(appConf); err != nil {
		Logger.Fatalf("err:%s\n", err)
	}

	appConf.Logging.SetLogLevel()

	s, _ := json.Marshal(appConf)

	Logger.Info(string(s))
}

// SetMySQLSetting 初始化mysql：账号密码解密、格式校验等
func (appConf *AppConfig) SetMySQLSetting() {
	//appConf.Mysql.Host = "10.4.106.129"
	//appConf.Mysql.Port = 30006
	//appConf.Mysql.Username = "root"
	//appConf.Mysql.Password = "root"
}
