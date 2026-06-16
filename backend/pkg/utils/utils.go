package utils

import (
	"centraliz-backend/pkg/redis"
	"context"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

// StringPtr 返回字符串指针
func StringPtr(s string) *string {
	return &s
}

// Int32Ptr 返回int32指针
func Int32Ptr(i int) *int32 {
	val := int32(i)
	return &val
}

// Int64Ptr 返回int64指针
func Int64Ptr(i int) *int64 {
	val := int64(i)
	return &val
}

// Float64Ptr 返回float64指针
func Float64Ptr(f float64) *float64 {
	return &f
}

// IsValidPassword 验证密码强度
func IsValidPassword(password string) bool {
	if len(password) < 6 || len(password) > 20 {
		return false
	}
	matched, _ := regexp.MatchString(`\p{Han}`, password)
	return !matched
}

// GenerateRandomPassword 生成随机密码
func GenerateRandomPassword(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// IsValidPhone 验证手机号格式
func IsValidPhone(phone string) bool {
	phone = strings.TrimSpace(phone)
	pattern := `^1[3-9]\d{9}$`
	matched, _ := regexp.MatchString(pattern, phone)
	return matched
}

// IsValidEmail 验证邮箱格式
func IsValidEmail(email string) bool {
	email = strings.TrimSpace(email)
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

// VerifyCode 验证验证码
func VerifyCode(key, code string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	storedCode, err := redis.RDB.Get(ctx, key).Result()
	if err != nil {
		return false
	}

	return storedCode == code
}

// DeleteUsedCode 删除已使用的验证码
func DeleteUsedCode(key string) {
	if key == "account" {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	redis.RDB.Del(ctx, key)
}
