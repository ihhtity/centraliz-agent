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
	api.Use(middleware.RateLimit()) // 所有API路由应用限流中间件
	{
		// 公共路由（不需要认证）
		// 通用功能 - 用户和商家共用
		api.POST("/send-code", controller.SendCode)

		// 用户相关
		api.POST("/user/login", controller.UserLogin)
		api.POST("/user/register", controller.UserRegister)
		api.POST("/user/reset-password", controller.UserResetPassword)

		// 商家相关
		api.POST("/merch/login", controller.MerchLogin)
		api.POST("/merch/register", controller.MerchRegister)
		api.POST("/merch/reset-password", controller.MerchResetPassword)

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

			// 商家相关路由
			merch := auth.Group("/merch")
			{
				merch.GET("/profile", controller.GetMerchProfile)
				merch.PUT("/profile", controller.UpdateMerchProfile)
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

			// 分组相关路由
			group := auth.Group("/group")
			{
				group.GET("/list", controller.GetGroupList)
				group.GET("/:id", controller.GetGroupDetail)
				group.POST("", controller.CreateGroup)
				group.PUT("/:id", controller.UpdateGroup)
				group.DELETE("/:id", controller.DeleteGroup)
			}

			// 房间相关路由
			room := auth.Group("/room")
			{
				room.GET("/list", controller.GetRoomList)
				room.GET("/:id", controller.GetRoomDetail)
				room.POST("", controller.CreateRoom)
				room.PUT("/:id", controller.UpdateRoom)
				room.DELETE("/:id", controller.DeleteRoom)
				room.POST("/bind-device", controller.BindDevice)
				room.POST("/unbind-device", controller.UnbindDevice)
			}

			// 订单相关路由
			order := auth.Group("/order")
			{
				order.GET("/list", controller.GetOrderList)
				order.GET("/:id", controller.GetOrderDetail)
				order.POST("", controller.CreateOrder)
				order.PUT("/:id", controller.UpdateOrder)
			}
		}
	}
}
