# api-server Docker 化运行面 PRD

## 背景

`services/api-server` 当前已经具备小程序主链路核心接口，但本地运行仍依赖手工启动服务，且 `/readyz` 只是静态占位，无法反映 `PostgreSQL`、`Redis` 等中间件的真实可用性，不利于团队联调和后续部署收敛。

## 目标

在不改变现有业务接口语义的前提下，为 `api-server` 增加可直接落地的容器化运行面：

- 提供 `api-server` 镜像构建能力
- 提供 `docker compose` 编排，统一拉起 `api-server`、`postgres`、`redis`
- 将 `/readyz` 改为真实依赖探测，反映中间件是否可用
- 补齐运行配置模板和中文运维说明，降低后续接手成本

## 非目标

- 不在本轮把业务仓储从内存实现迁移到真实 `PostgreSQL`
- 不接入真实对象存储、消息队列或生产级监控栈
- 不覆盖生产 `k8s`、灰度发布、备份恢复等完整运维体系

## 设计约束

- 业务数据访问仍保持现有模块边界，不为 Docker 化绕过 `repo` 层
- `/readyz` 要能表达“依赖未就绪”并映射到正确 HTTP 状态码
- 配置必须通过环境变量注入，不能把敏感信息硬编码进仓库
- 文档、计划、后端说明需要同步更新，保持事实来源一致

## 本轮范围

### 1. 运行时配置

- 新增 `PostgreSQL` 与 `Redis` 依赖探测配置
- 为依赖启用状态、地址、库名、超时等参数提供环境变量

### 2. 系统就绪检查

- `/readyz` 支持真实探测 `postgres`、`redis`
- 依赖不可用时返回 `503`
- 保留对象存储为可选占位状态，避免假装已接入

### 3. Docker 编排

- `services/api-server/Dockerfile`
- `ops/docker/api-server/compose.yaml`
- `ops/docker/api-server/.env.example`

### 4. 文档与校验

- 更新 `services/api-server/README.md`
- 更新 `ops/README.md`
- 更新 `docs/PLANS.md` 与执行计划
- 增加针对 readiness 聚合与 `503` 行为的测试

## 验收标准

- 执行 `docker compose -f ops/docker/api-server/compose.yaml up --build` 可以启动 `api-server`、`postgres`、`redis`
- `GET /readyz` 在依赖健康时返回 `200`，依赖故障时返回 `503`
- 本地 `go test ./...` 通过
- 运行文档能说明启动、配置和健康检查方式

## 风险

- 当前业务仓储仍为内存实现，容器化后不会自动获得持久化业务数据
- 对象存储未接入，涉及文件上传的完整链路仍不是本轮目标
- 本地 Docker 版本差异可能导致 `depends_on.condition` 支持行为不同，需要文档明确要求使用新版 `docker compose`
