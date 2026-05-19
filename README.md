# Centraliz 中控管理系统开发文档

## 1. 项目概述

### 1.1 项目背景
Centraliz 是一个支持多语言、面向设备管理与用户使用的中后台管理系统，适用于需要集中化控制和管理硬件资源（如储物柜、房间、设备等）的场景。

### 1.2 产品定位
- **中控存柜产品**: 主要用于物品存储管理，支持个人/企业用户的存取需求
- **中控无人零售柜产品**: 主要用于商品零售，支持扫码购买、支付结算等功能

### 1.3 目标平台
- 微信小程序（用户端）
- 微信H5（管理端，支持公众号嵌入）
- 浏览器直接访问（管理端）
- 跨境App用户端（复用现有App）

## 2. 技术架构

### 2.1 前端技术栈
- **框架**: Vue 3 + Composition API
- **构建工具**: Vite 5.x
- **状态管理**: Pinia（需集成）
- **UI组件库**: @climblee/uv-ui
- **跨端框架**: Uni-app 3.x
- **国际化**: vue-i18n 9.x
- **网络请求**: 自定义封装（utils/request.js）

### 2.2 后端技术栈
- **Web框架**: Gin (Go)
- **数据库**: MySQL/PostgreSQL
- **缓存**: Redis
- **消息队列**: RabbitMQ/Kafka（可选）
- **文件存储**: 阿里云OSS/本地存储

### 2.3 系统架构图
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   用户端        │    │   管理端        │    │   硬件设备      │
│  (小程序/App)   │◄──►│   (H5/Web)      │◄──►│ (蓝牙/WiFi)     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                      │                      │
         ▼                      ▼                      ▼
┌─────────────────────────────────────────────────────────────┐
│                    Centraliz API Gateway                    │
└─────────────────────────────────────────────────────────────┘
         │                      │                      │
         ▼                      ▼                      ▼
┌───────────────┐    ┌─────────────────┐    ┌─────────────────┐
│  User Service │    │ Device Service  │    │ Payment Service │
└───────────────┘    └─────────────────┘    └─────────────────┘
         │                      │                      │
         ▼                      ▼                      ▼
┌─────────────────────────────────────────────────────────────┐
│                    Database & Cache Layer                   │
└─────────────────────────────────────────────────────────────┘
```

## 3. 功能模块设计

### 3.1 用户认证模块

#### 3.1.1 注册方式
- **手机号验证码注册**
  - 支持国际手机号格式
  - 验证码有效期：5分钟
  - 每日发送限制：10次
- **邮箱验证码注册**
  - 支持主流邮箱服务商
  - 验证码有效期：30分钟
  - 邮件模板支持多语言
- **账号密码注册**
  - 密码强度要求：8-20位，包含大小写字母、数字、特殊字符
  - 支持密码找回功能

#### 3.1.2 登录方式
- 手机号+验证码登录
- 邮箱+验证码登录
- 账号+密码登录
- 微信授权登录（小程序/H5）
- 第三方OAuth登录（可扩展）

#### 3.1.3 账号安全
- **密码保存功能**
  - 支持浏览器记住密码
  - 支持自动填充
  - 创建桌面快捷方式（PWA支持）
- **会话管理**
  - JWT Token认证
  - Refresh Token机制
  - 多设备登录管理
- **安全策略**
  - 密码错误次数限制
  - 异地登录提醒
  - 敏感操作二次验证

### 3.2 设备管理模块

#### 3.2.1 中控存柜功能
- **柜子状态监控**
  - 在线/离线状态
  - 电池电量显示
  - 门锁状态（开/关）
- **存取操作**
  - 扫码开柜
  - 预约存取
  - 临时授权
- **异常处理**
  - 超时未取提醒
  - 强制开柜
  - 故障上报

#### 3.2.2 中控无人零售柜功能
- **商品管理**
  - 商品上下架
  - 库存监控
  - 价格设置
- **购买流程**
  - 扫码选购
  - 支付集成（微信/支付宝）
  - 取货验证
- **补货管理**
  - 补货提醒
  - 补货记录
  - 销售统计

### 3.3 管理端功能

#### 3.3.1 H5管理端特性
- **响应式设计**
  - 适配手机/平板/PC
  - 触摸优化
- **公众号集成**
  - JS-SDK调用
  - 微信分享
  - 微信支付
- **蓝牙通信**
  - 蓝牙设备扫描
  - 蓝牙数据传输
  - WiFi配置推送

#### 3.3.2 核心管理功能
- **账户管理**
  - 用户列表
  - 权限分配
  - 操作日志
- **设备管理**
  - 设备分组
  - 批量操作
  - 远程控制
- **规则管理**
  - 计费规则
  - 使用规则
  - 报警规则
- **房间管理**
  - 房间布局
  - 设备关联
  - 状态监控

### 3.4 国际化支持
- **支持语言**
  - 中文（简体/繁体）
  - 英语、日语、韩语
  - 法语、德语、西班牙语
  - 俄语、泰语、土耳其语
  - 阿拉伯语、葡萄牙语、意大利语
- **本地化特性**
  - 日期时间格式
  - 数字货币格式
  - RTL布局支持（阿拉伯语）

## 4. 开发流程

### 4.1 环境搭建

#### 4.1.1 前端环境
```bash
# 进入前端目录
cd d:\c-full-stack\centraliz\frontend

# 安装依赖
npm install

# 安装Pinia（当前项目缺少）
npm install pinia @pinia-plugin-persistedstate

# 启动开发服务器
npm run dev:h5
```

#### 4.1.2 后端环境
```bash
# Go环境要求
Go 1.20+

# 项目结构
backend/
├── cmd/
├── internal/
│   ├── handler/
│   ├── service/
│   ├── repository/
│   └── model/
├── pkg/
└── config/

# 启动后端服务
go run cmd/server/main.go
```

### 4.2 前端开发流程

#### 4.2.1 Pinia集成步骤
1. **创建store目录**
   ```
   src/
   └── stores/
       ├── index.js
       ├── user.js
       ├── device.js
       └── app.js
   ```

2. **配置Pinia主文件**
   ```javascript
   // src/stores/index.js
   import { createPinia } from 'pinia'
   import piniaPluginPersistedstate from 'pinia-plugin-persistedstate'
   
   const pinia = createPinia()
   pinia.use(piniaPluginPersistedstate)
   
   export default pinia
   ```

3. **在main.js中集成**
   ```javascript
   // src/main.js
   import { createSSRApp } from "vue"
   import App from "./App.vue"
   import uvUI from '@climblee/uv-ui'
   import i18n from './locales/index'
   import { Request } from './utils/request'
   import pinia from './stores' // 新增
   
   export function createApp() {
     const app = createSSRApp(App)
     app.use(uvUI)
     app.use(i18n)
     app.use(pinia) // 新增
     Request(app)
     return { app }
   }
   ```

#### 4.2.2 蓝牙通信实现
```javascript
// src/utils/bluetooth.js
export class BluetoothManager {
  constructor() {
    this.device = null
    this.server = null
    this.characteristic = null
  }
  
  async connect(deviceId) {
    // 微信小程序蓝牙API
    if (uni.getSystemInfoSync().platform === 'devtools') {
      // 开发者工具模拟
      return true
    }
    
    try {
      await uni.closeBluetoothAdapter()
      await uni.openBluetoothAdapter()
      const devices = await this.scanDevices()
      const targetDevice = devices.find(d => d.deviceId === deviceId)
      
      if (targetDevice) {
        await uni.createBLEConnection({ deviceId })
        this.device = targetDevice
        await this.discoverServices()
        return true
      }
    } catch (error) {
      console.error('蓝牙连接失败:', error)
      return false
    }
  }
  
  async sendWiFiConfig(ssid, password) {
    const data = JSON.stringify({ ssid, password })
    const buffer = new TextEncoder().encode(data)
    
    await uni.writeBLECharacteristicValue({
      deviceId: this.device.deviceId,
      serviceId: this.server.serviceId,
      characteristicId: this.characteristic.characteristicId,
      value: buffer
    })
  }
}
```

#### 4.2.3 多语言注册表单组件
``vue
<!-- src/components/AuthForm.vue -->
<template>
  <view class="auth-form">
    <uv-tabs :list="tabList" :current="currentTab" @change="handleTabChange" />
    
    <!-- 手机号注册 -->
    <view v-if="currentTab === 0" class="form-content">
      <uv-input v-model="phone" placeholder="请输入手机号" />
      <view class="verification-code">
        <uv-input v-model="phoneCode" placeholder="验证码" />
        <uv-button 
          :disabled="!canSendPhoneCode" 
          @click="sendPhoneCode"
        >
          {{ phoneCodeBtnText }}
        </uv-button>
      </view>
    </view>
    
    <!-- 邮箱注册 -->
    <view v-if="currentTab === 1" class="form-content">
      <uv-input v-model="email" placeholder="请输入邮箱" />
      <view class="verification-code">
        <uv-input v-model="emailCode" placeholder="验证码" />
        <uv-button 
          :disabled="!canSendEmailCode" 
          @click="sendEmailCode"
        >
          {{ emailCodeBtnText }}
        </uv-button>
      </view>
    </view>
    
    <!-- 账号密码注册 -->
    <view v-if="currentTab === 2" class="form-content">
      <uv-input v-model="username" placeholder="用户名" />
      <uv-input v-model="password" password placeholder="密码" />
      <uv-input v-model="confirmPassword" password placeholder="确认密码" />
    </view>
    
    <uv-button @click="handleSubmit">注册</uv-button>
    <uv-checkbox v-model="rememberAccount">记住账号</uv-checkbox>
  </view>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const currentTab = ref(0)
const tabList = computed(() => [
  { name: t('phoneRegister') },
  { name: t('emailRegister') },
  { name: t('accountRegister') }
])

// 表单数据
const phone = ref('')
const phoneCode = ref('')
const email = ref('')
const emailCode = ref('')
const username = ref('')
const password = ref('')
const confirmPassword = ref('')
const rememberAccount = ref(false)

// 验证码倒计时
const phoneCodeCountdown = ref(0)
const emailCodeCountdown = ref(0)

const canSendPhoneCode = computed(() => phoneCodeCountdown.value === 0 && phone.value)
const canSendEmailCode = computed(() => emailCodeCountdown.value === 0 && email.value)

const phoneCodeBtnText = computed(() => 
  phoneCodeCountdown.value > 0 ? `${phoneCodeCountdown.value}s` : t('sendCode')
)
const emailCodeBtnText = computed(() => 
  emailCodeCountdown.value > 0 ? `${emailCodeCountdown.value}s` : t('sendCode')
)

// 发送验证码逻辑
const sendPhoneCode = async () => {
  // 调用API发送短信验证码
  phoneCodeCountdown.value = 60
  startCountdown('phone')
}

const sendEmailCode = async () => {
  // 调用API发送邮件验证码
  emailCodeCountdown.value = 60
  startCountdown('email')
}

const startCountdown = (type) => {
  const countdown = type === 'phone' ? phoneCodeCountdown : emailCodeCountdown
  const timer = setInterval(() => {
    if (countdown.value > 0) {
      countdown.value--
    } else {
      clearInterval(timer)
    }
  }, 1000)
}

const handleSubmit = async () => {
  // 根据当前tab提交不同类型的注册
  if (currentTab.value === 0) {
    // 手机号注册
  } else if (currentTab.value === 1) {
    // 邮箱注册
  } else {
    // 账号密码注册
  }
  
  // 保存账号信息（如果用户选择记住账号）
  if (rememberAccount.value) {
    // 保存到localStorage或Pinia持久化
  }
}
</script>
```

### 4.3 后端开发流程

#### 4.3.1 Gin路由设计
```go
// internal/handler/router.go
func RegisterRoutes(r *gin.Engine) {
    api := r.Group("/api/v1")
    
    // 认证相关
    auth := api.Group("/auth")
    {
        auth.POST("/register/phone", handler.RegisterByPhone)
        auth.POST("/register/email", handler.RegisterByEmail)
        auth.POST("/register/account", handler.RegisterByAccount)
        auth.POST("/login/phone", handler.LoginByPhone)
        auth.POST("/login/email", handler.LoginByEmail)
        auth.POST("/login/account", handler.LoginByAccount)
        auth.POST("/verify-code", handler.VerifyCode)
    }
    
    // 设备相关
    device := api.Group("/devices")
    {
        device.GET("/", handler.ListDevices)
        device.GET("/:id", handler.GetDevice)
        device.POST("/", handler.CreateDevice)
        device.PUT("/:id", handler.UpdateDevice)
        device.DELETE("/:id", handler.DeleteDevice)
        device.POST("/:id/control", handler.ControlDevice)
    }
    
    // 蓝牙通信
    bluetooth := api.Group("/bluetooth")
    {
        bluetooth.POST("/connect", handler.ConnectBluetooth)
        bluetooth.POST("/send-config", handler.SendWiFiConfig)
        bluetooth.GET("/status", handler.GetBluetoothStatus)
    }
}
```

#### 4.3.2 JWT认证中间件
```go
// internal/middleware/auth.go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := c.GetHeader("Authorization")
        if tokenString == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
            c.Abort()
            return
        }
        
        // 验证JWT token
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
            }
            return []byte(config.JWTSecret), nil
        })
        
        if err != nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }
        
        // 将用户信息存入上下文
        if claims, ok := token.Claims.(jwt.MapClaims); ok {
            c.Set("user_id", claims["user_id"])
            c.Set("role", claims["role"])
        }
        
        c.Next()
    }
}
```

### 4.4 微信H5特有功能开发

#### 4.4.1 公众号JS-SDK集成
``javascript
// src/utils/wechat.js
export class WechatSDK {
  async initJSSDK() {
    try {
      // 获取签名配置
      const config = await this.getSignatureConfig()
      
      // 初始化JS-SDK
      wx.config({
        debug: false,
        appId: config.appId,
        timestamp: config.timestamp,
        nonceStr: config.nonceStr,
        signature: config.signature,
        jsApiList: [
          'scanQRCode',
          'openLocation',
          'getLocation',
          'chooseImage',
          'uploadImage',
          'onMenuShareTimeline',
          'onMenuShareAppMessage'
        ]
      })
      
      wx.ready(() => {
        console.log('微信JS-SDK初始化成功')
      })
      
      wx.error((res) => {
        console.error('微信JS-SDK初始化失败:', res)
      })
    } catch (error) {
      console.error('获取微信配置失败:', error)
    }
  }
  
  async getSignatureConfig() {
    // 调用后端API获取签名配置
    const response = await uni.request({
      url: '/api/v1/wechat/signature',
      method: 'POST',
      data: {
        url: window.location.href.split('#')[0]
      }
    })
    return response.data
  }
  
  scanQRCode() {
    return new Promise((resolve, reject) => {
      wx.scanQRCode({
        needResult: 1,
        desc: 'scanQRCode desc',
        success: (res) => {
          resolve(res.resultStr)
        },
        fail: (err) => {
          reject(err)
        }
      })
    })
  }
}
```

#### 4.4.2 PWA快捷方式支持
``json
// public/manifest.json
{
  "name": "Centraliz管理端",
  "short_name": "Centraliz",
  "description": "中控存柜和无人零售柜管理平台",
  "start_url": "/",
  "display": "standalone",
  "background_color": "#ffffff",
  "theme_color": "#007AFF",
  "icons": [
    {
      "src": "/icons/icon-192x192.png",
      "sizes": "192x192",
      "type": "image/png"
    },
    {
      "src": "/icons/icon-512x512.png",
      "sizes": "512x512",
      "type": "image/png"
    }
  ]
}
```

```
// src/main.js - 添加PWA支持
if ('serviceWorker' in navigator) {
  window.addEventListener('load', () => {
    navigator.serviceWorker.register('/sw.js')
      .then(registration => {
        console.log('SW registered: ', registration)
      })
      .catch(registrationError => {
        console.log('SW registration failed: ', registrationError)
      })
  })
}
```

## 5. 部署方案

### 5.1 前端部署
- **H5管理端**: 部署到Nginx/Apache静态服务器
- **微信小程序**: 通过微信公众平台上传代码
- **跨境App**: 集成到现有App的WebView中

### 5.2 后端部署
- **生产环境**: Docker容器化部署
- **负载均衡**: Nginx反向代理
- **SSL证书**: Let's Encrypt免费证书
- **监控告警**: Prometheus + Grafana

### 5.3 环境配置
```
# docker-compose.yml
version: '3.8'
services:
  backend:
    build: ./backend
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=mysql
      - REDIS_HOST=redis
      - JWT_SECRET=your-secret-key
    depends_on:
      - mysql
      - redis
  
  mysql:
    image: mysql:8.0
    environment:
      - MYSQL_ROOT_PASSWORD=rootpassword
      - MYSQL_DATABASE=centraliz
    volumes:
      - mysql_data:/var/lib/mysql
  
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

volumes:
  mysql_data:
```

## 6. 测试策略

### 6.1 前端测试
- **单元测试**: Vitest + Vue Test Utils
- **E2E测试**: Cypress（H5）/ Uni-app自动化测试（小程序）
- **兼容性测试**: 多浏览器、多设备测试

### 6.2 后端测试
- **单元测试**: Go testing包
- **集成测试**: Testcontainers
- **API测试**: Postman/Newman

### 6.3 蓝牙功能测试
- **真机测试**: iOS/Android设备
- **模拟测试**: 微信开发者工具
- **边界测试**: 蓝牙断连、信号弱等情况

## 7. 安全考虑

### 7.1 数据安全
- HTTPS全站加密
- 敏感数据加密存储
- API请求签名验证

### 7.2 用户隐私
- GDPR合规
- 用户数据最小化原则
- 隐私政策透明化

### 7.3 硬件安全
- 蓝牙通信加密
- WiFi配置安全传输
- 设备固件签名验证

## 8. 性能优化

### 8.1 前端优化
- 代码分割（Code Splitting）
- 图片懒加载
- 本地缓存策略
- 骨架屏加载

### 8.2 后端优化
- 数据库索引优化
- Redis缓存热点数据
- 异步处理耗时操作
- CDN静态资源加速

## 9. 维护与监控

### 9.1 日志系统
- 前端错误监控（Sentry）
- 后端日志收集（ELK Stack）
- 操作审计日志

### 9.2 健康检查
- 服务健康状态API
- 数据库连接池监控
- 蓝牙设备在线状态监控

### 9.3 版本管理
- 语义化版本控制
- 自动化发布流程
- 回滚机制

## 10. 附录

### 10.1 API文档规范
使用Swagger/OpenAPI 3.0规范编写API文档

### 10.2 代码规范
- ESLint + Prettier（前端）
- gofmt + golangci-lint（后端）
- Git提交规范（Conventional Commits）

### 10.3 第三方服务集成
- 短信服务：阿里云短信/腾讯云短信
- 邮件服务：SendGrid/Mailgun
- 支付服务：微信支付/支付宝
- 文件存储：阿里云OSS/腾讯云COS

## 11. 数据库模型详细说明

### 11.1 表结构概览
根据实际MySQL数据库 `centraliz` 的表结构，系统包含以下10个核心表：
- `codes` - 验证码管理表
- `devices` - 设备信息表  
- `group` - 分组管理表
- `huifu_account` - 汇付账户表
- `merch` - 商家主账户表
- `merch_pay` - 商家支付订单表
- `orders` - 用户订单表
- `room` - 房间信息表
- `submerch` - 商家子账户表
- `user` - 用户信息表

### 11.2 模型字段映射关系

#### Codes (验证码表)
- **ID**: uint (主键，自增)
- **Phone**: *string (手机号，可为空，索引)
- **Email**: *string (邮箱，可为空，索引)  
- **Code**: string (验证码，非空，索引)
- **Type**: int8 (类型：1注册/2登录/3找回密码等)
- **Status**: int8 (状态：0未使用/1已使用)
- **ExpireAt**: time.Time (过期时间，索引)

#### Device (设备表)
- **ID**: uint64 (主键，bigint(20) unsigned，自增)
- **Name**: string (设备名称，非空)
- **MerchID**: int32 (商家外键，int(20))
- **RoomID**: *uint64 (房间外键，可为空，int(20) unsigned)
- **GroupID**: *uint64 (分组外键，可为空，int(20) unsigned)
- **OrdersID**: *int32 (订单外键，可为空，int(20))
- **Status**: string (状态：0空闲/1租用/2维修)
- **Type**: string (设备类型)

#### Group (分组表)
- **ID**: uint64 (主键，bigint(20) unsigned，自增)
- **Name**: *string (分组名称，可为空)
- **MerchID**: int32 (商家外键，int(11))
- **RuleID**: *int32 (规则外键，可为空，int(11))
- **Phone**: *string (客服手机号，可为空)
- **Count**: *uint32 (房间数量，可为空，int(10) unsigned)
- **Type**: *string (类型：0存柜/1零售，可为空)
- **Location**: *string (位置信息，可为空)

#### HuifuAccount (汇付账户表)
- **ID**: int32 (主键，int(11)，自增)
- **MerchID**: *int32 (商家ID，可为空)
- **Code**: *string (汇付编码，可为空)
- **Sharing**: string (多方分账，text类型)
- **Account**: *string (账号，可为空)
- **Phone**: *string (手机号，可为空)
- **Name**: *string (姓名，可为空)
- **Identity**: *string (身份证，可为空)
- **Card**: *string (银行卡，可为空)
- **Encrypt**: *string (营业执照编码，可为空)
- **Storename**: *string (店名，可为空)
- **Area**: *string (经营地址，text类型，可为空)
- **Picture**: string (店铺图片，text类型)
- **Remarks**: string (使用场景描述，text类型)
- **Type**: *string (账号类型，可为空)
- **Choose**: *string (选择状态：0未选择/1已选择)
- **Share**: string (分账开关：0关闭/1开启)
- **Rate**: *float64 (分账比率，可为空)

#### Merch (商家主账户表)
- **ID**: int32 (主键，int(11)，自增)
- **Account**: string (账号，非空)
- **Password**: string (密码，非空)
- **Email**: *string (邮箱，可为空)
- **Phone**: *string (手机号，可为空)
- **Role**: *string (角色：0商家/1管理者/2代理商)
- **Status**: *string (状态：0白名单/1黑名单)
- **LogAt**: *time.Time (最后登录时间，可为空)
- **CreatedAt**: *time.Time (创建时间，可为空)

#### MerchPay (商家支付订单表)
- **ID**: int32 (主键，int(11)，自增)
- **Code**: *string (订单号，可为空，唯一索引)
- **MerID**: *int32 (商家外键，可为空，索引)
- **Name**: *string (商品名称，可为空，索引)
- **Gateway**: *string (充值控电ID，可为空)
- **ReqDate**: *string (汇付支付时间，可为空)
- **HfSeqID**: *string (原交易全局流水号，可为空)
- **OriginalPrice**: *float64 (订单原价，decimal(10,2))
- **Price**: *float64 (实际支付金额，decimal(10,2))
- **LockTotal**: *int32 (锁数量，默认0)
- **Type**: *string (订单类型：0短信/1广告/2年费)
- **Status**: *string (订单状态：0待支付/1已支付/2已关闭)
- **Remarks**: *string (订单备注，可为空)
- **OpenID**: *string (小程序openid，可为空)
- **CreateTime**: BigIntTime (创建时间，bigint)
- **FinishTime**: BigIntTime (套餐结束时间，bigint)
- **DelectTime**: BigIntTime (删除时间，bigint)

#### Order (用户订单表)
- **ID**: uint64 (主键，bigint(20) unsigned，自增)
- **MerchID**: int32 (商家外键，非空)
- **UserID**: int32 (用户外键，非空)
- **DevicesID**: int32 (设备外键，非空)
- **RoomID**: int32 (房间外键，非空)
- **GroupID**: int32 (分组外键，非空)
- **Name**: *string (订单名称，可为空)
- **Code**: *string (订单编号，可为空)
- **Status**: *string (状态：0未支付/1已支付/3申请退款/4已退款/5拒绝退款)
- **Amount**: *float64 (商品数量，默认0)
- **Duration**: *int64 (使用时长，默认0)
- **Price**: *float64 (支付金额，decimal(10,2))
- **Deposit**: *float64 (押金，decimal(10,2))
- **StartTime**: *time.Time (订单开始时间，可为空)
- **EndTime**: *time.Time (订单结束时间，可为空)
- **ReqDate**: *string (支付时间日期，可为空)
- **UserPhone**: *string (用户手机号，可为空)
- **MerchPhone**: *string (商家手机号，可为空)
- **CreatedAt**: *time.Time (创建时间，可为空)

#### Room (房间表)
- **ID**: uint64 (主键，bigint(20) unsigned，自增)
- **Name**: *string (房间名称，可为空)
- **MerchID**: int32 (商家外键，非空，默认0)
- **GroupID**: *int32 (分组外键，可为空，默认0)
- **RuleID**: *int32 (规则外键，可为空，默认0)
- **Tag**: *string (房间标签，可为空)
- **Status**: *string (状态：0空闲/1租用/2维修)
- **CreatedAt**: *time.Time (创建时间，可为空)
- **UpdatedAt**: *time.Time (更新时间，可为空)

#### SubMerch (商家子账户表)
- **ID**: int32 (主键，int(11)，自增)
- **Account**: string (账号，非空)
- **Password**: string (密码，非空)
- **MerchID**: int32 (商家外键，非空)
- **Email**: *string (邮箱，可为空)
- **Phone**: *string (手机号，可为空)
- **Role**: *string (角色：0商家/1管理者/2代理商)
- **Status**: *string (状态：0白名单/1黑名单)
- **Rule**: *string (使用权限，text类型)
- **LogAt**: *time.Time (最后登录时间，可为空)
- **CreatedAt**: *time.Time (创建时间，可为空)

#### User (用户表)
- **ID**: uint64 (主键，bigint(20) unsigned，自增)
- **Name**: string (用户名，非空)
- **Account**: string (账号，非空)
- **Password**: string (密码，非空)
- **Email**: *string (邮箱，可为空)
- **Phone**: *string (手机号，可为空)
- **MerchID**: int32 (商家外键，非空，默认0)
- **RoomID**: *int32 (房间外键，可为空，索引)
- **OrdersID**: *int32 (订单外键，可为空，索引)
- **Status**: *string (状态：0白名单/1黑名单)
- **CreatedAt**: *time.Time (创建时间，可为空)
- **UpdatedAt**: *time.Time (更新时间，可为空)

### 11.3 数据类型映射规范
- **MySQL bigint(20) unsigned** → **Go uint64**
- **MySQL int(11)/int(20)** → **Go int32**
- **MySQL varchar(N)** → **Go string** (带gorm size标签)
- **MySQL text** → **Go string** (无size限制)
- **MySQL datetime/datetime(3)** → **Go time.Time**
- **MySQL decimal(M,N)** → **Go float64** (带gorm type:decimal(M,N)标签)
- **MySQL bigint** (时间戳) → **Go 自定义BigIntTime类型**

### 11.4 索引优化策略
- 主键字段自动创建PRIMARY KEY索引
- 外键字段创建普通索引以优化JOIN查询
- 高频查询字段（如手机号、邮箱、验证码）创建索引
- 唯一约束字段创建UNIQUE INDEX
- 可空字段使用指针类型(*string, *int32等)以正确映射NULL值

此数据库模型完整覆盖了Centraliz系统的业务需求，支持设备管理、用户认证、订单处理、商家管理、支付集成等核心功能。

## 12. 前端页面优化记录

### 12.1 登录页面优化
- **支持三种登录方式**：手机号验证码登录、邮箱验证码登录、账号密码登录
- **角色切换功能**：支持管理员和普通用户角色切换
- **忘记密码链接**：直接跳转到独立的找回密码页面
- **注册链接**：提供快速注册入口

### 12.2 注册页面新增
- **三种注册方式**：
  - 手机号验证码注册（需输入手机号、验证码、密码）
  - 邮箱验证码注册（需输入邮箱、验证码、密码）  
  - 账号密码注册（需输入账号、用户名、密码）
- **表单验证**：完整的字段格式验证和一致性检查
- **角色支持**：根据当前选择的角色进行注册
- **用户协议**：支持用户协议同意确认

### 12.3 忘记密码页面新增
- **两种找回方式**：手机号验证码找回、邮箱验证码找回
- **密码重置**：支持新密码设置和确认
- **表单验证**：完整的验证码和密码强度验证

### 12.4 多语言支持
- 所有新增页面均支持完整的中文翻译
- 字段提示、错误信息、按钮文本均有对应的多语言配置
- 支持未来扩展其他语言版本

### 12.5 路由配置
- 在pages.json中正确配置了所有新页面的路由
- 保持了原有的导航样式和页面标题
- 确保页面间的正常跳转和参数传递

这些优化确保了前端页面与后端数据库模型的完全匹配，提供了完整的用户认证体验。
