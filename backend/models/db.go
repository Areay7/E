package models

import (
	"fmt"
	"log"
	"time"

	"cross-border-admin/config"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB() error {
	var dialector gorm.Dialector
	dbConfig := config.AppConfig.Database

	switch dbConfig.Type {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
			dbConfig.Username,
			dbConfig.Password,
			dbConfig.Host,
			dbConfig.Port,
			dbConfig.DBName,
			dbConfig.Charset,
		)
		dialector = mysql.Open(dsn)

	case "postgres":
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
			dbConfig.Host,
			dbConfig.Port,
			dbConfig.Username,
			dbConfig.Password,
			dbConfig.DBName,
			dbConfig.SSLMode,
			dbConfig.Timezone,
		)
		dialector = postgres.Open(dsn)

	default:
		return fmt.Errorf("不支持的数据库类型: %s", dbConfig.Type)
	}

	// 配置日志级别
	var logLevel logger.LogLevel
	switch dbConfig.LogLevel {
	case "silent":
		logLevel = logger.Silent
	case "error":
		logLevel = logger.Error
	case "warn":
		logLevel = logger.Warn
	case "info":
		logLevel = logger.Info
	default:
		logLevel = logger.Info
	}

	// 连接数据库
	var err error
	DB, err = gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return fmt.Errorf("数据库连接失败: %w", err)
	}

	// 获取底层 sql.DB
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("获取数据库实例失败: %w", err)
	}

	// 设置连接池
	sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)
	sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 自动迁移
	if dbConfig.AutoMigrate {
		if err := autoMigrate(); err != nil {
			return fmt.Errorf("数据库迁移失败: %w", err)
		}
		log.Println("数据库自动迁移完成")
	} else {
		log.Println("数据库自动迁移已关闭，跳过 AutoMigrate")
	}

	log.Println("数据库初始化成功")
	return nil
}

// autoMigrate 自动迁移所有表
func autoMigrate() error {
	return DB.AutoMigrate(
		&User{},
		&Order{},
		&OrderItem{},
		&Product{},
		&ProductVariation{},
		&Inventory{},
		&InventoryLog{},
		&Logistics{},
		&PlatformConfig{},
		&APILog{},
		&SyncTask{},
	)
}
