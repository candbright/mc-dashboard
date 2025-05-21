package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

// Config 应用配置
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	SMS      SMSConfig      `mapstructure:"sms"`
	MC       MCConfig       `mapstructure:"mc"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port    int           `mapstructure:"port"`
	Mode    string        `mapstructure:"mode"`
	Timeout time.Duration `mapstructure:"timeout"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver       string `mapstructure:"driver"`
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Username     string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
	DBName       string `mapstructure:"dbname"`
	Charset      string `mapstructure:"charset"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret      string `mapstructure:"secret"`
	ExpireHours int    `mapstructure:"expire_hours"`
}

// SMSConfig 短信配置
type SMSConfig struct {
	Provider     string `mapstructure:"provider"`
	AccessKey    string `mapstructure:"access_key"`
	AccessSecret string `mapstructure:"access_secret"`
	SignName     string `mapstructure:"sign_name"`
	TemplateCode string `mapstructure:"template_code"`
}

type MCConfig struct {
	Path string `mapstructure:"path"`
}

var (
	// GlobalConfig 全局配置实例
	GlobalConfig Config
)

// Init 初始化配置
func Init() error {
	// 设置配置文件路径
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("configs")

	// 读取环境变量
	viper.AutomaticEnv()
	viper.SetEnvPrefix("MC_DASHBOARD")

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 解析配置到结构体
	if err := viper.Unmarshal(&GlobalConfig); err != nil {
		return fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 从环境变量覆盖配置
	overrideFromEnv()

	return nil
}

// overrideFromEnv 从环境变量覆盖配置
func overrideFromEnv() {
	// 数据库配置
	if host := os.Getenv("MC_DASHBOARD_DB_HOST"); host != "" {
		GlobalConfig.Database.Host = host
	}
	if port := os.Getenv("MC_DASHBOARD_DB_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			GlobalConfig.Database.Port = p
		}
	}
	if user := os.Getenv("MC_DASHBOARD_DB_USER"); user != "" {
		GlobalConfig.Database.Username = user
	}
	if pass := os.Getenv("MC_DASHBOARD_DB_PASSWORD"); pass != "" {
		GlobalConfig.Database.Password = pass
	}
	if dbname := os.Getenv("MC_DASHBOARD_DB_NAME"); dbname != "" {
		GlobalConfig.Database.DBName = dbname
	}

	// JWT配置
	if secret := os.Getenv("MC_DASHBOARD_JWT_SECRET"); secret != "" {
		GlobalConfig.JWT.Secret = secret
	}

	// 服务器配置
	if port := os.Getenv("MC_DASHBOARD_SERVER_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			GlobalConfig.Server.Port = p
		}
	}

	// MC配置
	if path := os.Getenv("MC_DASHBOARD_MC_PATH"); path != "" {
		GlobalConfig.MC.Path = path
	}
}

// GetDSN 获取数据库连接字符串
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.DBName,
		c.Charset,
	)
}
