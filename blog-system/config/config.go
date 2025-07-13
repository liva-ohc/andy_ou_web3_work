package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"strings"
)

type Config struct {
	Name string
}

// 全局日志记录器
var Logger *zap.SugaredLogger

func Init(cfg string) error {
	c := Config{
		Name: cfg,
	}

	// 初始化配置文件
	if err := c.initConfig(); err != nil {
		return err
	}

	// 初始化日志包
	if err := c.initLog(); err != nil {
		return err
	}

	// 监控配置文件变化并热加载程序
	c.watchConfig()

	return nil
}

func (c *Config) initConfig() error {
	if c.Name != "" {
		// 指定配置文件，捷信指定的文件
		viper.SetConfigFile(c.Name)
	} else {
		// 默认文件
		viper.AddConfigPath("conf")   // 目录
		viper.SetConfigName("config") // 文件名
	}
	// 设置配置文件格式为yaml格式
	viper.SetConfigType("yaml")
	// 自动匹配环境变量
	viper.AutomaticEnv()
	// 读取环境变量的前缀为APISERVER
	viper.SetEnvPrefix("APISERVER")
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	if err := viper.ReadInConfig(); err != nil {
		// viper解析文件错误
		return err
	}
	return nil
}

func (c *Config) initLog() error {
	// 从配置中获取日志设置
	logLevel := viper.GetString("log.level")
	logFile := viper.GetString("log.file")

	// 设置日志级别
	var level zapcore.Level
	switch strings.ToLower(logLevel) {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}

	// 创建核心日志器
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), // JSON 格式
		zapcore.NewMultiWriteSyncer( // 多输出目标
			zapcore.AddSync(os.Stdout), // 控制台输出
			zapcore.AddSync(&lumberjack.Logger{ // 文件输出带日志轮转
				Filename:   logFile,
				MaxSize:    100, // 最大100MB
				MaxBackups: 3,   // 保留3个备份
				MaxAge:     30,  // 保留30天
			}),
		),
		level,
	)

	// 创建日志器
	logger := zap.New(core, zap.AddCaller())
	defer logger.Sync() // 确保程序退出前刷新日志

	// 创建带Sugar的日志器
	Logger = logger.Sugar()

	// 记录日志初始化成功
	Logger.Infow("日志初始化成功",
		"level", level.String(),
		"file", logFile,
	)

	return nil
}

func (c *Config) watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		Logger.Infof("配置文件已更新: %s", in.Name)
		// 重新加载日志配置
		if err := c.initLog(); err != nil {
			Logger.Errorf("重新加载日志配置失败: %v", err)
		}
	})
}
