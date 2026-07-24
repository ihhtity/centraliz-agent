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

	options := &redis.Options{
		Addr:     config.GetRedisAddr(),
		DB:       redisConfig.DB,
		PoolSize: redisConfig.PoolSize,
		MinIdleConns: 5,
		MaxRetries: 3,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolTimeout:  4 * time.Second,
	}

	if redisConfig.Password != "" {
		options.Password = redisConfig.Password
	}

	RDB = redis.NewClient(options)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return RDB.Ping(ctx).Err()
}

func Set(key string, value interface{}, expiration time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return RDB.Set(ctx, key, value, expiration).Err()
}

func Get(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return RDB.Get(ctx, key).Result()
}

func Del(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return RDB.Del(ctx, key).Err()
}

func Exists(key string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := RDB.Exists(ctx, key).Result()
	return result > 0, err
}

func HSet(key string, values ...interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return RDB.HSet(ctx, key, values...).Err()
}

func HGet(key, field string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return RDB.HGet(ctx, key, field).Result()
}

func HGetAll(key string) (map[string]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return RDB.HGetAll(ctx, key).Result()
}

func HDel(key string, fields ...string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return RDB.HDel(ctx, key, fields...).Err()
}

func Incr(key string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return RDB.Incr(ctx, key).Result()
}

func IncrBy(key string, increment int64) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return RDB.IncrBy(ctx, key, increment).Result()
}

func Decr(key string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return RDB.Decr(ctx, key).Result()
}

func DecrBy(key string, decrement int64) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return RDB.DecrBy(ctx, key, decrement).Result()
}

func LPush(key string, values ...interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return RDB.LPush(ctx, key, values...).Err()
}

func RPush(key string, values ...interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return RDB.RPush(ctx, key, values...).Err()
}

func LPop(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return RDB.LPop(ctx, key).Result()
}

func RPop(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return RDB.RPop(ctx, key).Result()
}

func LLen(key string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return RDB.LLen(ctx, key).Result()
}

func LRange(key string, start, stop int64) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return RDB.LRange(ctx, key, start, stop).Result()
}

func Expire(key string, expiration time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return RDB.Expire(ctx, key, expiration).Err()
}

func TTL(key string) (time.Duration, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return RDB.TTL(ctx, key).Result()
}

func Keys(pattern string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return RDB.Keys(ctx, pattern).Result()
}

func FlushDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return RDB.FlushDB(ctx).Err()
}

func GetPoolStats() *redis.PoolStats {
	if RDB == nil {
		return nil
	}
	return RDB.PoolStats()
}

func IsConnected() bool {
	if RDB == nil {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := RDB.Ping(ctx).Result()
	return err == nil
}