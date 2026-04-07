# TaoShop

一个适合练手的 Go 全栈电商项目，主题接近淘宝首页风格，包含注册、登录、商品浏览、购买和订单记录。

## 技术栈

- Go 1.26
- Gin
- MySQL
- Redis
- 原生 HTML / CSS / JavaScript 前端
- JWT 鉴权

## 功能点

- 用户注册 / 登录
- 商品列表查询
- 下单购买
- 我的订单列表
- Redis 商品缓存
- Redis 基于 IP 的简单限流
- MySQL 事务扣减库存
- 单体分层结构，可继续拆成微服务

## 目录结构

```text
cmd/server            程序入口
internal/app          启动与路由
internal/http         handler 与 middleware
internal/service      业务逻辑
internal/repository   数据访问
internal/database     MySQL 连接与建表初始化
internal/cache        Redis 连接
web                   前端页面
```

## 本地启动

1. 准备 MySQL 和 Redis。
2. 按 `.env.example` 设置环境变量。
3. 启动服务：

```powershell
$env:MYSQL_DSN="root:root@tcp(127.0.0.1:3306)/taoshop?parseTime=true&charset=utf8mb4&loc=Local"
$env:REDIS_ADDR="127.0.0.1:6379"
$env:JWT_SECRET="your-secret"
go run ./cmd/server
```

4. 打开 [http://localhost:8080](http://localhost:8080)。

## Docker

仓库提供了 `docker-compose.yml` 用来快速启动 MySQL 和 Redis：

```powershell
docker compose up -d
```

如果本机没有安装 Docker，也可以直接使用本地安装好的 MySQL / Redis。

## 后续扩展建议

- 把商品、用户、订单拆成独立微服务
- 增加消息队列处理异步订单
- 引入 Elasticsearch 做商品搜索
- 增加 Nginx 网关和多实例部署
- 增加 Prometheus / Grafana 监控
- 增加分布式 ID、分布式锁和库存预扣减
