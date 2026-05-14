# api-server

## 技术选型

- 当前实现：`PostgreSQL + Redis`
- 目标架构：`MySQL + MongoDB + Redis + BMFS`
- `MySQL` 侧 ORM：`GORM`
- `Go`
- `Gin`
- `OpenAPI`

## 目标

作为统一业务入口，对三类客户端提供一致的业务语义和权限裁决。

## 重构设计状态

- 当前代码实现仍以 `PostgreSQL + Redis` 为主
- 已完成目标重构设计建档：核心主数据迁移到 `MySQL`，社区内容/审核/审计/配置迁移到 `MongoDB`，状态流转改由 `BMFS` 统一承载
- 设计详情见：
  - [ARCHITECTURE.md](/Users/liangluo/code/weouc2026/ARCHITECTURE.md)
  - [0004-storage-state-refactor.md](/Users/liangluo/code/weouc2026/docs/exec-plans/active/0004-storage-state-refactor.md)
  - [api-server-storage-state-refactor.md](/Users/liangluo/code/weouc2026/.agent/PRDS/api-server-storage-state-refactor.md)

## 推荐结构

```text
cmd/
internal/
├── modules/
├── providers/
└── platform/
api/
migrations/
```

## 模块内部层级

```text
types -> config -> repo -> service -> runtime -> transport
```

## 约束

- 数据访问只走 `repo`
- 外部系统只走 `providers`
- 权限与状态机只保留在后端
- `MySQL` 仓储统一使用 `GORM`

## 安全基线

- `MySQL` 只允许 `GORM` 参数绑定或受控 scope，禁止拼接 SQL
- `MongoDB` 只允许后端构造白名单查询，禁止透传原始查询对象
- 排序、筛选、分页字段必须白名单化
- 客户端只能提交动作，不能直接写 `status`、角色、权限或敏感可见性字段
- 登录、验证码、发布、审核、搜索接口需要限流、防刷和幂等控制
- 错误响应禁止泄露 SQL、Mongo 查询、堆栈和密钥

## 当前已落地

- `cmd/api-server`：服务启动入口
- `internal/platform`：配置、日志、请求 ID、统一错误响应、Bearer Token / 头部双通道鉴权上下文
- `internal/platform/bmfs`：BMFS 状态驱动引擎，声明式状态机定义、guard 守卫、onTransition 钩子、can_xxx 自动派生；6 类社区内容状态机已定义并接入 service 层
- `internal/modules/system`：`/healthz`、`/readyz`、`/api/v1/system/profile`，以及 `mysql/mongo/redis` 依赖就绪探测
- `internal/modules/iam`：`/api/auth/wechat/login`、`/api/student`、`/api/edu/send-captcha`，当前主链路为 `GORM + MySQL + Redis`
- `internal/modules/academic`：`/api/academic/semesters`、`/api/academic/schedule`、`/api/academic/exams`、`/api/academic/grades`
- `internal/modules/portal`：`/api/portal/home`、公告列表/详情、管理端轮播 CRUD、管理端公告发布/更新/删除接口，当前主链路为 `MongoDB`
- `internal/modules/notification`：站内通知列表、未读计数、已读回执，以及管理端通知发布接口，当前主链路为 `MongoDB`
- `internal/modules/analytics`：`/api/admin/analytics/dashboard`、`/api/admin/analytics/audit-logs`，以及统一审计日志读取能力，当前主链路为 `MongoDB`
- `internal/modules/campus_life`：`/api/feed/list`、二手/跑腿/资料/失物招领/拼车/组局列表与关键详情/交互接口，以及校园生活审核接口
- `internal/modules/file_center`：`/api/upload/cos-sts`、`/api/upload/presigned-get`
- `internal/providers/*_provider`：微信与教务 mock provider
- `packages/contracts/openapi/api-server.yaml`：当前统一契约源文件
- `migrations/*.sql`：保留历史 PostgreSQL 迁移文件；当前主链路的 `IAM` 建表由 `GORM AutoMigrate` 承载

## 当前接口能力

- 微信小程序登录和 Bearer Token 会话
- 当前用户资料获取、教务验证码发送、绑定与解绑
- 当前已绑定学生的学期、课程表、考试安排、成绩单查询
- 门户首页轮播、公告列表、公告详情、管理端轮播 CRUD、管理端公告发布/更新/删除
- 当前用户通知列表、未读数量、通知已读回执、管理端通知发布
- 审计日志列表、基础运营看板
- 首页动态流
- 二手列表、详情、收藏、发布
- 跑腿列表、详情、发布、接单、取消发布、取消接单
- 资料列表、详情、发布
- 失物招领列表、详情、发布
- 拼车列表、详情、发布
- 组局列表、详情、发布、报名、取消报名、取消组局
- 校园生活审核列表、审核状态更新（覆盖二手、跑腿、资料、失物招领、拼车、组局），审核接口已改为接收 Action 字段（`review_approve` / `review_reject`）
- COS 临时上传凭证与对象下载预签名

说明：

- 当前运行时主链路已切到 `MySQL + MongoDB + Redis`
- `IAM` 通过 `GORM + MySQL` 持久化用户主资料，`Redis` 继续承载会话与验证码
- `portal`、`notification`、`analytics` 已切到 `MongoDB`
- 微信登录、教务绑定、教务课表/考试/成绩读取当前均通过 mock provider 模拟，不依赖真实外部系统
- `/readyz` 已接入 `mysql`、`mongo`、`redis`、`object_storage` 健康探测；依赖未就绪时会返回 `503`
- 当 `API_SERVER_IAM_BACKEND=mysql_redis` 时，登录会话与验证码继续落到 `Redis`，用户主资料走 `GORM + MySQL`
- 当 `API_SERVER_CAMPUS_LIFE_BACKEND=mongo` 时，二手、跑腿、资料、失物招领、拼车、组局会持久化到 `MongoDB`
- 当 `API_SERVER_PORTAL_BACKEND=mongo` 时，门户轮播与公告会从 `MongoDB` 读取
- 当 `API_SERVER_NOTIFICATION_BACKEND=mongo` 时，通知发布与已读状态会持久化到 `MongoDB`
- 当 `API_SERVER_ANALYTICS_BACKEND=mongo` 时，统一审计日志会持久化到 `MongoDB`
- 已支持腾讯 COS 直传；业务层只存对象路径，读取时由后端签发可访问 URL
- 新发布的校园生活内容默认进入 `reviewing`；公开列表与详情只暴露 `published` 内容，发布者和具备 `campus_life:moderate` 权限的管理员可继续查看待审内容

## 目标重构方向

- `iam`：从 `PostgreSQL + Redis` 演进为 `MySQL + Redis`
- `campus_life`：从 `memory / postgres` 演进为 `MongoDB + BMFS`（**BMFS 已落地**）
- `portal`：从 `memory / postgres` 演进为 `MongoDB`
- `notification`：通知内容与投放配置演进为 `MongoDB`；会话态与热点态仍可使用 `Redis`
- `analytics`：统一审计、业务日志、状态迁移日志演进为 `MongoDB`
- 所有状态规则收口到 `BMFS`，以单一 `status` 承载审核与业务阶段，`service` 只负责编排与权限裁决（**已落地**）
- 本轮重构按新模型直切设计，不要求迁移历史数据或兼容旧链路
- 新链路替换后可直接删除旧 `postgres/review_status/legacy service` 代码，不保留并存分支

## 本地运行

```bash
cd services/api-server
go run ./cmd/api-server
```

默认监听地址为 `0.0.0.0:8080`。

## Docker 联调

```bash
cd /Users/liangluo/code/weouc2026
docker compose -f ops/docker/api-server/compose.yaml up --build
```

如需覆盖默认端口、密码或版本号，可参考 [ops/docker/api-server/.env.example](/Users/liangluo/code/weouc2026/ops/docker/api-server/.env.example)。

常用环境变量：

- `API_SERVER_ENV`
- `API_SERVER_PORT`
- `API_SERVER_VERSION`
- `API_SERVER_AUTH_USER_ID_HEADER`
- `API_SERVER_AUTH_ROLES_HEADER`
- `API_SERVER_AUTH_PERMISSIONS_HEADER`
- `API_SERVER_AUTH_ACADEMIC_BOUND_HEADER`
- `API_SERVER_AUTH_ACCESS_TOKEN_TTL`
- `API_SERVER_MYSQL_ENABLED`
- `API_SERVER_MYSQL_HOST`
- `API_SERVER_MYSQL_PORT`
- `API_SERVER_MYSQL_DATABASE`
- `API_SERVER_MYSQL_USER`
- `API_SERVER_MYSQL_PASSWORD`
- `API_SERVER_MYSQL_PARAMS`
- `API_SERVER_MYSQL_HEALTHCHECK_TIMEOUT`
- `API_SERVER_IAM_BACKEND`
- `API_SERVER_CAMPUS_LIFE_BACKEND`
- `API_SERVER_PORTAL_BACKEND`
- `API_SERVER_NOTIFICATION_BACKEND`
- `API_SERVER_ANALYTICS_BACKEND`
- `API_SERVER_AUTO_MIGRATE`
- `API_SERVER_POSTGRES_ENABLED`
- `API_SERVER_POSTGRES_HOST`
- `API_SERVER_POSTGRES_PORT`
- `API_SERVER_POSTGRES_DATABASE`
- `API_SERVER_POSTGRES_USER`
- `API_SERVER_POSTGRES_PASSWORD`
- `API_SERVER_POSTGRES_SSL_MODE`
- `API_SERVER_POSTGRES_HEALTHCHECK_TIMEOUT`
- `API_SERVER_REDIS_ENABLED`
- `API_SERVER_REDIS_HOST`
- `API_SERVER_REDIS_PORT`
- `API_SERVER_REDIS_USERNAME`
- `API_SERVER_REDIS_PASSWORD`
- `API_SERVER_REDIS_DB`
- `API_SERVER_REDIS_HEALTHCHECK_TIMEOUT`
- `API_SERVER_COS_ENABLED`
- `API_SERVER_COS_SECRET_ID`
- `API_SERVER_COS_SECRET_KEY`
- `API_SERVER_COS_BUCKET`
- `API_SERVER_COS_REGION`
- `API_SERVER_COS_PATH_PREFIX`
- `API_SERVER_COS_STS_DURATION`
- `API_SERVER_COS_PRESIGNED_GET_TTL`
- `API_SERVER_COS_HEALTHCHECK_TIMEOUT`

## 本地校验

```bash
cd services/api-server
go test ./...
```

如需本地使用当前主链路，可配合 Docker 中间件启动：

```bash
cd /Users/liangluo/code/weouc2026
docker compose -f ops/docker/api-server/compose.yaml up -d mysql mongo redis

cd /Users/liangluo/code/weouc2026/services/api-server
API_SERVER_IAM_BACKEND=mysql_redis \
API_SERVER_CAMPUS_LIFE_BACKEND=mongo \
API_SERVER_PORTAL_BACKEND=mongo \
API_SERVER_NOTIFICATION_BACKEND=mongo \
API_SERVER_ANALYTICS_BACKEND=mongo \
API_SERVER_AUTO_MIGRATE=true \
API_SERVER_MYSQL_ENABLED=true \
API_SERVER_MYSQL_HOST=127.0.0.1 \
API_SERVER_MYSQL_PORT=3306 \
API_SERVER_MYSQL_DATABASE=weouc \
API_SERVER_MYSQL_USER=weouc \
API_SERVER_MYSQL_PASSWORD=weouc \
API_SERVER_MONGO_ENABLED=true \
API_SERVER_MONGO_URI='mongodb://127.0.0.1:27017/?directConnection=true' \
API_SERVER_MONGO_DATABASE=weouc \
API_SERVER_REDIS_ENABLED=true \
API_SERVER_REDIS_HOST=127.0.0.1 \
API_SERVER_REDIS_PORT=6379 \
go run ./cmd/api-server
```

如需本地启用 COS 文件管理，可继续补充：

如需本地试验完整主链路，至少补齐：

```bash
API_SERVER_MYSQL_ENABLED=true \
API_SERVER_MYSQL_HOST=127.0.0.1 \
API_SERVER_MYSQL_PORT=3306 \
API_SERVER_MYSQL_DATABASE=weouc \
API_SERVER_MYSQL_USER=weouc \
API_SERVER_MYSQL_PASSWORD=weouc \
API_SERVER_MONGO_ENABLED=true \
API_SERVER_MONGO_URI='mongodb://127.0.0.1:27017/?directConnection=true' \
API_SERVER_MONGO_DATABASE=weouc \
API_SERVER_REDIS_ENABLED=true \
API_SERVER_REDIS_HOST=127.0.0.1 \
API_SERVER_REDIS_PORT=6379 \
API_SERVER_IAM_BACKEND=mysql_redis \
API_SERVER_CAMPUS_LIFE_BACKEND=mongo \
API_SERVER_PORTAL_BACKEND=mongo \
API_SERVER_NOTIFICATION_BACKEND=mongo \
API_SERVER_ANALYTICS_BACKEND=mongo
```

```bash
API_SERVER_COS_ENABLED=true \
API_SERVER_COS_SECRET_ID=your-secret-id \
API_SERVER_COS_SECRET_KEY=your-secret-key \
API_SERVER_COS_BUCKET=your-bucket-1250000000 \
API_SERVER_COS_REGION=ap-guangzhou \
API_SERVER_COS_PATH_PREFIX=miniapp
```

## 调试接口

```bash
curl http://localhost:8080/healthz
curl http://localhost:8080/readyz
curl -X POST http://localhost:8080/api/auth/wechat/login \
  -H 'Content-Type: application/json' \
  -d '{"code":"wx-code-001","app_id":"wx-dev-app"}'
```

拿到 `token` 后，可继续验证：

说明：

- 当前不再提供 `memory / postgres` 兼容后端
- 使用 `mongo` 主链路时，首次启动默认无演示数据，需先发布后再访问详情

```bash
curl http://localhost:8080/api/market/list
curl http://localhost:8080/api/carpool/list
curl http://localhost:8080/api/meetup/list
curl http://localhost:8080/api/portal/home
curl http://localhost:8080/api/market/detail/market-101
curl http://localhost:8080/api/student -H "Authorization: Bearer <token>"
curl http://localhost:8080/api/notification/list -H "Authorization: Bearer <token>"
curl http://localhost:8080/api/admin/analytics/dashboard \
  -H 'X-User-ID: admin-001' \
  -H 'X-User-Permissions: analytics:view'
curl -X POST http://localhost:8080/api/edu/send-captcha \
  -H "Authorization: Bearer <token>" \
  -H 'Content-Type: application/json' \
  -d '{"sid":"20260001"}'
curl -X POST http://localhost:8080/api/student \
  -H "Authorization: Bearer <token>" \
  -H 'Content-Type: application/json' \
  -d '{"student_id":"20260001","password":"password-001","captcha":"123456"}'
curl http://localhost:8080/api/academic/semesters \
  -H "Authorization: Bearer <token>"
curl http://localhost:8080/api/academic/schedule \
  -H "Authorization: Bearer <token>"
curl http://localhost:8080/api/academic/grades?semester_id=2025-2026-2 \
  -H "Authorization: Bearer <token>"
curl http://localhost:8080/api/upload/cos-sts \
  -H "Authorization: Bearer <token>"
curl -X POST http://localhost:8080/api/upload/presigned-get \
  -H "Authorization: Bearer <token>" \
  -H 'Content-Type: application/json' \
  -d '{"path":"miniapp/market/u-1/20260510/example.png"}'
curl -X POST http://localhost:8080/api/meetup/publish \
  -H "Authorization: Bearer <token>" \
  -H 'Content-Type: application/json' \
  -d '{"category":"study","title":"明晚图书馆自习搭子","desc":"一起准备高数小测","location":"图书馆四楼 A 区","start_at":"2026-05-12T19:00:00+08:00","deadline_at":"2026-05-12T18:30:00+08:00","max_participants":4,"fee_text":"免费","tags":["高数","自习"],"contact":"wx-study-001"}'
curl -X POST http://localhost:8080/api/admin/portal/notices/publish \
  -H 'X-User-ID: admin-001' \
  -H 'X-User-Permissions: portal:publish' \
  -H 'Content-Type: application/json' \
  -d '{"title":"明晚停机维护通知","summary":"今晚 23:00 至明日 01:00 维护","content":"发布与审核链路将短暂只读。","audience":"all","tags":["运维","通知"],"pinned":true}'
curl -X POST http://localhost:8080/api/admin/notification/publish \
  -H 'X-User-ID: admin-001' \
  -H 'X-User-Permissions: notification:publish' \
  -H 'Content-Type: application/json' \
  -d '{"title":"今晚 23 点起暂停内容发布","content":"发布、审核和消息写入链路将短暂进入只读维护窗口。","category":"system","target_scope":"all","action_url":"/pages/home/index"}'
curl http://localhost:8080/api/admin/analytics/audit-logs?action=campus_life.review.update \
  -H 'X-User-ID: admin-001' \
  -H 'X-User-Permissions: analytics:view'
curl http://localhost:8080/api/admin/campus-life/review/list \
  -H 'X-User-ID: admin-001' \
  -H 'X-User-Permissions: campus_life:moderate'
curl -X POST http://localhost:8080/api/admin/campus-life/review/update \
  -H 'X-User-ID: admin-001' \
  -H 'X-User-Permissions: campus_life:moderate' \
  -H 'Content-Type: application/json' \
  -d '{"content_type":"meetup","content_id":"meetup-101","action":"review_approve"}'
```
