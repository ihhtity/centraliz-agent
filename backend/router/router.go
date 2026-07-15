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
		// 集控设备相关
		api.POST("/device/add", controller.AddDeviceControl)
		api.POST("/device/common", controller.DeviceCommon)

		// 设备控制日志相关
		api.POST("/devicelog/cablelogs", controller.CreateDeviceLog)
		// 微信相关
		api.POST("/wechat/login", controller.WXLogin)
		api.GET("/wechat/callback", controller.WXMPLoginCallback)
		api.GET("/wechat/jssdk/config", controller.WXJSSDKConfig)
		api.POST("/wechat/scan/image", controller.WXScanImage)
		// api.POST("/wechat/userinfo", controller.GetWechatUserInfo)
		// api.POST("/wechat/bind", controller.BindWechatUser)
		// api.GET("/wechat/unbind", controller.UnbindWechatUser)
		// api.POST("/wechat/update", controller.UpdateWechatUserInfo)
		// api.POST("/wechat/wechataccount", controller.GetWechatAccountInfo)
		// 通用功能 - 用户和商家共用
		api.POST("/send-code", controller.SendCode)
		// 用户相关
		api.POST("/user/login", controller.UserLogin)
		api.POST("/user/register", controller.UserRegister)
		// 商家相关
		api.POST("/merch/login", controller.MerchLogin)
		api.POST("/merch/register", controller.MerchRegister)
		api.POST("/merch/reset-password", controller.MerchResetPassword)
		// 汇付天下API路由（公共路由，不需要认证）
		huifuAPI := api.Group("/huifu")
		{
			// 	// 商户进件
			// 	merchant := huifuAPI.Group("/merchant")
			// 	{
			// 		merchant.POST("/register/enterprise", (&controller.HuifuAPIController{}).RegisterEnterprise)
			// 		merchant.POST("/register/personal", (&controller.HuifuAPIController{}).RegisterPersonal)
			// 		merchant.PUT("/modify", (&controller.HuifuAPIController{}).ModifyMerchant)
			// 	}

			// 	// 扫码支付
			qrpay := huifuAPI.Group("/qrpay")
			{
				qrpay.POST("/h5/wechat", (&controller.HuifuController{}).WxH5Pay)
				qrpay.POST("/mini/wechat", (&controller.HuifuController{}).Wxminipay)
				qrpay.POST("/jsapi/wechat", (&controller.HuifuController{}).WxJsApiPay)
				// 		qrpay.POST("/h5/alipay", (&controller.HuifuAPIController{}).CreateH5AlipayPay)
				// 		qrpay.POST("/mini/wechat/extend", (&controller.HuifuAPIController{}).CreateMiniWechatPayWithExtend)
				// 		qrpay.POST("/mini/alipay", (&controller.HuifuAPIController{}).CreateMiniAlipayPay)
				// 		qrpay.POST("/query", (&controller.HuifuAPIController{}).QueryQrPay)
				// 		qrpay.GET("/query/:trans_id", (&controller.HuifuAPIController{}).QueryQrPayByTransId)
				// 		qrpay.POST("/refund", (&controller.HuifuAPIController{}).QrPayRefund)
				// 		qrpay.POST("/refund/query", (&controller.HuifuAPIController{}).QrPayRefundQuery)
			}

			// 	// 分账相关
			// 	profit := huifuAPI.Group("/profit")
			// 	{
			// 		profit.POST("/share", (&controller.HuifuAPIController{}).ProfitShare)
			// 		profit.GET("/query/:profit_id", (&controller.HuifuAPIController{}).QueryProfit)
			// 	}

			// 	// 延时交易相关
			// 	delayed := huifuAPI.Group("/delayed")
			// 	{
			// 		delayed.POST("/confirm", (&controller.HuifuAPIController{}).DelayedConfirm)
			// 		delayed.POST("/confirm/query", (&controller.HuifuAPIController{}).DelayedConfirmQuery)
			// 		delayed.POST("/refund", (&controller.HuifuAPIController{}).DelayedRefund)
			// 		delayed.POST("/refund/query", (&controller.HuifuAPIController{}).DelayedRefundQuery)
			// 	}

			// 	// 支付回调相关
			// 	callback := huifuAPI.Group("/callback")
			// 	{
			// 		callback.POST("/payment", (&controller.HuifuAPIController{}).PaymentCallback)
			// 		callback.POST("/payment/parse", (&controller.HuifuAPIController{}).PaymentCallbackParse)
			// 	}
		}

		// 需要认证的路由
		auth := api.Group("")
		auth.Use(middleware.JWTAuth())
		{
			// 用户相关路由
			user := auth.Group("/user")
			{
				// 用户端个人信息相关路由
				user.GET("/profile/:id", controller.GetProfile)           // 获取用户个人信息
				user.PUT("/profile", controller.UpdateProfile)            // 更新用户个人信息
				user.PUT("/profile/email", controller.UserBindEmail)      // 绑定用户邮箱
				user.DELETE("/profile/email", controller.UserUnbindEmail) // 解绑用户邮箱
				// 用户端柜子相关路由
				user.GET("/room/:id", controller.GetUserRoomDetail)                       // 获取单个柜子详情
				user.GET("/room/group/:groupId", controller.GetUserRoomListByGroup)       // 获取分组下的柜子列表
				user.GET("/room/merchant/:merchId", controller.GetUserRoomListByMerchant) // 获取商家全部分组下的柜子列表
				// 用户端订单操作路由
				user.POST("/order/list", controller.GetUserOrderList)      // 获取用户订单列表
				user.GET("/order/:id", controller.GetUserOrderDetail)      // 获取单个订单详情
				user.PUT("/order/:id/refund", controller.UserApplyRefund)  // 申请退款
				user.POST("/order/complete", controller.UserCompleteOrder) // 完成订单
				user.POST("/order/create", controller.UserCreateOrder)     // 创建订单
				user.POST("/order/payment", controller.UserOrderPayment)   // 订单支付
				user.POST("/order/renew", controller.UserOrderRenew)       // 订单续费
				user.POST("/order/end", controller.UserOrderEnd)           // 结束订单

				// 用户端押金相关路由
				user.POST("/deposit/check", controller.CheckDepositStatus) // 检查押金状态
				user.POST("/deposit/pay", controller.PayDeposit)           // 支付押金
				user.POST("/deposit/refund", controller.RefundDeposit)     // 退还押金
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
				device.POST("/bind-group", controller.BindDeviceToGroup)
				device.POST("/unbind-group", controller.UnbindDeviceFromGroup)

				// 设备日志相关路由
				device.POST("/log", controller.CreateDeviceLog)
				device.GET("/log/list", controller.GetDeviceLogList)
				device.GET("/log/:id", controller.GetDeviceLogDetail)
				device.PUT("/log/:id", controller.UpdateDeviceLog)
				device.DELETE("/log/:id", controller.DeleteDeviceLog)
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
				room.GET("/:id/qrcode", controller.GenerateQRCode)
			}

			// 订单相关路由
			order := auth.Group("/order")
			{
				order.GET("/list", controller.GetOrderList)
				order.GET("/:id", controller.GetOrderDetail)
				order.POST("", controller.CreateOrder)
				order.PUT("/:id", controller.UpdateOrder)
				order.DELETE("/:id", controller.DeleteOrder)

				// 退款相关路由
				order.GET("/refund/list", controller.GetRefundList)
				order.PUT("/:id/refund/approve", controller.ApproveRefund)
				order.PUT("/:id/refund/reject", controller.RejectRefund)

				// 押金相关路由
				order.GET("/deposit/list", controller.GetDepositList)
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
				submerch.GET("/list", (&controller.SubMerchController{}).GetList)
				submerch.GET("/detail", (&controller.SubMerchController{}).GetDetail)
				submerch.POST("", (&controller.SubMerchController{}).Create)
				submerch.PUT("", (&controller.SubMerchController{}).Update)
				submerch.DELETE("/:id", (&controller.SubMerchController{}).Delete)
			}

			// 商家消费订单相关路由
			merchPay := auth.Group("/merch-pay")
			{
				merchPay.POST("/create", controller.CreateMerchPay)
				merchPay.GET("/list", controller.GetMerchPayList)
				merchPay.GET("/detail", controller.GetMerchPayDetail)
				merchPay.POST("/cancel", controller.CancelMerchPay)
				merchPay.POST("/pay", controller.PayMerchPay)
			}

			// 规则相关路由
			rule := auth.Group("/rule")
			{
				rule.GET("/list", controller.GetRuleList)
				rule.GET("/:id", controller.GetRuleDetail)
				rule.POST("", controller.CreateRule)
				rule.PUT("/:id", controller.UpdateRule)
				rule.DELETE("/:id", controller.DeleteRule)
			}
		}
	}
}
