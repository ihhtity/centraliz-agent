# Centraliz 后端服务

基于 Gin 框架开发的 RESTful API 服务，为前端提供设备管理、用户管理、订单管理等功能。

## 技术栈

- **框架**: Gin
- **数据库**: GORM + MySQL (数据库名: centraliz, 用户名: root, 密码: 123456)
- **缓存**: Redis
- **认证**: JWT
- **语言**: Go 1.21+

## 项目结构

```
backend/
├── config/          # 配置文件
├── controller/      # 控制器层
├── middleware/      # 中间件
├── model/           # 数据模型
├── pkg/             # 工具包
│   ├── db/          # 数据库相关
│   ├── jwt/         # JWT认证相关
│   ├── redis/       # Redis缓存相关
│   └── response/    # 响应处理
├── router/          # 路由配置
├── go.mod           # Go模块配置
└── main.go          # 主入口文件
```

## 架构改进方案

### 1. 配置管理（Viper + 结构体映射，环境变量覆盖）

当前实现已使用Viper进行配置管理，支持`.env`文件和环境变量。建议改进：
- 添加配置验证和类型安全
- 支持多环境配置文件（dev/test/prod）
- 实现配置热更新机制

### 2. 环境区分配置（开发/测试/生产）

当前仅支持单一`.env`配置。建议创建多环境配置：
- `config/dev.env` - 开发环境
- `config/test.env` - 测试环境  
- `config/prod.env` - 生产环境
- 通过`APP_ENV`环境变量切换

### 3. 依赖注入与组件初始化

当前采用手动初始化方式。建议引入Wire或重构为更清晰的依赖注入模式：
- 创建`di`包管理依赖注入
- 实现组件工厂模式
- 支持优雅启停

### 4. 中间件配置增强

当前仅有JWT认证中间件。建议添加：
- CORS跨域中间件
- Recovery错误恢复中间件
- Zap日志中间件
- 限流中间件（如令牌桶）
- 请求ID追踪中间件

### 5. 数据库与缓存配置优化

当前数据库连接较为简单。建议改进：
- 配置连接池参数（MaxOpenConns, MaxIdleConns, ConnMaxLifetime）
- 实现读写分离（主从复制）
- 添加数据库迁移工具（如goose或gormigrate）
- Redis连接池配置优化

### 6. 日志配置

当前无专门日志配置。建议集成Zap日志库：
- 支持JSON/Console格式输出
- 按级别和时间分割日志文件
- 实现采样日志避免性能问题
- 集成到中间件中记录请求日志

### 7. 路由设计优化

当前路由设计良好，已有版本控制和分组。建议：
- 添加Swagger文档自动生成
- 实现更细粒度的中间件绑定
- 添加路由指标监控

### 8. 错误处理与响应统一

当前已有基础响应封装。建议：
- 定义标准错误码体系
- 实现全局错误处理中间件
- 统一异常处理流程

### 9. 优雅启停与健康检查

当前无优雅启停机制。建议添加：
- HTTP服务器优雅关闭
- 信号处理（SIGTERM, SIGINT）
- 健康检查端点（/health, /ready）
- 就绪探针和存活探针

### 10. 配置热更新

当前配置在启动时加载。建议：
- 监听配置文件变化
- 动态更新配置（如日志级别、限流参数）
- 线程安全的配置访问

### 11. 测试与文档配置

建议添加：
- 单元测试和集成测试框架
- Swagger API文档自动生成
- Postman集合导出
- 性能基准测试

## 已实现的改进

✅ **配置管理**：已实现多环境YAML配置（dev/test/prod），支持环境变量覆盖
✅ **日志系统**：已集成Zap日志库，支持JSON/Console格式，文件分割
✅ **中间件增强**：已添加CORS、Recovery、Logger、RateLimit中间件
✅ **优雅启停**：已实现HTTP服务器优雅关闭和信号处理
✅ **健康检查**：已添加/health和/ready端点
✅ **数据库优化**：已配置连接池参数和Redis连接池
✅ **限流机制**：已实现基于令牌桶的限流中间件

## 使用说明

### 环境配置

1. **开发环境**（默认）：
   ```bash
   APP_ENV=dev go run main.go
   ```

2. **生产环境**：
   ```bash
   APP_ENV=prod go run main.go
   ```

3. **测试环境**：
   ```bash
   APP_ENV=test go run main.go
   ```

### 配置文件优先级

1. 环境变量（最高优先级）
2. YAML配置文件（config/{APP_ENV}.yaml）
3. 默认值（最低优先级）

### 健康检查端点

- **存活检查**：`GET /health`
- **就绪检查**：`GET /ready`

### 日志配置

支持以下日志级别：
- debug
- info  
- warn
- error

支持输出格式：
- console（带颜色）
- json（结构化日志）

### 限流配置

在配置文件中设置：
```yaml
rate_limit:
  requests_per_second: 100
  burst: 200
```

限流仅在非调试模式下启用。

## 快速开始

### 1. 安装依赖

确保已安装 Go 1.21+，然后运行：

```bash
go mod tidy
```

### 2. 数据库准备

1. 确保 MySQL 服务已启动
2. 创建数据库 `centraliz`：
   ```sql
   CREATE DATABASE centraliz CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
   ```
3. 确保 Redis 服务已启动

### 3. 启动服务

```bash
go run main.go
```

服务将在 `http://localhost:8080` 启动。

### 4. 测试API

使用 curl 测试登录接口：

```bash
curl -X POST http://localhost:8080/api/v1/user/login
```

## 配置说明

- **环境变量**: 复制 `.env.example` 为 `.env` 并根据需要修改配置
- **JWT密钥**: 在 `config/config.go` 中配置，生产环境务必修改
- **服务器端口**: 默认 8080，在 `config/config.go` 中可修改

## 注意事项

- 当前实现为示例代码，实际使用时需要：
  - 实现完整的数据库操作逻辑
  - 添加输入验证和错误处理
  - 实现真实的用户认证逻辑
  - 添加日志记录
  - 配置生产环境的安全设置
- **重要**: 生产环境中必须修改默认的 JWT 密钥和数据库密码