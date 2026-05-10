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
- `http://localhost:8080/readyz`：依赖就绪；当 `postgres` 或 `redis` 不可用时返回 `503`

## 当前默认行为

- compose 默认把 `API_SERVER_IAM_BACKEND` 设为 `postgres_redis`
- compose 默认开启 `API_SERVER_AUTO_MIGRATE=true`，服务启动时会自动执行内置 IAM 迁移
- 当前只有 IAM 状态会持久化到 `PostgreSQL + Redis`；`campus_life` 业务数据仍为内存种子数据

## 当前限制

- 只有 IAM 状态具备持久化能力，重启 `api-server` 后 `campus_life` 运行期发布数据仍会丢失
- 当前未接入真实对象存储，因此文件上传链路仍未闭环
