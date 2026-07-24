# Centraliz 中控管理系统

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

### 1.4 项目架构
```
┌─────────────────────────────────────────────────────────────┐
│                     Centraliz Agent                         │
├──────────────────┬──────────────────┬───────────────────────┤
│     admin        │    frontend      │       backend         │
│   (PC管理后台)     │  (H5/小程序)      │     (API服务)         │
│  React + TS      │  Vue3 + Uni-app  │    Go + Gin           │
│   Ant Design     │    uv-ui         │     GORM + Redis      │
│   ECharts        │    i18n          │    Huifu Pay SDK      │
└──────────────────┴──────────────────┴───────────────────────┘
```

### 1.5 项目目录结构

```
centraliz-agent/
├── .gitignore                    # Git忽略配置
├── README.md                     # 项目说明文档
├── AGENT_DEV.md                  # Agent智能开发助手开发文档
├── admin/                        # PC管理后台（React + TypeScript）
│   ├── index.html                # 入口HTML文件
│   ├── package.json              # 项目依赖配置
│   ├── vite.config.ts            # Vite配置文件
│   ├── tsconfig.json             # TypeScript配置文件
│   └── src/                      # 项目源代码目录
│       ├── components/           # 公共组件
│       │   ├── BatchActions.tsx    # 批量操作组件
│       │   ├── CustomPagination.tsx # 自定义分页组件
│       │   ├── ExportButton.tsx    # 导出按钮组件
│       │   └── SearchBar.tsx       # 搜索栏组件
│       ├── layout/               # 布局组件
│       │   ├── Layout.tsx          # 布局组件
│       │   └── Layout.scss         # 布局组件样式文件
│       ├── pages/                # 页面组件
│       │   ├── Dashboard.tsx       # 首页组件
│       │   ├── DeviceManage.tsx    # 设备管理组件
│       │   ├── DeviceLogManage.tsx # 设备日志管理组件
│       │   ├── DeviceLogManage.tsx # 设备日志管理组件
│       │   ├── GroupManage.tsx     # 分组管理组件
│       │   ├── MerchManage.tsx     # 商品管理组件
│       │   ├── MerchPayManage.tsx  # 商品支付管理组件
│       │   ├── OrderManage.tsx     # 订单管理组件
│       │   ├── RoomManage.tsx      # 房间管理组件
│       │   ├── RoomImgManage.tsx   # 房间图片管理组件
│       │   ├── RoomTagManage.tsx   # 房间标签管理组件
│       │   ├── RuleManage.tsx      # 规则管理组件
│       │   ├── SubMerchManage.tsx  # 子商品管理组件
│       │   ├── UserManage.tsx      # 用户管理组件
│       │   ├── WxUserManage.tsx    # 微信用户管理组件
│       │   ├── HuifuAccountManage.tsx # 恒付账号管理组件
│       │   ├── Login.tsx           # 登录组件
│       │   └── Login.scss          # 登录组件样式文件
│       ├── api/                  # API接口定义
│       │   └── index.ts            # API接口定义文件
│       ├── types/                # TypeScript类型定义
│       │   └── index.ts            # 类型定义文件入口
│       ├── utils/                # 工具函数
│       │   ├── request.ts          # HTTP请求工具函数
│       │   ├── export.ts           # 导出工具函数
│       │   └── pagination.ts       # 分页工具函数
│       ├── App.tsx                 # 应用入口组件
│       ├── main.tsx                # 应用入口文件组件
│       └── index.scss              # 全局样式文件
├── backend/                      # 后端服务（Go + Gin）
│   ├── main.go                   # 入口文件
│   ├── go.mod                    # 依赖依赖配置
│   ├── go.sum                    # 依赖依赖配置
│   ├── config/                   # 配置文件
│   │   ├── config.yaml           # 配置文件YAML格式
│   │   ├── access_token.json     # 访问令牌JSON文件
│   │   └── huifu.json            # 恒付账号JSON文件
│   ├── controller/               # 控制器层
│   │   ├── admin.go              # 管理员控制器文件
│   │   ├── admin_import.go     # 管理员导入控制器文件
│   │   ├── device.go           # 设备控制器文件
│   │   ├── devicelog.go        # 设备日志控制器文件
│   │   ├── group.go            # 分组控制器文件
│   │   ├── health.go           # 健康检查控制器文件
│   │   ├── huifu.go            # 汇付账号控制器文件
│   │   ├── merch.go            # 商品控制器文件
│   │   ├── merch_pay.go      # 商品支付控制器文件
│   │   ├── order.go          # 订单控制器文件
│   │   ├── room.go           # 房间控制器文件
│   │   ├── rule.go           # 规则控制器文件
│   │   ├── submerch.go       # 子商品控制器文件
│   │   ├── user.go           # 用户控制器文件
│   │   └── wechat.go         # 微信控制器文件
│   ├── logic/                    # 业务逻辑层
│   │   ├── device.go           # 设备业务逻辑文件
│   │   ├── devicelog.go        # 设备日志业务逻辑文件
│   │   ├── group.go            # 分组业务逻辑文件
│   │   ├── merch.go            # 商品业务逻辑文件
│   │   ├── order.go          # 订单业务逻辑文件
│   │   ├── room.go          # 房间业务逻辑文件
│   │   ├── rule.go           # 规则业务逻辑文件
│   │   └── user.go          # 用户业务逻辑文件
│   ├── dao/                      # 数据访问层
│   │   └── mysql/               # MySQL数据库访问层
│   │       ├── mysql.go           # MySQL数据库访问层文件
│   │       ├── device.go           # 设备数据库访问层文件
│   │       ├── devicelog.go        # 设备日志数据库访问层文件
│   │       ├── group.go            # 分组数据库访问层文件
│   │       ├── merch.go            # 商品数据库访问层文件
│   │       ├── order.go            # 订单数据库访问层文件
│   │       ├── room.go             # 房间数据库访问层文件
│   │       ├── rule.go             # 规则数据库访问层文件
│   │       └── user.go             # 用户数据库访问层文件
│   ├── model/                    # 数据模型
│   │   ├── device.go             # 设备数据模型文件
│   │   ├── devicelog.go          # 设备日志数据模型文件
│   │   ├── group.go              # 分组数据模型文件
│   │   ├── huifu_account.go      # 汇付账号数据模型文件
│   │   ├── merch.go              # 商品数据模型文件
│   │   ├── merch_pay.go          # 商品支付数据模型文件
│   │   ├── order.go              # 订单数据模型文件
│   │   ├── room.go             # 房间数据模型文件
│   │   ├── room_image.go       # 房间图片数据模型文件
│   │   ├── room_tag.go         # 房间标签数据模型文件
│   │   ├── rule.go             # 规则数据模型文件
│   │   ├── submerch.go         # 子商品数据模型文件
│   │   ├── user.go             # 用户数据模型文件
│   │   └── wechat_user.go      # 微信用户数据模型文件
│   ├── router/                   # 路由定义
│   │   └── router.go             # 路由定义文件
│   ├── middleware/               # 中间件
│   │   ├── auth.go               # 认证中间件文件
│   │   ├── cors.go               # CORS中间件文件
│   │   ├── logger.go             # 日志中间件文件
│   │   └── rate_limit.go         # 限流中间件文件
│   ├── pkg/                      # 公共包
│   │   ├── config/               # 配置文件包
│   │   ├── db/                   # 数据库包
│   │   ├── errno/                # 错误码包
│   │   ├── jwt/                  # JWT包
│   │   ├── log/                  # 日志包
│   │   ├── mail/                 # 邮件包
│   │   ├── redis/                # Redis包
│   │   ├── response/             # 响应包
│   │   └── utils/                # 工具包
│   ├── huifu/                    # 汇付天下支付SDK封装
│   │   ├── wxh5.go               # 汇付天下支付H5支付文件
│   │   ├── wxjsapi.go            # 汇付天下支付JSAPI支付文件
│   │   └── wxmini.go             # 汇付天下支付小程序支付文件
│   ├── static/                   # 静态文件
│   ├── uploads/                  # 上传文件目录
│   ├── start.bat                 # Windows启动脚本
│   └── start_daemon.bat          # Windows守护进程启动脚本
└── frontend/                     # 用户端（Vue3 + Uni-app）
    ├── index.html                # 用户端首页文件
    ├── package.json              # 用户端项目配置文件
    ├── vite.config.js            # 用户端Vite配置文件
    ├── vue.config.js             # 用户端Vue配置文件
    ├── shims-uni.d.ts            # 用户端Vue类型定义文件
    └── src/                      # 用户端源代码目录
        ├── components/           # 公共组件
        │   ├── BatteryIndicator.vue  # 电池指示器组件
        │   ├── LockerAgreement.vue   # 锁器协议组件
        │   └── RoomDetailPopup.vue   # 房间详情弹窗组件
        ├── pages/                # 页面组件
        │   ├── login/            # 登录相关
        │   │   ├── login.vue       # 登录页面
        │   │   ├── register.vue    # 注册页面
        │   │   └── forgot-password.vue # 忘记密码页面
        │   ├── user/             # 用户端页面
        │   │   ├── index/          # 用户端首页
        │   │   ├── order/          # 订单相关
        │   │   ├── profile/        # 用户信息相关
        │   │   ├── wallet/         # 余额相关
        │   │   └── language/       # 国际化配置（13种语言）
        │   └── admin/            # H5管理端页面
        │       ├── account/        # 账号相关
        │       ├── device/         # 设备相关
        │       ├── group/          # 分组相关
        │       ├── room/           # 房间相关
        │       ├── order/          # 订单相关
        │       ├── rule/           # 规则相关
        │       ├── huifu/          # 汇付相关
        │       ├── expense/        # 支出相关
        │       ├── value-added/    # 增值相关相关
        │       └── profile/        # 用户信息相关
        ├── locales/              # 国际化配置（13种语言）
        │   ├── index.js          # 默认语言配置
        │   ├── zh-CN.js          # 中文配置
        │   ├── zh-TW.js          # 繁体中文配置
        │   ├── en-US.js          # 英文配置
        │   ├── ja-JP.js          # 日文配置
        │   ├── ko-KR.js          # 韩文配置
        │   ├── fr-FR.js          # 法文配置
        │   ├── de-DE.js          # 德文配置
        │   ├── es-ES.js          # 西班牙文配置
        │   ├── it-IT.js          # 意大利文配置
        │   ├── pt-BR.js          # 葡萄牙文配置
        │   ├── ru-RU.js          # 俄文配置
        │   ├── tr-TR.js          # 土耳其文配置
        │   └── ar-SA.js          # 阿拉伯文配置
        ├── utils/                # 工具函数
        │   ├── request.js        # 网络请求工具函数
        │   └── utils.js          # 通用工具函数
        ├── static/               # 静态资源
        ├── App.vue               # 应用入口组件
        ├── main.js               # 应用入口文件
        ├── manifest.json         # 应用清单文件
        ├── pages.json            # 页面配置文件
        └── uni.scss              # 应用样式文件
```

## 2. 技术架构

### 2.1 PC管理后台 (admin)

| 分类 | 技术 | 版本 |
|------|------|------|
| 框架 | React | 18.2.0 |
| 语言 | TypeScript | 5.3.0 |
| UI组件库 | Ant Design | 5.12.0 |
| 图表库 | ECharts | 6.1.0 |
| 网络请求 | Axios | 1.6.0 |
| 路由 | React Router DOM | 6.20.0 |
| 构建工具 | Vite | 5.0.0 |
| 样式 | SCSS | 1.101.0 |
| Excel导出 | XLSX + FileSaver | 0.18.5 / 2.0.5 |

**开发端口**: `http://localhost:4000`

### 2.2 用户端前端 (frontend)

| 分类 | 技术 | 版本 |
|------|------|------|
| 框架 | Vue 3 | 3.4.21 |
| 跨端框架 | Uni-app | 3.0.0-5000720260410001 |
| UI组件库 | uv-ui (@climblee/uv-ui) | 1.1.20 |
| 国际化 | vue-i18n | 9.14.4 |
| 构建工具 | Vite | 5.2.8 |

**开发端口**: `http://localhost:5173`

### 2.3 后端服务 (backend)

| 分类 | 技术 | 版本 |
|------|------|------|
| 语言 | Go | 1.25.0 |
| Web框架 | Gin | 1.10.0 |
| ORM | GORM | 1.25.7 |
| 数据库 | MySQL | 8.0+ |
| 缓存 | Redis | 7+ |
| JWT | golang-jwt/jwt/v4 | 4.5.0 |
| 日志 | Zap | 1.28.0 |
| 配置 | Viper | 1.18.2 |
| 支付SDK | huifurepo/bspay-go-sdk | 1.0.28 |

**开发端口**: `http://localhost:3300`

### 2.4 系统架构图

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   PC管理后台    │    │   用户端(H5/小程序) │    │   硬件设备      │
│  (React + TS)   │◄──►│  (Vue3 + UniApp) │◄──►│ (蓝牙/WiFi)     │
│  http://4000    │    │  http://5173     │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                      │                      │
         │                      │                      │
         ▼                      ▼                      ▼
┌─────────────────────────────────────────────────────────────┐
│                 Centraliz Backend (Go + Gin)               │
│                     http://localhost:3300                   │
└─────────────────────────────────────────────────────────────┘
         │                      │                      │
         ▼                      ▼                      ▼
┌───────────────┐    ┌─────────────────┐    ┌─────────────────┐
│    MySQL      │    │     Redis       │    │   Huifu Pay     │
│  localhost:3306│    │  localhost:6379 │    │  汇付天下支付   │
└───────────────┘    └─────────────────┘    └─────────────────┘
```

## 3. 功能模块设计

### 3.1 用户认证模块

#### 3.1.1 注册方式
- **手机号验证码注册**：支持国际手机号格式，验证码有效期5分钟
- **邮箱验证码注册**：支持主流邮箱服务商，验证码有效期30分钟
- **账号密码注册**：密码强度要求8-20位

#### 3.1.2 登录方式
- 手机号+验证码登录
- 邮箱+验证码登录
- 账号+密码登录
- 微信授权登录（小程序/H5）

#### 3.1.3 账号安全
- JWT Token认证
- Refresh Token机制
- 密码错误次数限制
- 敏感操作二次验证

### 3.2 设备管理模块

#### 3.2.1 中控存柜功能
- **柜子状态监控**：在线/离线状态、电池电量、门锁状态
- **存取操作**：扫码开柜、预约存取、临时授权
- **异常处理**：超时未取提醒、强制开柜、故障上报

#### 3.2.2 中控无人零售柜功能
- **商品管理**：商品上下架、库存监控、价格设置
- **购买流程**：扫码选购、支付集成、取货验证
- **补货管理**：补货提醒、补货记录、销售统计

### 3.3 管理端功能

#### 3.3.1 PC管理后台特性
- **房间管理**：房间列表、创建/编辑/删除、标签管理、图片管理
- **设备管理**：设备列表、创建设备、绑定分组、设备日志
- **分组管理**：分组列表、创建/编辑/删除、规则绑定
- **规则管理**：计费规则、使用规则、报警规则
- **订单管理**：订单列表、退款审批、押金管理
- **商家管理**：商家列表、子商户管理、汇付账号管理
- **用户管理**：用户列表、微信用户管理

#### 3.3.2 H5管理端特性
- **响应式设计**：适配手机/平板/PC
- **公众号集成**：JS-SDK调用、微信分享、微信支付
- **蓝牙通信**：蓝牙设备扫描、数据传输、WiFi配置推送

#### 3.3.3 智能开发助手 (Agent)
- **PC端**：仿 ChatGPT 风格聊天界面，支持 Markdown、代码高亮、文件 diff 预览
- **移动端**：轻量级助手入口，支持快捷指令
- **核心能力**：代码编写、接口调试、图表生成、Bug修复
- **安全机制**：高风险操作确认弹窗、影响文件预览、人工确认机制

### 3.4 国际化支持
- **支持语言**：中文（简体/繁体）、英语、日语、韩语、法语、德语、西班牙语、俄语、泰语、土耳其语、阿拉伯语、葡萄牙语、意大利语
- **本地化特性**：日期时间格式、数字货币格式、RTL布局支持（阿拉伯语）

## 4. 开发流程

### 4.1 环境搭建

#### 4.1.1 后端环境

```bash
# Go环境要求
Go 1.20+

# 进入后端目录
cd d:\c-full-stack\centraliz-agent\backend

# 安装依赖
go mod tidy

# 启动后端服务（热加载模式）
go run main.go
```

**配置文件**: `backend/config/config.yaml`

| 配置项 | 开发环境 | 说明 |
|--------|----------|------|
| server.port | :3300 | 服务端口 |
| database.host | localhost | 数据库地址 |
| database.port | 3306 | 数据库端口 |
| database.name | centraliz | 数据库名 |
| database.username | root | 数据库用户名 |
| database.password | 123456 | 数据库密码 |
| redis.host | localhost | Redis地址 |
| redis.port | 6379 | Redis端口 |
| jwt.secret | centraliz-production-jwt-secret-key-2026 | JWT密钥 |

#### 4.1.2 PC管理后台环境

```bash
# 进入admin目录
cd d:\c-full-stack\centraliz-agent\admin

# 安装依赖
npm install

# 启动开发服务器
npm run dev

# 构建生产版本
npm run build

# 类型检查
npm run typecheck
```

**开发端口**: `http://localhost:4000`

**API代理**: `/api` -> `http://localhost:3300`

#### 4.1.3 用户端前端环境

```bash
# 进入frontend目录
cd d:\c-full-stack\centraliz-agent\frontend

# 安装依赖
npm install

# 启动H5开发服务器
npm run dev:h5

# 启动微信小程序开发
npm run dev:mp-weixin

# 构建H5版本
npm run build:h5

# 构建微信小程序版本
npm run build:mp-weixin
```

**开发端口**: `http://localhost:5173`

**API代理**: `/api/v1` -> `http://localhost:3300`

### 4.2 后端路由设计

后端API采用RESTful风格，路由前缀为 `/api/v1`，分为公共路由和需要认证的路由。

#### 4.2.1 公共路由（无需认证）

| 路由 | 方法 | 功能 |
|------|------|------|
| `/api/v1/device/add` | POST | 添加集控设备 |
| `/api/v1/device/common` | POST | 集控设备通用接口 |
| `/api/v1/devicelog/cablelogs` | POST | 创建设备日志 |
| `/api/v1/wechat/login` | POST | 微信登录 |
| `/api/v1/wechat/callback` | GET | 微信登录回调 |
| `/api/v1/wechat/jssdk/config` | GET | 获取JS SDK配置 |
| `/api/v1/wechat/scan/image` | POST | 扫描二维码 |
| `/api/v1/send-code` | POST | 发送验证码 |
| `/api/v1/user/login` | POST | 用户登录 |
| `/api/v1/user/register` | POST | 用户注册 |
| `/api/v1/merch/login` | POST | 商家登录 |
| `/api/v1/merch/register` | POST | 商家注册 |
| `/api/v1/merch/reset-password` | POST | 商家重置密码 |
| `/api/v1/huifu/qrpay/h5/wechat` | POST | H5微信支付 |
| `/api/v1/huifu/qrpay/mini/wechat` | POST | 小程序支付 |
| `/api/v1/huifu/qrpay/jsapi/wechat` | POST | JS API支付 |

#### 4.2.2 需要认证的路由（JWT）

| 路由分组 | 功能 |
|----------|------|
| `/api/v1/user/*` | 用户个人信息、房间、订单、押金相关 |
| `/api/v1/merch/*` | 商家个人信息、密码修改、绑定邮箱/手机 |
| `/api/v1/device/*` | 设备管理、设备日志 |
| `/api/v1/group/*` | 分组管理 |
| `/api/v1/room/*` | 房间管理、房间标签、房间图片 |
| `/api/v1/order/*` | 订单管理、退款审批、押金管理 |
| `/api/v1/huifu/*` | 汇付账号管理 |
| `/api/v1/submerch/*` | 子账号管理 |
| `/api/v1/merch-pay/*` | 商家消费订单管理 |
| `/api/v1/rule/*` | 规则管理 |
| `/api/v1/admin/*` | PC管理后台专属接口 |

#### 4.2.3 PC管理后台路由

| 路由分组 | 功能 |
|----------|------|
| `/api/v1/admin/stats` | 首页统计数据 |
| `/api/v1/admin/room/*` | 房间管理（增删改查、批量操作、导入） |
| `/api/v1/admin/device/*` | 设备管理（增删改查、批量操作、导入） |
| `/api/v1/admin/group/*` | 分组管理（增删改查、批量操作、导入） |
| `/api/v1/admin/rule/*` | 规则管理（增删改查、批量操作、导入） |
| `/api/v1/admin/order/*` | 订单管理（增删改查、批量操作） |
| `/api/v1/admin/merch/*` | 商家管理（增删改查、批量操作、导入） |
| `/api/v1/admin/devicelog/*` | 设备日志管理（查询、删除） |
| `/api/v1/admin/huifu/*` | 汇付账号管理（增删改查、批量操作、导入） |
| `/api/v1/admin/merchpay/*` | 商户支付管理（查询、删除） |
| `/api/v1/admin/roomimg/*` | 房间图片管理（增删改查、批量操作、导入） |
| `/api/v1/admin/roomtag/*` | 房间标签管理（增删改查、批量操作、导入） |
| `/api/v1/admin/submerch/*` | 子商户管理（增删改查、批量操作、导入） |
| `/api/v1/admin/user/*` | 用户管理（增删改查、批量操作） |
| `/api/v1/admin/wxuser/*` | 微信用户管理（增删改查、批量操作） |
| `/api/v1/admin/assistant/*` | 智能开发助手（聊天、会话管理、操作确认） |

## 5. 数据库模型

### 5.1 表结构概览

系统包含以下核心数据表：

| 表名 | 说明 |
|------|------|
| `codes` | 验证码管理表 |
| `devices` | 设备信息表 |
| `group` | 分组管理表 |
| `huifu_account` | 汇付账户表 |
| `merch` | 商家主账户表 |
| `merch_pay` | 商家支付订单表 |
| `orders` | 用户订单表 |
| `room` | 房间信息表 |
| `submerch` | 商家子账户表 |
| `user` | 用户信息表 |
| `wechat_user` | 微信用户表 |

### 5.2 数据类型映射规范

| MySQL类型 | Go类型 | 说明 |
|-----------|--------|------|
| bigint(20) unsigned | uint64 | 主键/大整数 |
| int(11)/int(20) | int32 | 普通整数 |
| varchar(N) | string | 字符串（带gorm size标签） |
| text | string | 长文本 |
| datetime/datetime(3) | time.Time | 时间类型 |
| decimal(M,N) | float64 | 小数（带gorm type:decimal标签） |
| bigint | BigIntTime | 时间戳 |

## 6. 部署方案

### 6.1 前端部署

| 项目 | 部署方式 | 说明 |
|------|----------|------|
| admin | Nginx静态服务器 | 构建产物dist目录 |
| frontend(H5) | Nginx静态服务器 | 构建产物dist/build/h5目录 |
| frontend(小程序) | 微信公众平台 | 上传代码包 |

### 6.2 后端部署

```bash
# 构建后端二进制
cd d:\c-full-stack\centraliz-agent\backend
go build -o centraliz.exe main.go

# 运行
./centraliz.exe
```

**Windows启动脚本**: `start.bat` / `start_daemon.bat`

### 6.3 Docker部署

```yaml
# docker-compose.yml
version: '3.8'
services:
  backend:
    build: ./backend
    ports:
      - "3300:3300"
    environment:
      - DB_HOST=mysql
      - REDIS_HOST=redis
    depends_on:
      - mysql
      - redis

  mysql:
    image: mysql:8.0
    environment:
      - MYSQL_ROOT_PASSWORD=123456
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

## 7. 安全考虑

### 7.1 数据安全
- HTTPS全站加密
- 敏感数据加密存储
- API请求签名验证
- JWT Token认证

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
- Ant Design按需加载

### 8.2 后端优化
- 数据库索引优化
- Redis缓存热点数据
- 异步处理耗时操作
- CDN静态资源加速

## 9. 代码规范

### 9.1 前端
- ESLint + Prettier（frontend）
- TypeScript严格模式（admin）
- SCSS模块化

### 9.2 后端
- gofmt + golangci-lint
- GORM最佳实践

### 9.3 Git提交规范
- Conventional Commits
- 提交信息格式：`feat: xxx` / `fix: xxx` / `docs: xxx`

## 10. 第三方服务集成

| 服务类型 | 服务提供商 |
|----------|------------|
| 短信服务 | 阿里云短信 |
| 邮件服务 | QQ邮箱SMTP |
| 支付服务 | 汇付天下（微信支付） |
| 文件存储 | 本地存储/阿里云OSS |
| 微信平台 | 微信小程序/公众号 |
