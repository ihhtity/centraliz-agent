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

// 全局变量用于追踪服务器状态
var (
	serverStarted = false
)

func main() {
	// 初始化配置
	config.InitConfig()

	// 初始化日志
	log.InitLogger()
	defer log.Logger.Sync()

	// 初始化数据库
	log.Logger.Info("正在初始化数据库...")
	db.InitDB()
	log.Logger.Info("数据库初始化成功")

	// 初始化Redis
	log.Logger.Info("正在初始化Redis...")
	redis.InitRedis()
	log.Logger.Info("Redis初始化成功")

	// 初始化邮件服务（非必需，失败时警告但不退出）
	log.Logger.Info("正在初始化邮件服务...")
	if err := mail.InitMailService(); err != nil {
		log.Logger.Warn("邮件服务初始化失败，将影响验证码发送功能", zap.Error(err))
	} else {
		log.Logger.Info("邮件服务初始化成功")
	}

	// 创建gin实例
	if !config.IsDebugMode() {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()

	// 或者只暴露这一个文件（更安全，推荐）
	r.StaticFile("/MP_verify_PpANgPfV9JVRT3XV.txt", "./static/MP_verify_PpANgPfV9JVRT3XV.txt")

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
		} else if err == nil && serverStarted {
			log.Logger.Info("HTTP服务器已正常关闭")
		}
	}()

	// 标记服务器已启动
	serverStarted = true
	log.Logger.Info("服务器启动成功", zap.String("addr", server.Addr))

	// 等待中断信号优雅关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	sig := <-quit

	log.Logger.Info("收到关闭信号，正在关闭服务器...", zap.String("signal", sig.String()))

	// 给予10秒超时时间优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Logger.Error("服务器强制关闭", zap.Error(err))
	} else {
		log.Logger.Info("服务器已优雅退出")
	}
}
