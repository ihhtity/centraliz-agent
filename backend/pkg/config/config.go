package config

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// ServerConfig 服务器配置
type ServerConfig struct {
	Port  string `mapstructure:"port"`
	Host  string `mapstructure:"host"`
	Debug bool   `mapstructure:"debug"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host            string `mapstructure:"host"`
	Port            string `mapstructure:"port"`
	Name            string `mapstructure:"name"`
	Username        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret      string `mapstructure:"secret"`
	ExpireHours int    `mapstructure:"expire_hours"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
}

// CORSConfig CORS配置
type CORSConfig struct {
	AllowedOrigins []string `mapstructure:"allowed_origins"`
	AllowedMethods []string `mapstructure:"allowed_methods"`
	AllowedHeaders []string `mapstructure:"allowed_headers"`
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	RequestsPerSecond int `mapstructure:"requests_per_second"`
	Burst             int `mapstructure:"burst"`
}

// MailConfig 邮件配置
type MailConfig struct {
	Type       string `mapstructure:"mail_type"`
	Host       string `mapstructure:"mail_smtp_host"`
	Port       int    `mapstructure:"mail_smtp_port"`
	User       string `mapstructure:"mail_smtp_user"`
	Password   string `mapstructure:"mail_smtp_pass"`
	VerifyType string `mapstructure:"mail_verify_type"`
	From       string `mapstructure:"mail_from"`
}

// Config 主配置结构体
type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	Database  DatabaseConfig  `mapstructure:"database"`
	Redis     RedisConfig     `mapstructure:"redis"`
	JWT       JWTConfig       `mapstructure:"jwt"`
	Log       LogConfig       `mapstructure:"log"`
	CORS      CORSConfig      `mapstructure:"cors"`
	RateLimit RateLimitConfig `mapstructure:"rate_limit"`
	Mail      *MailConfig     `mapstructure:"mail"`
}

var (
	AppConfig *Config
	configMux sync.RWMutex
)

// InitConfig 初始化配置
func InitConfig() {
	// 设置默认配置（必须先调用，否则配置文件缺失时会出错）
	setDefaultConfig()

	// 配置Viper使用单一配置文件
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config/")
	viper.AddConfigPath("/etc/centraliz/")
	viper.AddConfigPath(".")

	// 读取环境变量
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 尝试读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("未找到配置文件 config.yaml，使用默认配置")
		} else {
			log.Fatalf("读取配置文件错误: %v", err)
		}
	} else {
		log.Printf("使用配置文件: %s", viper.ConfigFileUsed())
	}

	// 解析配置到结构体
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("无法将配置解析到结构体: %v", err)
	}

	configMux.Lock()
	AppConfig = &config
	configMux.Unlock()

	// 验证必要配置
	validateConfig()

	// 启动配置热更新监听（仅在调试模式下）
	if config.Server.Debug {
		startConfigWatcher()
	}
}

// setDefaultConfig 设置默认配置
func setDefaultConfig() {
	viper.SetDefault("server.port", ":8080")
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.debug", true)

	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", "3306")
	viper.SetDefault("database.name", "centraliz")
	viper.SetDefault("database.username", "root")
	viper.SetDefault("database.password", "123456")
	viper.SetDefault("database.max_open_conns", 25)
	viper.SetDefault("database.max_idle_conns", 5)
	viper.SetDefault("database.conn_max_lifetime", 3600)

	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", "6379")
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("redis.pool_size", 10)

	viper.SetDefault("jwt.secret", "dev-secret-key-change-in-production")
	viper.SetDefault("jwt.expire_hours", 24)

	viper.SetDefault("log.level", "debug")
	viper.SetDefault("log.format", "console")
	viper.SetDefault("log.output", "stdout")

	viper.SetDefault("cors.allowed_origins", []string{
		"http://localhost:8080",
		"http://localhost:5173",
		"http://127.0.0.1:5173",
	})
	viper.SetDefault("cors.allowed_methods", []string{
		"GET", "POST", "PUT", "DELETE", "OPTIONS",
	})
	viper.SetDefault("cors.allowed_headers", []string{
		"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With",
	})

	viper.SetDefault("rate_limit.requests_per_second", 100)
	viper.SetDefault("rate_limit.burst", 200)
}

// validateConfig 验证必要配置项
func validateConfig() {
	configMux.RLock()
	defer configMux.RUnlock()

	if AppConfig.Database.Name == "" {
		log.Fatal("数据库名称是必需的")
	}
	if AppConfig.Database.Username == "" {
		log.Fatal("数据库用户名是必需的")
	}
	if AppConfig.JWT.Secret == "" {
		log.Fatal("JWT密钥是必需的")
	}
}

// startConfigWatcher 启动配置文件监听器
func startConfigWatcher() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("配置文件已更改:", e.Name)
		var newConfig Config
		if err := viper.Unmarshal(&newConfig); err != nil {
			log.Printf("解析新配置错误: %v", err)
			return
		}

		configMux.Lock()
		AppConfig = &newConfig
		configMux.Unlock()

		log.Println("配置重新加载成功")
	})
}

// GetMySQLDSN 获取MySQL连接字符串
func GetMySQLDSN() string {
	configMux.RLock()
	defer configMux.RUnlock()

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		AppConfig.Database.Username,
		AppConfig.Database.Password,
		AppConfig.Database.Host,
		AppConfig.Database.Port,
		AppConfig.Database.Name)
}

// GetRedisAddr 获取Redis地址
func GetRedisAddr() string {
	configMux.RLock()
	defer configMux.RUnlock()

	return fmt.Sprintf("%s:%s", AppConfig.Redis.Host, AppConfig.Redis.Port)
}

// GetServerAddress 获取服务器地址
func GetServerAddress() string {
	configMux.RLock()
	defer configMux.RUnlock()

	if AppConfig.Server.Host == "" {
		return AppConfig.Server.Port
	}
	return fmt.Sprintf("%s%s", AppConfig.Server.Host, AppConfig.Server.Port)
}

// IsDebugMode 是否为调试模式
func IsDebugMode() bool {
	configMux.RLock()
	defer configMux.RUnlock()

	return AppConfig.Server.Debug
}
