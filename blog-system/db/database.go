package db

import (
	"blog-system/config"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

var DB *Database

type Database struct {
	Self   *gorm.DB
	Docker *gorm.DB
}

// 初始化数据库连接
func InitDB() {
	// 主数据库
	selfDB, err := createDBConnection("db")
	if err != nil {
		config.Logger.Fatalf("Failed to connect to self database: %v", err)
	}

	// Docker数据库（如果需要）
	var dockerDB *gorm.DB
	if viper.GetBool("docker_db.enabled") {
		dockerDB, err = createDBConnection("docker_db")
		if err != nil {
			config.Logger.Errorf("Failed to connect to docker database: %v", err)
		}
	}

	DB = &Database{
		Self:   selfDB,
		Docker: dockerDB,
	}
}

// 创建数据库连接
func createDBConnection(prefix string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		viper.GetString(prefix+".username"),
		viper.GetString(prefix+".password"),
		viper.GetString(prefix+".addr"),
		viper.GetString(prefix+".name"),
	)

	// 配置GORM日志
	gormLogger := logger.Default
	if viper.GetBool("gormlog") {
		gormLogger = logger.Default.LogMode(logger.Info)
	} else {
		gormLogger = logger.Default.LogMode(logger.Silent)
	}

	// 创建数据库连接
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// 获取通用数据库对象并设置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// 配置连接池
	sqlDB.SetMaxIdleConns(viper.GetInt(prefix + ".max_idle_conns"))
	sqlDB.SetMaxOpenConns(viper.GetInt(prefix + ".max_open_conns"))
	sqlDB.SetConnMaxLifetime(time.Duration(viper.GetInt(prefix+".conn_max_lifetime")) * time.Minute)

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	config.Logger.Infof("Connected to database: %s", viper.GetString(prefix+".name"))
	return db, nil
}
