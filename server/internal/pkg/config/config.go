package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config 应用配置结构
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Mail     MailConfig     `mapstructure:"mail"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Log      LogConfig      `mapstructure:"log"`
	OSS      OSS            `mapstructure:"oss"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	Mode         string        `mapstructure:"mode"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver          string `mapstructure:"driver"`
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	Database        string `mapstructure:"database"`
	Username        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	Charset         string `mapstructure:"charset"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
	LogLevel        string `mapstructure:"log_level"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Database int    `mapstructure:"database"`
	PoolSize int    `mapstructure:"pool_size"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret      string        `mapstructure:"secret"`
	AccessTime  time.Duration `mapstructure:"access_time"`
	RefreshTime time.Duration `mapstructure:"refresh_time"`
	Issuer      string        `mapstructure:"issuer"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"`
	Output     string `mapstructure:"output"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
	Compress   bool   `mapstructure:"compress"`
}

// MailConfig 邮件配置
type MailConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	FromName string `mapstructure:"from_name"`
	FromMail string `mapstructure:"from_mail"`
	UserName string `mapstructure:"username"`
	AuthCode string `mapstructure:"auth_code"`
}

type OSS struct {
	Provider   string     `mapstructure:"provider"`
	MinioConf  MinioConf  `mapstructure:"minio_conf"`
	AliyunConf AliyunConf `mapstructure:"aliyun_conf"`
}

type MinioConf struct {
	Endpoint        string `mapstructure:"endpoint"`
	AccessURL       string `mapstructure:"access_url"`
	AccessKeyId     string `mapstructure:"access_key_id"`
	SecretAccessKey string `mapstructure:"secret_access_key"`
	BucketName      string `mapstructure:"bucket_name"`
	UseSSL          bool   `mapstructure:"use_ssl"`
}

type AliyunConf struct {
	AccessURL       string `mapstructure:"access_url"`
	AccessKeyID     string `mapstructure:"access_key_id"`
	AccessKeySecret string `mapstructure:"access_key_secret"`
	Region          string `mapstructure:"region"`
	BucketName      string `mapstructure:"bucket_name"`
}

// Load 加载配置
func Load(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	return &config, nil
}
