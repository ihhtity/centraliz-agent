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
		api.POST("/device/add", controller.AddDeviceControl) // 添加集控设备
		api.POST("/device/common", controller.DeviceCommon)  // 集控设备通用接口

		// 设备控制日志相关
		api.POST("/devicelog/cablelogs", controller.CreateDeviceLog) // 创建集控设备日志
		// 微信相关
		api.POST("/wechat/login", controller.WXLogin)             // 微信登录
		api.GET("/wechat/callback", controller.WXMPLoginCallback) // 微信登录回调
		api.GET("/wechat/jssdk/config", controller.WXJSSDKConfig) // 获取JS SDK配置
		api.POST("/wechat/scan/image", controller.WXScanImage)    // 扫描二维码
		// api.POST("/wechat/userinfo", controller.GetWechatUserInfo) // 获取微信用户信息
		// api.POST("/wechat/bind", controller.BindWechatUser) // 绑定微信用户
		// api.GET("/wechat/unbind", controller.UnbindWechatUser) // 解绑微信用户
		// api.POST("/wechat/update", controller.UpdateWechatUserInfo) // 更新微信用户信息
		// api.POST("/wechat/wechataccount", controller.GetWechatAccountInfo) // 获取微信账号信息
		// 通用功能 - 用户和商家共用
		api.POST("/send-code", controller.SendCode) // 发送验证码
		// 用户相关
		api.POST("/user/login", controller.UserLogin)       // 用户登录
		api.POST("/user/register", controller.UserRegister) // 用户注册
		// 商家相关
		api.POST("/merch/login", controller.MerchLogin)                  // 商家登录
		api.POST("/merch/register", controller.MerchRegister)            // 商家注册
		api.POST("/merch/reset-password", controller.MerchResetPassword) // 商家重置密码
		// 汇付天下API路由（公共路由，不需要认证）
		huifuAPI := api.Group("/huifu")
		{
			// 	// 商户进件
			// 	merchant := huifuAPI.Group("/merchant")
			// 	{
			// 		merchant.POST("/register/enterprise", (&controller.HuifuAPIController{}).RegisterEnterprise) // 企业商户进件
			// 		merchant.POST("/register/personal", (&controller.HuifuAPIController{}).RegisterPersonal) // 个人商户进件
			// 		merchant.PUT("/modify", (&controller.HuifuAPIController{}).ModifyMerchant) // 修改商户信息
			// 	}

			// 	// 扫码支付
			qrpay := huifuAPI.Group("/qrpay")
			{
				qrpay.POST("/h5/wechat", (&controller.HuifuController{}).WxH5Pay)       // H5微信支付
				qrpay.POST("/mini/wechat", (&controller.HuifuController{}).Wxminipay)   // 小程序支付支付
				qrpay.POST("/jsapi/wechat", (&controller.HuifuController{}).WxJsApiPay) // JS API支付
				// 		qrpay.POST("/h5/alipay", (&controller.HuifuAPIController{}).CreateH5AlipayPay) // H5支付宝支付
				// 		qrpay.POST("/mini/wechat/extend", (&controller.HuifuAPIController{}).CreateMiniWechatPayWithExtend) // 小程序支付扩展
				// 		qrpay.POST("/mini/alipay", (&controller.HuifuAPIController{}).CreateMiniAlipayPay) // 小程序支付宝支付
				// 		qrpay.POST("/query", (&controller.HuifuAPIController{}).QueryQrPay) // 查询支付订单
				// 		qrpay.GET("/query/:trans_id", (&controller.HuifuAPIController{}).QueryQrPayByTransId) // 根据交易ID查询支付订单
				// 		qrpay.POST("/refund", (&controller.HuifuAPIController{}).QrPayRefund) // 退款
				// 		qrpay.POST("/refund/query", (&controller.HuifuAPIController{}).QrPayRefundQuery) // 查询退款订单
			}

			// 	// 分账相关
			// 	profit := huifuAPI.Group("/profit")
			// 	{
			// 		profit.POST("/share", (&controller.HuifuAPIController{}).ProfitShare) // 分账
			// 		profit.GET("/query/:profit_id", (&controller.HuifuAPIController{}).QueryProfit) // 查询分账订单
			// 	}

			// 	// 延时交易相关
			// 	delayed := huifuAPI.Group("/delayed")
			// 	{
			// 		delayed.POST("/confirm", (&controller.HuifuAPIController{}).DelayedConfirm) // 确认延时交易
			// 		delayed.POST("/confirm/query", (&controller.HuifuAPIController{}).DelayedConfirmQuery) // 查询延时交易确认订单
			// 		delayed.POST("/refund", (&controller.HuifuAPIController{}).DelayedRefund) // 退款
			// 		delayed.POST("/refund/query", (&controller.HuifuAPIController{}).DelayedRefundQuery) // 查询退款订单
			// 	}

			// 	// 支付回调相关
			// 	callback := huifuAPI.Group("/callback")
			// 	{
			// 		callback.POST("/payment", (&controller.HuifuAPIController{}).PaymentCallback) // 支付回调
			// 		callback.POST("/payment/parse", (&controller.HuifuAPIController{}).PaymentCallbackParse) // 支付回调解析
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
				merch.GET("/profile", controller.GetMerchProfile)          // 获取商家个人信息
				merch.POST("/profile/password", controller.ChangePassword) // 修改商家密码
				merch.PUT("/profile/email", controller.BindEmail)          // 绑定商家邮箱
				merch.DELETE("/profile/email", controller.UnbindEmail)     // 解绑商家邮箱
				merch.PUT("/profile/phone", controller.BindPhone)          // 绑定商家手机号
				merch.DELETE("/profile/phone", controller.UnbindPhone)     // 解绑商家手机号
			}

			// 设备相关路由
			device := auth.Group("/device")
			{
				device.GET("/list", controller.GetDeviceList)                  // 获取设备列表
				device.GET("/:id", controller.GetDeviceDetail)                 // 获取设备详情
				device.POST("", controller.CreateDevice)                       // 创建设备
				device.PUT("/:id", controller.UpdateDevice)                    // 更新设备
				device.DELETE("/:id", controller.DeleteDevice)                 // 删除设备
				device.POST("/bind-group", controller.BindDeviceToGroup)       // 绑定设备到分组
				device.POST("/unbind-group", controller.UnbindDeviceFromGroup) // 解绑设备从分组

				// 设备日志相关路由
				device.POST("/log", controller.CreateDeviceLog)       // 创建设备日志
				device.GET("/log/list", controller.GetDeviceLogList)  // 获取设备日志列表
				device.GET("/log/:id", controller.GetDeviceLogDetail) // 获取设备日志详情
				device.PUT("/log/:id", controller.UpdateDeviceLog)    // 更新设备日志
				device.DELETE("/log/:id", controller.DeleteDeviceLog) // 删除设备日志
			}

			// 分组相关路由
			group := auth.Group("/group")
			{
				group.GET("/list", controller.GetGroupList)  // 获取分组列表
				group.GET("/:id", controller.GetGroupDetail) // 获取分组详情
				group.POST("", controller.CreateGroup)       // 创建分组
				group.PUT("/:id", controller.UpdateGroup)    // 更新分组
				group.DELETE("/:id", controller.DeleteGroup) // 删除分组
			}

			// 房间相关路由
			room := auth.Group("/room")
			{
				room.GET("/list", controller.GetRoomList)          // 获取房间列表
				room.GET("/:id", controller.GetRoomDetail)         // 获取房间详情
				room.POST("", controller.CreateRoom)               // 创建房间
				room.PUT("/:id", controller.UpdateRoom)            // 更新房间
				room.DELETE("/:id", controller.DeleteRoom)         // 删除房间
				room.GET("/:id/qrcode", controller.GenerateQRCode) // 生成房间二维码

				// 房间标签相关路由
				tag := room.Group("/tag")
				{
					tag.GET("/list/:merchId", controller.GetRoomTagList) // 获取房间标签列表
					tag.POST("", controller.CreateRoomTag)               // 创建房间标签
					tag.PUT("/:id", controller.UpdateRoomTag)            // 更新房间标签
					tag.DELETE("/:id", controller.DeleteRoomTag)         // 删除房间标签
				}

				// 房间图片相关路由
				image := room.Group("/image")
				{
					image.GET("/list", controller.GetRoomImageList)   // 获取房间图片列表
					image.POST("", controller.CreateRoomImage)        // 创建房间图片
					image.PUT("/:id", controller.UpdateRoomImage)     // 更新房间图片
					image.DELETE("/:id", controller.DeleteRoomImage)  // 删除房间图片
					image.POST("/upload", controller.UploadRoomImage) // 上传房间图片
				}
			}

			// 订单相关路由
			order := auth.Group("/order")
			{
				order.GET("/list", controller.GetOrderList)  // 获取订单列表
				order.GET("/:id", controller.GetOrderDetail) // 获取订单详情
				order.POST("", controller.CreateOrder)       // 创建订单
				order.PUT("/:id", controller.UpdateOrder)    // 更新订单
				order.DELETE("/:id", controller.DeleteOrder) // 删除订单

				// 退款相关路由
				order.GET("/refund/list", controller.GetRefundList)        // 获取退款列表
				order.PUT("/:id/refund/approve", controller.ApproveRefund) // 审批退款
				order.PUT("/:id/refund/reject", controller.RejectRefund)   // 拒绝退款

				// 押金相关路由
				order.GET("/deposit/list", controller.GetDepositList) // 获取押金列表
			}

			// 收款账号相关路由
			huifu := auth.Group("/huifu")
			{
				huifu.GET("/list", (&controller.HuifuController{}).GetList)     // 获取收款账号列表
				huifu.GET("/:id", (&controller.HuifuController{}).GetDetail)    // 获取收款账号详情
				huifu.POST("", (&controller.HuifuController{}).Create)          // 创建收款账号
				huifu.PUT("", (&controller.HuifuController{}).Update)           // 更新收款账号
				huifu.DELETE("/:id", (&controller.HuifuController{}).Delete)    // 删除收款账号
				huifu.PUT("/choose", (&controller.HuifuController{}).SetChoose) // 选择收款账号
			}

			// 子账号相关路由
			submerch := auth.Group("/submerch")
			{
				submerch.GET("/list", (&controller.SubMerchController{}).GetList)     // 获取子账号列表
				submerch.GET("/detail", (&controller.SubMerchController{}).GetDetail) // 获取子账号详情
				submerch.POST("", (&controller.SubMerchController{}).Create)          // 创建子账号
				submerch.PUT("", (&controller.SubMerchController{}).Update)           // 更新子账号
				submerch.DELETE("/:id", (&controller.SubMerchController{}).Delete)    // 删除子账号
			}

			// 商家消费订单相关路由
			merchPay := auth.Group("/merch-pay")
			{
				merchPay.POST("/create", controller.CreateMerchPay)   // 创建商家消费订单
				merchPay.GET("/list", controller.GetMerchPayList)     // 获取商家消费订单列表
				merchPay.GET("/detail", controller.GetMerchPayDetail) // 获取商家消费订单详情
				merchPay.POST("/cancel", controller.CancelMerchPay)   // 取消商家消费订单
				merchPay.POST("/pay", controller.PayMerchPay)         // 支付商家消费订单
			}

			// 规则相关路由
			rule := auth.Group("/rule")
			{
				rule.GET("/list", controller.GetRuleList)  // 获取规则列表
				rule.GET("/:id", controller.GetRuleDetail) // 获取规则详情
				rule.POST("", controller.CreateRule)       // 创建规则
				rule.PUT("/:id", controller.UpdateRule)    // 更新规则
				rule.DELETE("/:id", controller.DeleteRule) // 删除规则
			}

			// PC管理后台路由
			admin := auth.Group("/admin")
			{
				// 首页统计
				admin.GET("/stats", controller.AdminGetDashboardStats)
				admin.GET("/stats/trend", controller.AdminGetTrendStats)

				// 房间管理
				adminRoom := admin.Group("/room")
				{
					adminRoom.GET("/list", controller.AdminGetRoomList)              // 获取房间列表(搜索分页)
					adminRoom.GET("/:id", controller.AdminGetRoomDetail)             // 获取房间详情
					adminRoom.POST("", controller.AdminCreateRoom)                   // 创建房间
					adminRoom.PUT("/:id", controller.AdminUpdateRoom)                // 更新房间
					adminRoom.DELETE("/:id", controller.AdminDeleteRoom)             // 删除房间
					adminRoom.POST("/batch-delete", controller.AdminBatchDeleteRoom) // 批量删除房间
					adminRoom.POST("/batch-update", controller.AdminBatchUpdateRoom) // 批量更新房间
					adminRoom.POST("/import", controller.AdminImportRoom)            // 导入房间
				}

				// 设备管理
				adminDevice := admin.Group("/device")
				{
					adminDevice.GET("/list", controller.AdminGetDeviceList)              // 获取设备列表(搜索分页)
					adminDevice.GET("/:id", controller.AdminGetDeviceDetail)             // 获取设备详情
					adminDevice.POST("", controller.AdminCreateDevice)                   // 创建设备
					adminDevice.PUT("/:id", controller.AdminUpdateDevice)                // 更新设备
					adminDevice.DELETE("/:id", controller.AdminDeleteDevice)             // 删除设备
					adminDevice.POST("/batch-delete", controller.AdminBatchDeleteDevice) // 批量删除设备
					adminDevice.POST("/batch-update", controller.AdminBatchUpdateDevice) // 批量更新设备
					adminDevice.POST("/import", controller.AdminImportDevice)            // 导入设备
				}

				// 分组管理
				adminGroup := admin.Group("/group")
				{
					adminGroup.GET("/list", controller.AdminGetGroupList)              // 获取分组列表(搜索分页)
					adminGroup.GET("/:id", controller.AdminGetGroupDetail)             // 获取分组详情
					adminGroup.POST("", controller.AdminCreateGroup)                   // 创建分组
					adminGroup.PUT("/:id", controller.AdminUpdateGroup)                // 更新分组
					adminGroup.DELETE("/:id", controller.AdminDeleteGroup)             // 删除分组
					adminGroup.POST("/batch-delete", controller.AdminBatchDeleteGroup) // 批量删除分组
					adminGroup.POST("/batch-update", controller.AdminBatchUpdateGroup) // 批量更新分组
					adminGroup.POST("/import", controller.AdminImportGroup)            // 导入分组
				}

				// 规则管理
				adminRule := admin.Group("/rule")
				{
					adminRule.GET("/list", controller.AdminGetRuleList)              // 获取规则列表(搜索分页)
					adminRule.GET("/:id", controller.AdminGetRuleDetail)             // 获取规则详情
					adminRule.POST("", controller.AdminCreateRule)                   // 创建规则
					adminRule.PUT("/:id", controller.AdminUpdateRule)                // 更新规则
					adminRule.DELETE("/:id", controller.AdminDeleteRule)             // 删除规则
					adminRule.POST("/batch-delete", controller.AdminBatchDeleteRule) // 批量删除规则
					adminRule.POST("/batch-update", controller.AdminBatchUpdateRule) // 批量更新规则
					adminRule.POST("/import", controller.AdminImportRule)            // 导入规则
				}

				// 订单管理
				adminOrder := admin.Group("/order")
				{
					adminOrder.GET("/list", controller.AdminGetOrderList)              // 获取订单列表(搜索分页)
					adminOrder.GET("/:id", controller.AdminGetOrderDetail)             // 获取订单详情
					adminOrder.PUT("/:id", controller.AdminUpdateOrder)                // 更新订单
					adminOrder.DELETE("/:id", controller.AdminDeleteOrder)             // 删除订单
					adminOrder.POST("/batch-delete", controller.AdminBatchDeleteOrder) // 批量删除订单
					adminOrder.POST("/batch-update", controller.AdminBatchUpdateOrder) // 批量更新订单
				}

				// 商家管理
				adminMerch := admin.Group("/merch")
				{
					adminMerch.GET("/list", controller.AdminGetMerchList)              // 获取商家列表(搜索分页)
					adminMerch.GET("/:id", controller.AdminGetMerchDetail)             // 获取商家详情
					adminMerch.POST("", controller.AdminCreateMerch)                   // 创建商家
					adminMerch.PUT("/:id", controller.AdminUpdateMerch)                // 更新商家
					adminMerch.DELETE("/:id", controller.AdminDeleteMerch)             // 删除商家
					adminMerch.POST("/batch-delete", controller.AdminBatchDeleteMerch) // 批量删除商家
					adminMerch.POST("/batch-update", controller.AdminBatchUpdateMerch) // 批量更新商家
					adminMerch.POST("/import", controller.AdminImportMerch)            // 导入商家
				}

				// 设备日志管理
				adminDeviceLog := admin.Group("/devicelog")
				{
					adminDeviceLog.GET("/list", controller.AdminGetDeviceLogList)              // 获取设备日志列表(搜索分页)
					adminDeviceLog.GET("/:id", controller.AdminGetDeviceLogDetail)             // 获取设备日志详情
					adminDeviceLog.DELETE("/:id", controller.AdminDeleteDeviceLog)             // 删除设备日志
					adminDeviceLog.POST("/batch-delete", controller.AdminBatchDeleteDeviceLog) // 批量删除设备日志
				}

				// 汇付账号管理
				adminHuifu := admin.Group("/huifu")
				{
					adminHuifu.GET("/list", controller.AdminGetHuifuAccountList)              // 获取汇付账号列表(搜索分页)
					adminHuifu.GET("/:id", controller.AdminGetHuifuAccountDetail)             // 获取汇付账号详情
					adminHuifu.POST("", controller.AdminCreateHuifuAccount)                   // 创建汇付账号
					adminHuifu.PUT("/:id", controller.AdminUpdateHuifuAccount)                // 更新汇付账号
					adminHuifu.DELETE("/:id", controller.AdminDeleteHuifuAccount)             // 删除汇付账号
					adminHuifu.POST("/batch-delete", controller.AdminBatchDeleteHuifuAccount) // 批量删除汇付账号
					adminHuifu.POST("/batch-update", controller.AdminBatchUpdateHuifuAccount) // 批量更新汇付账号
					adminHuifu.POST("/import", controller.AdminImportHuifuAccount)            // 导入汇付账号
				}

				// 商户支付管理
				adminMerchPay := admin.Group("/merchpay")
				{
					adminMerchPay.GET("/list", controller.AdminGetMerchPayList)              // 获取商户支付列表(搜索分页)
					adminMerchPay.GET("/:id", controller.AdminGetMerchPayDetail)             // 获取商户支付详情
					adminMerchPay.DELETE("/:id", controller.AdminDeleteMerchPay)             // 删除商户支付
					adminMerchPay.POST("/batch-delete", controller.AdminBatchDeleteMerchPay) // 批量删除商户支付
				}

				// 房间图片管理
				adminRoomImg := admin.Group("/roomimg")
				{
					adminRoomImg.GET("/list", controller.AdminGetRoomImageList)              // 获取房间图片列表(搜索分页)
					adminRoomImg.GET("/:id", controller.AdminGetRoomImageDetail)             // 获取房间图片详情
					adminRoomImg.POST("", controller.AdminCreateRoomImage)                   // 创建房间图片
					adminRoomImg.PUT("/:id", controller.AdminUpdateRoomImage)                // 更新房间图片
					adminRoomImg.DELETE("/:id", controller.AdminDeleteRoomImage)             // 删除房间图片
					adminRoomImg.POST("/batch-delete", controller.AdminBatchDeleteRoomImage) // 批量删除房间图片
					adminRoomImg.POST("/import", controller.AdminImportRoomImage)            // 导入房间图片
				}

				// 房间标签管理
				adminRoomTag := admin.Group("/roomtag")
				{
					adminRoomTag.GET("/list", controller.AdminGetRoomTagList)              // 获取房间标签列表(搜索分页)
					adminRoomTag.GET("/:id", controller.AdminGetRoomTagDetail)             // 获取房间标签详情
					adminRoomTag.POST("", controller.AdminCreateRoomTag)                   // 创建房间标签
					adminRoomTag.PUT("/:id", controller.AdminUpdateRoomTag)                // 更新房间标签
					adminRoomTag.DELETE("/:id", controller.AdminDeleteRoomTag)             // 删除房间标签
					adminRoomTag.POST("/batch-delete", controller.AdminBatchDeleteRoomTag) // 批量删除房间标签
					adminRoomTag.POST("/import", controller.AdminImportRoomTag)            // 导入房间标签
				}

				// 子商户管理
				adminSubMerch := admin.Group("/submerch")
				{
					adminSubMerch.GET("/list", controller.AdminGetSubMerchList)              // 获取子商户列表(搜索分页)
					adminSubMerch.GET("/:id", controller.AdminGetSubMerchDetail)             // 获取子商户详情
					adminSubMerch.POST("", controller.AdminCreateSubMerch)                   // 创建子商户
					adminSubMerch.PUT("/:id", controller.AdminUpdateSubMerch)                // 更新子商户
					adminSubMerch.DELETE("/:id", controller.AdminDeleteSubMerch)             // 删除子商户
					adminSubMerch.POST("/batch-delete", controller.AdminBatchDeleteSubMerch) // 批量删除子商户
					adminSubMerch.POST("/batch-update", controller.AdminBatchUpdateSubMerch) // 批量更新子商户
					adminSubMerch.POST("/import", controller.AdminImportSubMerch)            // 导入子商户
				}

				// 用户管理
				adminUser := admin.Group("/user")
				{
					adminUser.GET("/list", controller.AdminGetUserList)              // 获取用户列表(搜索分页)
					adminUser.GET("/:id", controller.AdminGetUserDetail)             // 获取用户详情
					adminUser.PUT("/:id", controller.AdminUpdateUser)                // 更新用户
					adminUser.DELETE("/:id", controller.AdminDeleteUser)             // 删除用户
					adminUser.POST("/batch-delete", controller.AdminBatchDeleteUser) // 批量删除用户
					adminUser.POST("/batch-update", controller.AdminBatchUpdateUser) // 批量更新用户
				}

				// 微信用户管理
				adminWxUser := admin.Group("/wxuser")
				{
					adminWxUser.GET("/list", controller.AdminGetWxUserList)              // 获取微信用户列表(搜索分页)
					adminWxUser.GET("/:id", controller.AdminGetWxUserDetail)             // 获取微信用户详情
					adminWxUser.PUT("/:id", controller.AdminUpdateWxUser)                // 更新微信用户
					adminWxUser.DELETE("/:id", controller.AdminDeleteWxUser)             // 删除微信用户
					adminWxUser.POST("/batch-delete", controller.AdminBatchDeleteWxUser) // 批量删除微信用户
				}
			}
		}
	}
}
