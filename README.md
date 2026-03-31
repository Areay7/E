# 跨境电商后台管理系统

商用级跨境电商后台管理系统，支持 Shopee、速卖通、TikTok 等多平台对接。

## 技术栈

### 后端
- **Go 1.22** - 高性能后端语言
- **Gin** - Web 框架
- **GORM** - ORM 框架
- **MySQL/PostgreSQL** - 生产级数据库
- **Viper** - 配置管理

### 前端
- **Vue 3** - 渐进式框架
- **TypeScript** - 类型安全
- **Element Plus** - UI 组件库
- **Vue Router** - 路由管理
- **Pinia** - 状态管理
- **Axios** - HTTP 客户端
- **Vite** - 构建工具

## 核心功能

### 1. 多平台管理
- **Shopee** - 东南亚电商平台
- **速卖通 (AliExpress)** - 全球速卖通
- **TikTok Shop** - 短视频电商

### 2. 订单管理
- 订单列表查询（支持多条件筛选）
- 订单详情查看
- 订单状态同步
- 批量发货处理
- 物流跟踪

### 3. 商品管理
- 商品列表管理
- 商品详情查看
- 商品变体管理
- 商品状态更新
- 跨平台商品同步

### 4. 库存管理
- 实时库存查询
- 库存更新
- 库存变动日志
- 多仓库支持

### 5. 平台 API 管理
- 统一的平台接口封装
- API 调用日志记录
- 签名认证处理
- 错误重试机制

### 6. 数据同步
- 定时同步订单
- 定时同步商品
- 定时同步库存
- 同步任务监控

## 项目结构

```
E/
├── backend/                    # Go 后端
│   ├── config/                 # 配置管理
│   │   ├── config.go          # 配置结构定义
│   │   └── config.yaml        # 配置文件
│   ├── models/                 # 数据模型
│   │   ├── db.go              # 数据库连接
│   │   └── models.go          # 业务模型
│   ├── platform/               # 平台 API 封装
│   │   ├── interface.go       # 平台接口定义
│   │   ├── registry.go        # 平台注册管理
│   │   ├── shopee/            # Shopee 实现
│   │   ├── aliexpress/        # 速卖通实现
│   │   └── tiktok/            # TikTok 实现
│   ├── service/                # 业务逻辑层
│   │   ├── order_service.go   # 订单服务
│   │   ├── product_service.go # 商品服务
│   │   └── inventory_service.go # 库存服务
│   ├── handlers/               # HTTP 处理器
│   │   ├── dashboard.go       # 数据看板
│   │   ├── platform.go        # 平台管理
│   │   └── business.go        # 业务接口
│   ├── routes/                 # 路由配置
│   │   └── routes.go
│   ├── main.go                 # 入口文件
│   └── go.mod
│
└── frontend/                   # Vue3 前端
    ├── src/
    │   ├── api/                # API 接口
    │   ├── components/         # 组件
    │   ├── layout/             # 布局
    │   ├── router/             # 路由
    │   ├── types/              # 类型定义
    │   ├── utils/              # 工具函数
    │   └── views/              # 页面
    └── package.json
```

## 数据库设计

### 核心表
- **orders** - 订单主表
- **order_items** - 订单明细
- **products** - 商品表
- **product_variations** - 商品变体
- **inventory** - 库存表
- **inventory_logs** - 库存变动日志
- **logistics** - 物流信息
- **platform_configs** - 平台配置
- **api_logs** - API 调用日志
- **sync_tasks** - 同步任务记录

## 快速开始

### 1. 数据库准备

**MySQL:**
```sql
CREATE DATABASE cross_border_ecommerce CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

**PostgreSQL:**
```sql
CREATE DATABASE cross_border_ecommerce ENCODING 'UTF8';
```

### 2. 配置文件

编辑 `backend/config/config.yaml`：

```yaml
database:
  type: mysql  # 或 postgres
  host: localhost
  port: 3306
  username: your_username
  password: your_password
  dbname: cross_border_ecommerce
```

### 3. 启动后端

```bash
cd backend

# 安装依赖
go mod tidy

# 启动服务
go run main.go
```

后端服务将在 `http://localhost:8080` 启动

### 4. 启动前端

```bash
cd frontend

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

前端服务将在 `http://localhost:3000` 启动

## API 接口

### 平台管理
- `GET /api/v1/platforms` - 获取所有平台配置
- `GET /api/v1/platforms/:platform/config` - 获取平台配置
- `PUT /api/v1/platforms/:platform/config` - 更新平台配置
- `POST /api/v1/platforms/:platform/sync/orders` - 同步订单
- `POST /api/v1/platforms/:platform/sync/products` - 同步商品

### 订单管理
- `GET /api/v1/orders` - 获取订单列表
- `GET /api/v1/orders/:id` - 获取订单详情
- `POST /api/v1/orders/:id/ship` - 订单发货

### 商品管理
- `GET /api/v1/products` - 获取商品列表
- `GET /api/v1/products/:id` - 获取商品详情

### 库存管理
- `GET /api/v1/inventory` - 获取库存列表
- `PUT /api/v1/inventory` - 更新库存

### 监控
- `GET /api/v1/api-logs` - API 调用日志
- `GET /api/v1/sync-tasks` - 同步任务列表
- `GET /api/v1/dashboard/summary` - 数据看板

## 平台对接说明

### Shopee
1. 在 Shopee Open Platform 注册应用
2. 获取 Partner ID 和 Partner Key
3. 配置到 `config.yaml` 的 `platforms.shopee` 部分

### 速卖通
1. 在阿里百川开放平台注册应用
2. 获取 App Key 和 App Secret
3. 配置到 `config.yaml` 的 `platforms.aliexpress` 部分

### TikTok Shop
1. 在 TikTok Shop Seller Center 注册应用
2. 获取 App Key 和 App Secret
3. 配置到 `config.yaml` 的 `platforms.tiktok` 部分

## 开发说明

### 添加新平台

1. 在 `platform/` 下创建新平台目录
2. 实现 `Platform` 接口
3. 在 `init()` 中注册平台
4. 更新配置文件

### 扩展功能

框架已预留扩展接口：
- 认证接口 (AuthAPI)
- 订单接口 (OrderAPI)
- 商品接口 (ProductAPI)
- 库存接口 (InventoryAPI)
- 物流接口 (LogisticsAPI)

## 生产部署

### 后端部署
```bash
# 编译
go build -o cross-border-admin main.go

# 运行
./cross-border-admin
```

### 前端部署
```bash
# 构建
npm run build

# 部署 dist 目录到 Nginx/CDN
```

## 注意事项

1. **配置文件安全**: 生产环境请妥善保管 `config.yaml`，不要提交到版本控制
2. **数据库连接**: 建议使用连接池，已在代码中配置
3. **API 限流**: 各平台 API 有调用频率限制，请注意控制
4. **日志管理**: API 调用日志会持续增长，建议定期清理
5. **错误处理**: 平台 API 调用失败会记录日志，需要监控和处理

## License

MIT
