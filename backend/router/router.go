package router

import (
	"centraliz-backend/controller"
	"centraliz-backend/middleware"
	"centraliz-backend/pkg/db"
	"centraliz-backend/pkg/redis"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	// 健康检查端点
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "centraliz-backend"})
	})

	r.GET("/ready", func(c *gin.Context) {
		// 检查数据库连接
		if db.DB == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy", "reason": "database not initialized"})
			return
		}

		// 检查Redis连接
		if redis.RDB == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy", "reason": "redis not initialized"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "ready", "service": "centraliz-backend"})
	})

	// API版本路由
	api := r.Group("/api/v1")
	{
		// 公共路由（不需要认证）
		api.POST("/user/login", controller.Login)
		api.POST("/user/register", controller.Register)

		// 需要认证的路由
		auth := api.Group("")
		auth.Use(middleware.JWTAuth())
		{
			// 用户相关路由
			user := auth.Group("/user")
			{
				user.GET("/profile", controller.GetProfile)
				user.PUT("/profile", controller.UpdateProfile)
			}

			// 设备相关路由
			device := auth.Group("/device")
			{
				device.GET("/list", controller.GetDeviceList)
				device.GET("/:id", controller.GetDeviceDetail)
				device.POST("", controller.CreateDevice)
				device.PUT("/:id", controller.UpdateDevice)
				device.DELETE("/:id", controller.DeleteDevice)
			}

			// 房间相关路由
			room := auth.Group("/room")
			{
				room.GET("/list", controller.GetRoomList)
				room.GET("/:id", controller.GetRoomDetail)
				room.POST("", controller.CreateRoom)
				room.PUT("/:id", controller.UpdateRoom)
				room.DELETE("/:id", controller.DeleteRoom)
			}

			// 订单相关路由
			order := auth.Group("/order")
			{
				order.GET("/list", controller.GetOrderList)
				order.GET("/:id", controller.GetOrderDetail)
				order.POST("", controller.CreateOrder)
				order.PUT("/:id", controller.UpdateOrder)
			}

			// 管理员相关路由
			admin := auth.Group("/admin")
			{
				admin.GET("/dashboard", controller.GetDashboardStats)
				admin.GET("/accounts", controller.GetAccountList)
				admin.POST("/accounts", controller.CreateAccount)
			}
		}
	}
}
