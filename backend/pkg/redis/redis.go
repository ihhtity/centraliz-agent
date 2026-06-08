package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"centraliz-backend/pkg/config"
)

var RDB *redis.Client

func InitRedis() error {
	redisConfig := config.AppConfig.Redis

	// 创建Redis客户端配置
	options := &redis.Options{
		Addr:     config.GetRedisAddr(),
		DB:       redisConfig.DB,
		PoolSize: redisConfig.PoolSize,
	}

	// 只有当密码不为空时才设置密码
	// 如果设置空密码，go-redis会尝试发送AUTH命令，导致"ERR Client sent AUTH, but no password is set"错误
	if redisConfig.Password != "" {
		options.Password = redisConfig.Password
	}

	RDB = redis.NewClient(options)

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return RDB.Ping(ctx).Err()
}

// Set 设置键值对
func Set(key string, value interface{}, expiration time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return RDB.Set(ctx, key, value, expiration).Err()
}

// Get 获取键值
func Get(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return RDB.Get(ctx, key).Result()
}

// Del 删除键
func Del(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return RDB.Del(ctx, key).Err()
}

// Exists 检查键是否存在
func Exists(key string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := RDB.Exists(ctx, key).Result()
	return result > 0, err
}