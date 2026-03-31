package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server    ServerConfig              `mapstructure:"server"`
	Database  DatabaseConfig            `mapstructure:"database"`
	JWT       JWTConfig                 `mapstructure:"jwt"`
	CORS      CORSConfig                `mapstructure:"cors"`
	Log       LogConfig                 `mapstructure:"log"`
	Platforms map[string]PlatformConfig `mapstructure:"platforms"`
	Sync      SyncConfig                `mapstructure:"sync"`
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Type         string `mapstructure:"type"`
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Username     string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
	DBName       string `mapstructure:"dbname"`
	Charset      string `mapstructure:"charset"`
	SSLMode      string `mapstructure:"sslmode"`
	Timezone     string `mapstructure:"timezone"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	LogLevel     string `mapstructure:"log_level"`
	AutoMigrate  bool   `mapstructure:"auto_migrate"`
}

type JWTConfig struct {
	Secret      string `mapstructure:"secret"`
	ExpireHours int    `mapstructure:"expire_hours"`
}

type CORSConfig struct {
	AllowOrigins     []string `mapstructure:"allow_origins"`
	AllowCredentials bool     `mapstructure:"allow_credentials"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	FilePath   string `mapstructure:"file_path"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}

type PlatformConfig struct {
	Enabled    bool   `mapstructure:"enabled"`
	PartnerID  string `mapstructure:"partner_id"`
	PartnerKey string `mapstructure:"partner_key"`
	ShopID     string `mapstructure:"shop_id"`
	AppKey     string `mapstructure:"app_key"`
	AppSecret  string `mapstructure:"app_secret"`
	APIURL     string `mapstructure:"api_url"`
}

type SyncConfig struct {
	OrderInterval     int `mapstructure:"order_interval"`
	ProductInterval   int `mapstructure:"product_interval"`
	InventoryInterval int `mapstructure:"inventory_interval"`
}

var AppConfig *Config

// LoadConfig 加载配置文件
func LoadConfig(configPath string) error {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 初始化 AppConfig
	AppConfig = &Config{}

	// 解析配置
	if err := viper.Unmarshal(AppConfig); err != nil {
		return fmt.Errorf("解析配置文件失败: %w", err)
	}

	log.Printf("配置文件加载成功，服务器端口: %d", AppConfig.Server.Port)
	return nil
}

// GetPlatformConfig 获取指定平台的配置
func GetPlatformConfig(platform string) (PlatformConfig, bool) {
	config, exists := AppConfig.Platforms[platform]
	return config, exists
}
