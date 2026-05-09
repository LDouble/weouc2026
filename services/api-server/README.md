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
- `internal/platform`：配置、日志、请求 ID、统一错误响应、鉴权上下文
- `internal/modules/system`：`/healthz`、`/readyz`、`/api/v1/system/profile`
- `packages/contracts/openapi/api-server.yaml`：当前基础契约源文件

## 本地运行

```bash
cd services/api-server
go run ./cmd/api-server
```

默认监听地址为 `0.0.0.0:8080`。

常用环境变量：

- `API_SERVER_ENV`
- `API_SERVER_PORT`
- `API_SERVER_VERSION`
- `API_SERVER_AUTH_USER_ID_HEADER`
- `API_SERVER_AUTH_ROLES_HEADER`
- `API_SERVER_AUTH_PERMISSIONS_HEADER`
- `API_SERVER_AUTH_ACADEMIC_BOUND_HEADER`

## 本地校验

```bash
cd services/api-server
go test ./...
```

## 调试接口

```bash
curl http://localhost:8080/healthz
curl http://localhost:8080/readyz
curl \
  -H 'X-User-ID: student-001' \
  -H 'X-User-Roles: student' \
  -H 'X-User-Permissions: contact:view' \
  -H 'X-Academic-Bound: true' \
  http://localhost:8080/api/v1/system/profile
```
