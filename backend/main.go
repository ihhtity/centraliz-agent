package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"centraliz-backend/middleware"
	"centraliz-backend/pkg/config"
	"centraliz-backend/pkg/db"
	"centraliz-backend/pkg/log"
	"centraliz-backend/pkg/mail"
	"centraliz-backend/pkg/redis"
	"centraliz-backend/router"
)

func main() {
	// 初始化配置
	config.InitConfig()

	// 初始化日志
	log.InitLogger()
	defer log.Logger.Sync()

	// 初始化数据库
	db.InitDB()

	// 初始化Redis
	redis.InitRedis()

	// 初始化邮件服务
	if err := mail.InitMailService(); err != nil {
		log.Logger.Fatal("邮件服务初始化失败", zap.Error(err))
	}

	// 创建gin实例
	if !config.IsDebugMode() {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()

	// 注册中间件
	r.Use(gin.Recovery())
	r.Use(middleware.Logger(log.Logger))
	r.Use(middleware.CORS())

	// 仅在生产环境启用限流
	if !config.IsDebugMode() && config.AppConfig.RateLimit.RequestsPerSecond > 0 {
		r.Use(middleware.RateLimit())
	}

	// 注册路由
	router.InitRouter(r)

	// 创建HTTP服务器
	server := &http.Server{
		Addr:    config.GetServerAddress(),
		Handler: r,
	}

	// 启动服务器
	go func() {
		log.Logger.Info("服务器启动中", zap.String("addr", server.Addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Logger.Fatal("服务器启动失败", zap.Error(err))
		}
	}()

	// 等待中断信号优雅关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Logger.Info("正在关闭服务器...")

	// 给予5秒超时时间优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Logger.Fatal("服务器强制关闭", zap.Error(err))
	}

	log.Logger.Info("服务器已退出")
}