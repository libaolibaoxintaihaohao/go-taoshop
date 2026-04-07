# TaoShop 学习指南

这份文档不是 API 手册，而是按“第一次完整做 Go 项目”的视角，带你理解这个项目为什么这样写、应该按什么顺序写、每一层负责什么，以及你现在最需要掌握的 Go / Gin / MySQL / Redis / 微服务基础。

## 1. 先建立总图

这个项目可以先把它看成 5 层：

1. `web/`
   前端页面，负责展示商品、登录注册、下单和查看订单。
2. `internal/http/handler`
   HTTP 接口入口。负责接收请求、校验参数、返回 JSON。
3. `internal/service`
   业务逻辑层。负责“真正要做什么”。
4. `internal/repository`
   数据访问层。负责和 MySQL 交互。
5. `internal/database` / `internal/cache`
   基础设施层。负责连接 MySQL、Redis。

你可以把一次“下单”理解成这样一条链路：

前端点击购买 -> Gin 路由收到请求 -> Handler 解析参数 -> Service 开事务 -> Repository 查商品和扣库存 -> Repository 写订单 -> Service 清 Redis 缓存 -> Handler 返回结果给前端

这就是一个典型的后端项目执行流程。

## 2. 为什么我会按这个顺序写代码

如果你第一次写完整项目，推荐按下面顺序，而不是想到哪写到哪。

### 第一步：先定需求和数据模型

你一开始要先想清楚：

- 这个网站有什么页面和功能
- 后端需要哪些接口
- 数据库里要有哪些表

这个项目一开始的核心需求只有 4 个：

- 注册
- 登录
- 商品列表
- 下单和查看订单

所以数据库最少要有 3 张表：

- `users`
- `products`
- `orders`

这一步对应代码主要是：

- [models.go](C:\Users\HP\Desktop\go_test\internal\models\models.go)
- [schema.go](C:\Users\HP\Desktop\go_test\internal\database\schema.go)

为什么先做这一步：

因为数据库表和业务对象定不下来，后面的接口、服务、前端都容易反复推翻。

### 第二步：先把基础设施接起来

任何真实项目都要先能“活着启动”：

- 读取配置
- 连接 MySQL
- 连接 Redis

这一步对应：

- [config.go](C:\Users\HP\Desktop\go_test\internal\config\config.go)
- [mysql.go](C:\Users\HP\Desktop\go_test\internal\database\mysql.go)
- [redis.go](C:\Users\HP\Desktop\go_test\internal\cache\redis.go)

为什么第二步做这个：

因为如果数据库和缓存都连不上，业务代码写再多也跑不起来。

### 第三步：先写 repository，而不是先写页面

很多新手会先写页面，但后端项目更稳的做法是先把数据访问写出来。

原因很简单：

- `repository` 决定你怎么查表
- `service` 依赖 `repository`
- `handler` 依赖 `service`
- 前端又依赖 `handler` 提供的 API

所以应该先把底层铺好。

对应文件：

- [user_repository.go](C:\Users\HP\Desktop\go_test\internal\repository\user_repository.go)
- [product_repository.go](C:\Users\HP\Desktop\go_test\internal\repository\product_repository.go)
- [order_repository.go](C:\Users\HP\Desktop\go_test\internal\repository\order_repository.go)

### 第四步：写 service，表达业务规则

`service` 是整个项目最重要的层，因为它承载“为什么这样做”。

比如：

- 注册时要加密密码
- 登录时要校验密码并签发 JWT
- 商品列表优先查 Redis
- 下单时必须开事务、锁库存、扣库存、写订单

对应文件：

- [auth_service.go](C:\Users\HP\Desktop\go_test\internal\service\auth_service.go)
- [catalog_service.go](C:\Users\HP\Desktop\go_test\internal\service\catalog_service.go)
- [order_service.go](C:\Users\HP\Desktop\go_test\internal\service\order_service.go)

为什么 service 单独拆出来：

因为“业务规则”不应该塞进 HTTP 控制器里，也不应该塞进 SQL 访问层里。

### 第五步：写 handler 和 middleware

到这一步，后端已经基本可用了。然后再写：

- 参数绑定
- 返回 JSON
- 鉴权
- 限流

对应文件：

- [auth_handler.go](C:\Users\HP\Desktop\go_test\internal\http\handler\auth_handler.go)
- [catalog_handler.go](C:\Users\HP\Desktop\go_test\internal\http\handler\catalog_handler.go)
- [order_handler.go](C:\Users\HP\Desktop\go_test\internal\http\handler\order_handler.go)
- [auth.go](C:\Users\HP\Desktop\go_test\internal\http\middleware\auth.go)
- [rate_limit.go](C:\Users\HP\Desktop\go_test\internal\http\middleware\rate_limit.go)

为什么这一步放后面：

因为 handler 只是“入口壳子”，它本身不应该承载核心业务。

### 第六步：最后再接前端页面

当前端接入时，你已经有了稳定的 API：

- `POST /api/auth/register`
- `POST /api/auth/login`
- `GET /api/products`
- `POST /api/orders`
- `GET /api/orders/me`

这时再写前端，效率最高。

对应文件：

- [index.html](C:\Users\HP\Desktop\go_test\web\index.html)
- [app.js](C:\Users\HP\Desktop\go_test\web\assets\app.js)
- [styles.css](C:\Users\HP\Desktop\go_test\web\assets\styles.css)

## 3. 为什么项目要分层

如果你以前主要看 C++ 代码，第一次看这种 Web 项目，最容易困惑的是“为什么不直接全写在 main 里”。

答案是：因为 Web 项目复杂度很快就会涨。

分层的目的是让不同代码只做自己的事：

- `main`
  负责启动程序
- `app`
  负责把所有东西组装起来
- `handler`
  负责 HTTP 输入输出
- `service`
  负责业务逻辑
- `repository`
  负责数据库读写
- `database/cache`
  负责连接基础设施
- `web`
  负责界面

这样分开之后，你读代码就知道应该从哪里找问题。

比如：

- 接口 404，看路由
- 参数错误，看 handler
- 库存扣错，看 service
- SQL 有问题，看 repository
- 数据库连不上，看 database

## 4. Go 语法先抓最常用的，不要一口气学完

你现在做这个项目，最需要理解的是以下几个语法点。

### 包和导入

```go
package service

import (
    "context"
    "database/sql"
)
```

含义：

- `package service` 表示这个文件属于 `service` 包
- `import` 表示导入别的包

Go 里按目录组织包，目录感非常强。

### 结构体 struct

```go
type User struct {
    ID       int64
    Username string
    Email    string
}
```

这相当于 C++ 里的数据结构，用来描述对象。

在这个项目里：

- `User` 表示用户
- `Product` 表示商品
- `Order` 表示订单

### 方法 method

```go
func (s *AuthService) Login(ctx context.Context, email, password string) (string, *models.User, error)
```

这行拆开理解：

- `func` 表示定义函数
- `(s *AuthService)` 表示这是挂在 `AuthService` 上的方法
- `ctx context.Context, email, password string` 是参数
- `(string, *models.User, error)` 是返回值

Go 支持多返回值，这点和很多语言很不一样。

这个方法会同时返回：

- token
- user
- error

### 指针

```go
func NewAuthService(users *repository.UserRepository, jwtSecret string) *AuthService
```

`*AuthService` 表示返回一个指针。

第一次你可以简单理解为：

- 用指针可以避免复制大对象
- 方法里能直接操作同一个服务实例

你现在不用把指针学得特别深，先知道服务对象一般都用指针传递就够了。

### if err != nil

Go 项目里你会反复看到：

```go
if err != nil {
    return nil, err
}
```

这是 Go 最经典的错误处理方式。

意思是：

- 如果出错
- 就立刻返回

这种写法虽然看起来朴素，但非常清晰。

### make

```go
orders := make([]models.Order, 0)
```

意思是创建一个空切片。

你可以先把切片理解成“动态数组”。

### JSON tag

```go
Username string `json:"username"`
```

这个反引号里的内容叫 tag。

作用是告诉 Go：

- 这个字段序列化成 JSON 时用 `username`

所以前端拿到的数据字段才会是小写下划线风格。

### 结构体字面量

```go
user := &models.User{
    Username: username,
    Email: email,
}
```

表示创建一个 `User` 对象，并给字段赋值。

### defer

```go
defer tx.Rollback()
```

意思是：函数结束前执行 `Rollback()`。

在事务里常见写法是：

- 先 `BeginTx`
- 马上 `defer Rollback`
- 如果最后成功，再 `Commit`

这样能避免中途报错却忘记回滚。

## 5. Gin 到底是什么

Gin 可以先简单理解为：

“Go 语言里的一个轻量 Web 框架，用来快速写 HTTP 接口。”

你当前主要用到它的 4 个能力：

### 路由

```go
api.POST("/auth/login", authHandler.Login)
```

表示：

- 当收到 `POST /api/auth/login`
- 就执行 `authHandler.Login`

### 上下文 `*gin.Context`

```go
func (h *AuthHandler) Login(c *gin.Context)
```

`c` 里包含：

- 请求参数
- 请求头
- 响应写出能力

你可以从里面取 JSON、取 header、返回 JSON。

### 参数绑定

```go
if err := c.ShouldBindJSON(&req); err != nil {
```

意思是把请求体 JSON 自动解析到 `req` 结构体中。

### 返回 JSON

```go
c.JSON(http.StatusOK, gin.H{"message": "login success"})
```

`gin.H` 本质上是一个 map，用来快速返回 JSON。

## 6. Redis 在这里为什么存在

如果没有 Redis，项目也能跑。

但加 Redis 是为了让你提前接触真实项目里非常常见的两类用途。

### 用途 1：缓存商品列表

对应：

- [catalog_service.go](C:\Users\HP\Desktop\go_test\internal\service\catalog_service.go)

逻辑是：

1. 先查 Redis
2. 如果命中，直接返回
3. 如果没命中，查 MySQL
4. 再把结果写回 Redis

这叫缓存旁路模式，实际项目里非常常见。

### 用途 2：限流

对应：

- [rate_limit.go](C:\Users\HP\Desktop\go_test\internal\http\middleware\rate_limit.go)

逻辑是：

1. 用 IP 作为 key
2. 访问一次计数加一
3. 超过阈值就拒绝

这是最基础的限流思路。

## 7. MySQL 在这里为什么比文件或内存更合适

因为电商数据是强业务数据：

- 用户要持久化
- 订单要持久化
- 库存要准确

这类数据适合放关系型数据库。

尤其是下单时，我们希望：

- 扣库存
- 写订单

这两个动作要么都成功，要么都失败。

所以我们用了事务。

对应：

- [order_service.go](C:\Users\HP\Desktop\go_test\internal\service\order_service.go)

## 8. JWT 是什么，为什么登录后不用 session

JWT 可以先理解成：

“服务器签发的一个带签名的用户身份令牌”

登录成功后，后端会返回 token，前端把它存在 `localStorage`，之后请求时放进：

```text
Authorization: Bearer <token>
```

后端再在中间件里解析它。

对应：

- [auth_service.go](C:\Users\HP\Desktop\go_test\internal\service\auth_service.go)
- [auth.go](C:\Users\HP\Desktop\go_test\internal\http\middleware\auth.go)

## 9. 什么是分布式，什么是微服务

你现在先不要把这两个词想得太神秘。

### 单体项目

现在这个项目本质上是单体：

- 一个 Go 服务
- 一个进程
- 连接 MySQL 和 Redis

优点：

- 容易启动
- 容易调试
- 适合练手

### 微服务

微服务是把一个大系统拆成多个服务，比如：

- 用户服务
- 商品服务
- 订单服务
- 支付服务

每个服务单独部署、单独扩缩容。

### 分布式

分布式更宽泛，指的是系统不只跑在一台机器或一个进程里，而是分散在多个节点上协作。

比如：

- 多台应用服务器
- 独立 Redis
- 独立 MySQL
- 消息队列
- 网关

### 为什么你现在不应该一上来就真拆微服务

因为第一次做完整项目时，先学清楚这些更重要：

- 请求怎么进来
- 代码怎么分层
- 数据怎么流动
- 库存为什么要事务
- 缓存为什么要失效

如果这些都没掌握，直接拆微服务只会把复杂度放大。

所以现在这个结构是“单体实现，但按微服务思路分层”，是最适合练手的阶段。

## 10. 你应该怎么读这个项目

推荐顺序如下：

1. 先读 [README.md](C:\Users\HP\Desktop\go_test\README.md)
2. 再读 [main.go](C:\Users\HP\Desktop\go_test\cmd\server\main.go)
3. 再读 [server.go](C:\Users\HP\Desktop\go_test\internal\app\server.go)
4. 再读 [models.go](C:\Users\HP\Desktop\go_test\internal\models\models.go)
5. 再读 `repository`
6. 再读 `service`
7. 再读 `handler`
8. 最后读前端 `web/`

为什么这样读：

因为你先知道“怎么启动”和“有哪些模块”，再看局部实现，不容易迷路。

## 11. 你应该怎么继续练

第一轮练习建议你自己做这些改动：

1. 给商品增加分类字段
2. 增加订单状态流转，比如 `created`、`paid`、`shipped`
3. 增加管理员新增商品接口
4. 增加分页查询商品接口
5. 给 Redis 增加购物车缓存
6. 给订单增加“取消订单”
7. 给商品增加搜索功能

这些练习会强迫你同时改：

- model
- schema
- repository
- service
- handler
- frontend

这正是完整项目训练的核心。

## 12. 最后给你的一个实用建议

你以前从 C++ 阅读项目结构入手，这个习惯很好，但 Web 项目更强调“调用链”。

以后你看这类 Go 项目，不要只按目录看，要按“一个请求从哪里来，到哪里去”看。

比如你可以专门跟一条链：

- 登录链路
- 商品查询链路
- 下单链路

每跟完一条链，你就会比单纯看零散文件进步更快。

如果你愿意，我下一步可以继续给你做两件事里的任意一种：

1. 我按“注册 / 登录 / 下单”三条请求链路，逐文件带你精读这个项目。
2. 我把这个项目继续升级成“适合你练微服务”的第二版，并告诉你下一步该怎么拆。

## 13. 再往底层看一层：一个 HTTP 请求到底经历了什么

如果你以前更熟 C++，那你现在可以把 Web 项目理解成：

- 操作系统先接收到网络数据
- Go 运行时里的 HTTP 服务器负责监听端口
- Gin 在标准库 `net/http` 之上做了更好用的封装
- 你的 handler 最终被调用

在这个项目里，请求大致是这样流动的：

1. 浏览器向 `http://localhost:8080/api/products` 发请求
2. Go 服务监听 `8080` 端口
3. Gin 路由匹配到 `GET /api/products`
4. 先执行中间件
5. 再执行 `CatalogHandler.ListProducts`
6. handler 调 service
7. service 查 Redis 或 MySQL
8. 结果一路返回给浏览器

你可以把 Gin 理解成一层“请求分发器 + 中间件执行器 + 参数解析器”。

### 为什么 handler 不直接写 SQL

因为一个 HTTP 请求里通常有三类逻辑：

- 协议逻辑
  比如 JSON 参数怎么解析、HTTP 状态码怎么返回
- 业务逻辑
  比如库存不够该不该报错
- 存储逻辑
  比如 SQL 怎么写

把这三件事搅在一起，后面就很难维护。

所以这个项目故意拆成：

- handler 管协议
- service 管业务
- repository 管存储

这是你以后看大多数 Go 后端项目时最常见的套路之一。

## 14. Gin 中间件到底是什么

中间件本质上就是：

“在真正业务处理前后插一段公共逻辑”

在这个项目里用了两个中间件：

- 跨域中间件
- 限流中间件

还有一个鉴权中间件挂在受保护路由上。

你可以把它想成洋葱模型：

1. 请求先进最外层中间件
2. 再往里进下一个中间件
3. 最后到 handler
4. handler 返回后，再一层层退出

### 为什么中间件适合做鉴权和限流

因为它们和具体业务无关，但很多接口都要用。

比如：

- 登录接口可能不用鉴权
- 下单接口必须鉴权
- 所有接口都可能需要限流

所以中间件是最合适的位置。

## 15. Go 的接口、组合、依赖注入，先建立直觉

这个项目目前还没有大量使用接口抽象，但你很快会在真实项目里遇到。

### 什么是依赖注入

比如这段：

```go
authService := service.NewAuthService(userRepo, cfg.JWTSecret)
```

这里的意思是：

- `AuthService` 自己不去创建 `userRepo`
- 而是由外部把依赖传给它

这就叫依赖注入。

为什么这么做：

- 模块之间更解耦
- 更好测试
- 更容易替换实现

以后如果你要把 `UserRepository` 换成别的实现，这种写法就很好改。

### Go 为什么常强调组合而不是继承

Go 没有传统面向对象里那种复杂继承体系，它更喜欢：

- 结构体组合
- 接口约束行为

这会让项目结构更平，也更容易控制复杂度。

## 16. 数据库事务、锁、并发下单，为什么这是重点

电商场景最典型的风险就是超卖。

比如：

- 商品库存只有 1
- 两个用户同时下单

如果你只是先查库存，再单独更新库存，中间没有锁和事务，就可能两个都买成功。

### 这个项目是怎么避免明显超卖的

看 [order_service.go](C:\Users\HP\Desktop\go_test\internal\service\order_service.go) 和 [product_repository.go](C:\Users\HP\Desktop\go_test\internal\repository\product_repository.go)。

关键点是：

1. 开事务
2. `SELECT ... FOR UPDATE`
3. 检查库存
4. 扣库存
5. 写订单
6. 提交事务

### `FOR UPDATE` 是什么

它会在事务中把查询到的行锁住。

简单理解：

- 当前事务没提交之前
- 别的事务不能同时改这行关键数据

所以它特别适合库存这类场景。

### 为什么先 `defer tx.Rollback()`，最后再 `Commit()`

这是 Go 很典型的防御式写法。

因为事务中途任意一步失败，都应该回滚。

先写 `defer tx.Rollback()` 可以避免你漏掉回滚逻辑。

## 17. Redis 缓存为什么不是“加了就一定更快”

很多初学者会把 Redis 理解成“万能加速器”，这不准确。

Redis 适合解决的是：

- 热点数据频繁读取
- 数据允许短时间不绝对实时
- 数据查询成本较高

在这个项目里，商品列表适合缓存，因为：

- 会被很多人频繁访问
- 读取比写入多
- 商品列表允许短时间缓存

### 缓存带来的问题是什么

最经典的问题叫缓存一致性。

比如：

1. 商品库存被更新了
2. Redis 里还是旧数据
3. 用户读到过期库存

所以这里用了最简单的策略：

- 查商品时先查 Redis
- 下单成功后删掉商品列表缓存

这叫“更新数据库后删除缓存”。

### 为什么不是先删缓存再更新数据库

因为那样中间有窗口期，别的请求可能又把旧数据写回缓存。

当然，真实大项目里还会有更复杂的一致性方案，但你现在先记住这条就够：

优先考虑“先更新数据库，再删缓存”。

## 18. JWT 的底层理解

JWT 不是 session 存储本身，它更像一个“自带信息、可验签”的令牌。

它通常有三部分：

- Header
- Payload
- Signature

### 这个项目里 JWT 在做什么

登录成功后：

- 服务端生成 token
- token 里放用户 ID 和过期时间
- 用密钥签名

后续请求时：

- 前端把 token 放到 `Authorization` 头里
- 后端中间件验证签名
- 验签通过就知道是谁在请求

### JWT 的优点

- 服务端不用像传统 session 那样强依赖内存会话表
- 比较适合前后端分离

### JWT 的风险

- 一旦签发出去，在过期前通常都有效
- 如果泄露，别人就能冒充用户

所以生产环境里你还会考虑：

- 更短过期时间
- refresh token
- 黑名单
- HTTPS

## 19. 高并发和“能承受很多用户”到底意味着什么

你提到希望项目“也能承受很多用户”，这句话在工程上要拆开看。

真正支撑高并发不是只靠某一个框架，而是靠整套设计：

- 应用可水平扩容
- Redis 扛热点读
- 数据库索引合理
- 事务范围尽量小
- 慢查询少
- 限流保护系统
- 网关和负载均衡
- 异步任务削峰

所以“用了 Gin”不等于高并发，“用了 Redis”也不等于高并发。

你现在这个项目只是打下了几个基础点：

- 商品缓存
- 简单限流
- 事务下单
- 分层结构

这些是以后继续往高并发演进的基础。

## 20. 如果以后真的拆微服务，这个项目会怎么拆

你可以先有个脑图，但先别急着真拆。

比较自然的拆法是：

- 用户服务
  负责注册、登录、用户资料
- 商品服务
  负责商品列表、库存、商品管理
- 订单服务
  负责下单、订单查询
- API 网关
  对外统一入口

然后再引入：

- Redis 做缓存和分布式锁
- MQ 做异步下单、通知、削峰
- Nacos/Consul 做服务发现
- Prometheus/Grafana 做监控

但在你当前阶段，更重要的是先把单体的调用链吃透。

## 21. 面试八股题：先记和这个项目强相关的

下面这些不是死记硬背题，而是你看着当前项目最容易真正理解的题。

### 1. Gin 和标准库 `net/http` 的关系是什么

答：

Gin 不是替代 HTTP 协议，而是构建在 Go 标准库 `net/http` 之上的 Web 框架。它在标准库基础上封装了路由、中间件、参数绑定、JSON 返回等能力，让开发 API 更高效。

### 2. 为什么要做项目分层

答：

分层是为了隔离职责。handler 管 HTTP，service 管业务，repository 管数据库。这样修改业务逻辑时不容易影响协议层，修改 SQL 时也不会把 handler 搞乱，项目更容易维护和测试。

### 3. 为什么下单必须用事务

答：

因为下单涉及多个步骤，例如查库存、扣库存、写订单。这些步骤必须保持原子性，要么全部成功，要么全部失败，否则会出现库存扣了但订单没写成功，或者订单成功了但库存没扣的问题。

### 4. `SELECT ... FOR UPDATE` 的作用是什么

答：

它会在事务中对选中的行加排他锁，防止其他事务同时修改这行数据。库存扣减场景里常用它避免并发下单导致超卖。

### 5. Redis 为什么适合做缓存

答：

因为 Redis 基于内存，读写速度很快，适合存放热点数据。像商品列表这种读多写少的数据，放到 Redis 可以减少 MySQL 压力，提高响应速度。

### 6. Redis 和 MySQL 的区别是什么

答：

MySQL 是关系型数据库，适合持久化存储和复杂查询；Redis 是内存型键值存储，适合缓存、计数、会话、排行榜等高频访问场景。真实项目里两者常常配合使用，而不是互相替代。

### 7. 为什么是“更新数据库后删除缓存”

答：

因为如果先删缓存再更新数据库，中间可能有并发请求把旧数据重新写回缓存，导致脏数据。先更新数据库再删缓存，通常更安全，是常见的缓存一致性策略。

### 8. JWT 为什么能做登录鉴权

答：

因为 JWT 中包含用户标识和过期时间，并经过服务端密钥签名。服务端收到 token 后可以校验签名是否合法，从而判断用户身份是否可信。

### 9. JWT 和 Session 的区别是什么

答：

Session 通常把会话状态存在服务端，客户端只拿一个 session ID；JWT 则把一部分身份信息放进 token 里，服务端通过签名验证身份。JWT 更适合前后端分离和多服务场景，但撤销控制更复杂。

### 10. 中间件适合做什么

答：

适合做和具体业务无关、但多个接口都会用到的公共逻辑，例如鉴权、限流、日志、跨域、统一异常处理等。

### 11. 为什么说 Redis 限流是保护系统

答：

限流可以防止某个用户或某类请求在短时间内大量打进系统，避免服务、数据库或下游组件被压垮。它本质上是一种系统保护手段。

### 12. Go 为什么大量使用 `if err != nil`

答：

因为 Go 采用显式错误处理，不依赖大量异常机制。这样每一步出错都能在代码里清晰看到，流程更直接，更适合工程维护。

### 13. 为什么 service 层不应该直接处理 HTTP 细节

答：

因为 service 代表业务逻辑，它应该尽量独立于传输协议。这样同一套业务逻辑未来既能被 HTTP 调用，也能被 gRPC、消息队列消费端、定时任务复用。

### 14. 单体和微服务的区别是什么

答：

单体是把多个业务模块放在一个应用里部署；微服务是把不同业务能力拆成独立服务分别部署。单体开发和调试更简单，微服务扩展性更强但复杂度更高。

### 15. 微服务一定比单体好吗

答：

不一定。微服务适合业务复杂、团队规模较大、模块边界清晰的系统。对于早期项目或练手项目，单体通常更容易交付和理解。

## 22. 这些题怎么背更牢

不要把它们当成纯背诵题，建议你用“项目映射法”记忆。

比如：

- 看到事务，就立刻想到 [order_service.go](C:\Users\HP\Desktop\go_test\internal\service\order_service.go)
- 看到缓存，就立刻想到 [catalog_service.go](C:\Users\HP\Desktop\go_test\internal\service\catalog_service.go)
- 看到鉴权，就立刻想到 [auth.go](C:\Users\HP\Desktop\go_test\internal\http\middleware\auth.go)
- 看到分层，就立刻想到 `handler -> service -> repository`

这样你不是死背，而是“把题和代码绑在一起”。

## 23. 给你一套更有效的学习顺序

你现在最适合这样推进：

1. 先读这份文档，建立总图
2. 跟一遍登录链路
3. 跟一遍商品列表链路
4. 跟一遍下单链路
5. 自己改一个小功能
6. 再回来看这些八股题

这时你会发现，很多题其实不是背出来的，而是你已经能从代码里讲出来。

## 24. 你接下来最值得补的基础

如果你准备继续深入，优先补这几块：

- Go 的结构体、方法、接口、指针、错误处理
- HTTP 基础：请求头、状态码、JSON、Cookie、Authorization
- MySQL 基础：索引、事务、行锁、隔离级别
- Redis 基础：缓存、过期、计数器、分布式锁
- Web 项目调用链：路由、中间件、handler、service、repository

如果这些打牢了，再去看微服务、分布式、消息队列，就不会只停留在名词层面。
