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
	// 根路径健康检查（用于负载均衡器/进程管理工具的健康检查）
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "centraliz-backend", "version": "1.0.0"})
	})

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

		// 微信相关
		api.POST("/wechat/login", controller.WechatLogin)
		api.POST("/wechat/userinfo", controller.GetWechatUserInfo)
		api.POST("/wechat/bind", controller.BindWechatUser)
		api.GET("/wechat/unbind", controller.UnbindWechatUser)
		api.POST("/wechat/update", controller.UpdateWechatUserInfo)

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
				merch.POST("/profile/password", controller.ChangePassword)
				merch.PUT("/profile/email", controller.BindEmail)
				merch.DELETE("/profile/email", controller.UnbindEmail)
				merch.PUT("/profile/phone", controller.BindPhone)
				merch.DELETE("/profile/phone", controller.UnbindPhone)
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
				room.POST("/:id/open-lock", controller.OpenLock)
				room.GET("/:id/qrcode", controller.GenerateQRCode)
			}

			// 订单相关路由
			order := auth.Group("/order")
			{
				order.GET("/list", controller.GetOrderList)
				order.GET("/:id", controller.GetOrderDetail)
				order.POST("", controller.CreateOrder)
				order.PUT("/:id", controller.UpdateOrder)
			}

			// 收款账号相关路由
			huifu := auth.Group("/huifu")
			{
				huifu.GET("/list", (&controller.HuifuController{}).GetList)
				huifu.GET("/:id", (&controller.HuifuController{}).GetDetail)
				huifu.POST("", (&controller.HuifuController{}).Create)
				huifu.PUT("", (&controller.HuifuController{}).Update)
				huifu.DELETE("/:id", (&controller.HuifuController{}).Delete)
				huifu.PUT("/choose", (&controller.HuifuController{}).SetChoose)
			}

			// 子账号相关路由
			submerch := auth.Group("/submerch")
			{
				submerch.GET("/list", controller.GetSubMerchList)
				submerch.GET("/:id", controller.GetSubMerchDetail)
				submerch.POST("", controller.CreateSubMerch)
				submerch.PUT("/:id", controller.UpdateSubMerch)
				submerch.DELETE("/:id", controller.DeleteSubMerch)
			}
		}
	}
}
