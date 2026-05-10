# api-server Docker 联调

本目录用于本地拉起 `api-server`、`PostgreSQL`、`Redis`。

## 目录说明

- `compose.yaml`：容器编排入口
- `.env.example`：环境变量示例

## 启动方式

```bash
cd /Users/liangluo/code/weouc2026
docker compose -f ops/docker/api-server/compose.yaml up --build
```

如需自定义端口、密码或版本号，可先复制一份环境变量文件：

```bash
cd /Users/liangluo/code/weouc2026/ops/docker/api-server
cp .env.example .env
```

然后再回到仓库根目录执行：

```bash
cd /Users/liangluo/code/weouc2026
docker compose --env-file ops/docker/api-server/.env -f ops/docker/api-server/compose.yaml up --build
```

## 健康检查

- `http://localhost:8080/healthz`：进程存活
- `http://localhost:8080/readyz`：依赖就绪；当启用的 `postgres`、`redis` 或 `object_storage` 不可用时返回 `503`

## 当前默认行为

- compose 默认把 `API_SERVER_IAM_BACKEND` 设为 `postgres_redis`
- compose 默认开启 `API_SERVER_AUTO_MIGRATE=true`，服务启动时会自动执行内置 IAM 迁移
- 当前只有 IAM 状态会持久化到 `PostgreSQL + Redis`；`campus_life` 业务数据仍为内存种子数据
- COS 默认为关闭；只有补齐 `API_SERVER_COS_*` 环境变量后，文件上传链路才会启用

## COS 配置

如需启用腾讯 COS，请至少补齐以下变量：

- `API_SERVER_COS_ENABLED=true`
- `API_SERVER_COS_SECRET_ID`
- `API_SERVER_COS_SECRET_KEY`
- `API_SERVER_COS_BUCKET`
- `API_SERVER_COS_REGION`

可选项：

- `API_SERVER_COS_PATH_PREFIX`
- `API_SERVER_COS_STS_DURATION`
- `API_SERVER_COS_PRESIGNED_GET_TTL`
- `API_SERVER_COS_HEALTHCHECK_TIMEOUT`

## 当前限制

- 只有 IAM 状态具备持久化能力，重启 `api-server` 后 `campus_life` 运行期发布数据仍会丢失
- 当前未落文件元数据独立表，文件管理仍以“业务记录引用 COS 对象路径”为主
