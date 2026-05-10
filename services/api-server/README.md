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
- `internal/modules/iam`：`/api/auth/wechat/login`、`/api/student`、`/api/edu/send-captcha`
- `internal/modules/campus_life`：`/api/feed/list`、二手/跑腿/资料/失物招领列表与关键详情/交互接口
- `internal/providers/*_provider`：微信与教务 mock provider
- `packages/contracts/openapi/api-server.yaml`：当前统一契约源文件

## 当前接口能力

- 微信小程序登录和 Bearer Token 会话
- 当前用户资料获取、教务验证码发送、绑定与解绑
- 首页动态流
- 二手列表、详情、收藏、发布
- 跑腿列表、详情、发布、接单、取消发布、取消接单
- 资料列表、详情、发布
- 失物招领列表、详情、发布

说明：

- 当前阶段使用内存仓储和种子数据，服务重启后会话与运行期发布数据会丢失
- 微信登录、教务绑定通过 mock provider 模拟，不依赖真实外部系统
- `/readyz` 已接入 `postgres`、`redis` 健康探测；依赖未就绪时会返回 `503`
- `upload/cos-sts` 与真实对象存储直传尚未实现，因此“带真实文件上传”的完整发布联调仍有缺口

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

## 本地校验

```bash
cd services/api-server
go test ./...
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

```bash
curl http://localhost:8080/api/market/list
curl http://localhost:8080/api/market/detail/market-101
curl http://localhost:8080/api/student -H "Authorization: Bearer <token>"
curl -X POST http://localhost:8080/api/edu/send-captcha \
  -H "Authorization: Bearer <token>" \
  -H 'Content-Type: application/json' \
  -d '{"sid":"20260001"}'
curl -X POST http://localhost:8080/api/student \
  -H "Authorization: Bearer <token>" \
  -H 'Content-Type: application/json' \
  -d '{"student_id":"20260001","password":"password-001","captcha":"123456"}'
```
