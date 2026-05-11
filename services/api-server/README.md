# api-server

## 技术选型

- `Go`
- `Gin`
- `PostgreSQL`
- `Redis`
- `OpenAPI`

## 目标

作为统一业务入口，对三类客户端提供一致的业务语义和权限裁决。

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

## 当前已落地

- `cmd/api-server`：服务启动入口
- `internal/platform`：配置、日志、请求 ID、统一错误响应、Bearer Token / 头部双通道鉴权上下文
- `internal/modules/system`：`/healthz`、`/readyz`、`/api/v1/system/profile`，以及 `postgres/redis` 依赖就绪探测
- `internal/modules/iam`：`/api/auth/wechat/login`、`/api/student`、`/api/edu/send-captcha`，并支持 `PostgreSQL + Redis` 持久化
- `internal/modules/portal`：`/api/portal/home`、公告列表/详情，以及管理端公告发布接口
- `internal/modules/notification`：站内通知列表、未读计数、已读回执，以及管理端通知发布接口
- `internal/modules/campus_life`：`/api/feed/list`、二手/跑腿/资料/失物招领/拼车/组局列表与关键详情/交互接口，以及校园生活审核接口
- `internal/modules/file_center`：`/api/upload/cos-sts`、`/api/upload/presigned-get`
- `internal/providers/*_provider`：微信与教务 mock provider
- `packages/contracts/openapi/api-server.yaml`：当前统一契约源文件

## 当前接口能力

- 微信小程序登录和 Bearer Token 会话
- 当前用户资料获取、教务验证码发送、绑定与解绑
- 门户首页轮播、公告列表、公告详情、管理端公告发布
- 当前用户通知列表、未读数量、通知已读回执、管理端通知发布
- 首页动态流
- 二手列表、详情、收藏、发布
- 跑腿列表、详情、发布、接单、取消发布、取消接单
- 资料列表、详情、发布
- 失物招领列表、详情、发布
- 拼车列表、详情、发布
- 组局列表、详情、发布、报名、取消报名、取消组局
- 校园生活审核列表、审核状态更新（覆盖二手、跑腿、资料、失物招领、拼车、组局）
- COS 临时上传凭证与对象下载预签名

说明：

- 当前阶段 `IAM` 已支持 `PostgreSQL + Redis` 持久化；`campus_life` 已支持 `memory / postgres` 双后端
- `portal` 与 `notification` 当前先以 `memory` 内置种子数据提供联调入口，后续再切换持久化后端
- 微信登录、教务绑定通过 mock provider 模拟，不依赖真实外部系统
- `/readyz` 已接入 `postgres`、`redis`、`object_storage` 健康探测；依赖未就绪时会返回 `503`
- 当 `API_SERVER_IAM_BACKEND=postgres_redis` 时，登录会话与教务绑定资料会分别落到 `Redis` 和 `PostgreSQL`
- 当 `API_SERVER_CAMPUS_LIFE_BACKEND=postgres` 时，二手、跑腿、资料、失物招领、拼车、组局会持久化到 `PostgreSQL`
- 已支持腾讯 COS 直传；业务层只存对象路径，读取时由后端签发可访问 URL
- 新发布的校园生活内容默认进入 `reviewing`；公开列表与详情只暴露 `published` 内容，发布者和具备 `campus_life:moderate` 权限的管理员可继续查看待审内容

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
- `API_SERVER_IAM_BACKEND`
- `API_SERVER_CAMPUS_LIFE_BACKEND`
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

如需本地使用持久化 IAM 和 `campus_life` 仓储，可配合 Docker 中间件启动：

```bash
cd /Users/liangluo/code/weouc2026
docker compose -f ops/docker/api-server/compose.yaml up -d postgres redis

cd /Users/liangluo/code/weouc2026/services/api-server
API_SERVER_IAM_BACKEND=postgres_redis \
API_SERVER_CAMPUS_LIFE_BACKEND=postgres \
API_SERVER_AUTO_MIGRATE=true \
API_SERVER_POSTGRES_ENABLED=true \
API_SERVER_POSTGRES_HOST=127.0.0.1 \
API_SERVER_POSTGRES_PORT=5432 \
API_SERVER_POSTGRES_DATABASE=weouc \
API_SERVER_POSTGRES_USER=weouc \
API_SERVER_POSTGRES_PASSWORD=weouc \
API_SERVER_REDIS_ENABLED=true \
API_SERVER_REDIS_HOST=127.0.0.1 \
API_SERVER_REDIS_PORT=6379 \
go run ./cmd/api-server
```

如需本地启用 COS 文件管理，可继续补充：

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

- 若使用 `memory` 后端，可直接访问内置示例详情如 `market-101`
- 若使用 `postgres` 后端，首次启动默认无演示数据，需先发布后再访问详情

```bash
curl http://localhost:8080/api/market/list
curl http://localhost:8080/api/carpool/list
curl http://localhost:8080/api/meetup/list
curl http://localhost:8080/api/portal/home
curl http://localhost:8080/api/market/detail/market-101
curl http://localhost:8080/api/student -H "Authorization: Bearer <token>"
curl http://localhost:8080/api/notification/list -H "Authorization: Bearer <token>"
curl -X POST http://localhost:8080/api/edu/send-captcha \
  -H "Authorization: Bearer <token>" \
  -H 'Content-Type: application/json' \
  -d '{"sid":"20260001"}'
curl -X POST http://localhost:8080/api/student \
  -H "Authorization: Bearer <token>" \
  -H 'Content-Type: application/json' \
  -d '{"student_id":"20260001","password":"password-001","captcha":"123456"}'
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
curl http://localhost:8080/api/admin/campus-life/review/list \
  -H 'X-User-ID: admin-001' \
  -H 'X-User-Permissions: campus_life:moderate'
curl -X POST http://localhost:8080/api/admin/campus-life/review/update \
  -H 'X-User-ID: admin-001' \
  -H 'X-User-Permissions: campus_life:moderate' \
  -H 'Content-Type: application/json' \
  -d '{"content_type":"meetup","content_id":"meetup-101","review_status":"published"}'
```
