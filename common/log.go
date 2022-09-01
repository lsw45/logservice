package common

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

type Log = zap.SugaredLogger

var (
	defaultLogPath = ""
	Logger         *Log
)

// SetLogLevel 设置日志级别
func (logConf *Logging) SetLogLevel() {
	atomicLevel, err := zap.ParseAtomicLevel(logConf.LogLevel)

	if err != nil {
		Logger.Errorf("Parse Log Level failed:%s", logConf.LogLevel)
		return
	}
	logConf.atomicLevel = atomicLevel

	Logger = logConf.initLogger()
	Logger.Infof(logConf.LogLevel, logConf.atomicLevel)
}

// 初始化日志对象
func (logConf *Logging) initLogger() *Log {
	// 日志文件hook
	hook := &lumberjack.Logger{
		Filename:   logConf.LogFilePath,
		LocalTime:  true,               //日志文件名的时间格式为本地时间
		MaxAge:     logConf.MaxAge,     //文件保留的最长时间，单位为天
		MaxBackups: logConf.MaxBackups, // 旧文件保留的最大个数
		MaxSize:    logConf.MaxSize,    // 单个文件最大长度，单位是M
	}

	// 日志格式设定
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "linenum",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),                                          // 编码器配置
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(hook)), // 打印到控制台和文件
		logConf.atomicLevel, // 日志级别
	)

	// 设置初始化字段
	filed := zap.Fields(zap.String("serviceName", "ingest_manager"))
	// 构造日志
	var logger *zap.SugaredLogger
	if logConf.DevelopMode {
		// 开启开发模式，堆栈跟踪
		caller := zap.AddCaller()
		// 开启文件及行号
		development := zap.Development()
		logger = zap.New(core, caller, development, filed).Sugar()
	} else {
		logger = zap.New(core, filed).Sugar()
	}

	return logger
}

func init() {
	atomicLevel := zap.NewAtomicLevelAt(zap.InfoLevel)
	lf := Logging{
		LogFilePath: defaultLogPath,
		DevelopMode: true,
		MaxAge:      100,
		MaxBackups:  20,
		MaxSize:     100,
		atomicLevel: atomicLevel,
	}
	Logger = lf.initLogger()
}
