# 计划管理

## 当前激活计划

- [0001-foundation.md](/Users/liangluo/code/weouc2026/docs/exec-plans/active/0001-foundation.md)
- [0002-m1-m2-closure.md](/Users/liangluo/code/weouc2026/docs/exec-plans/active/0002-m1-m2-closure.md)
- [0004-storage-state-refactor.md](/Users/liangluo/code/weouc2026/docs/exec-plans/active/0004-storage-state-refactor.md)

## 已归档计划

- [0003-owner-viewer-differentiation.md](/Users/liangluo/code/weouc2026/docs/exec-plans/completed/0003-owner-viewer-differentiation.md)

## 当前产品约束

- 学生端当前主线为跑腿、组局、二手交易、资料、失物招领
- 默认以微信小程序作为首发与主登录入口
- 联系方式查看必须完成教务绑定

## 最近更新

- `2026-05-13`：`api-server` 已继续推进 `0004-storage-state-refactor` 的第二阶段实现：运行时默认依赖已切到 `MySQL + MongoDB + Redis`，`IAM` 改为 `GORM + MySQL + Redis` 主链路并接入 `AutoMigrate`，`portal/notification/analytics/campus_life` 已补 `MongoDB` 仓储装配入口，Docker 联调环境同步改为 `mysql + mongo + redis`；原依赖 `memory/postgres` 的接口级测试已改为显式集成测试跳过态，等待真实环境用例回补。
- `2026-05-13`：新增 `0004-storage-state-refactor` 计划，并完成后端“`MySQL + MongoDB + Redis + BMFS`”目标架构的设计建档：核心主数据转向 `MySQL`，社区内容/审核/审计/配置转向 `MongoDB`，状态流转由 `BMFS` 统一承载，审核态并入单一 `status`；当前仅完成设计与文档同步，且不再要求迁移历史数据、兼容旧链路或保留历史代码逻辑。
- `2026-05-13`：已为 `0004-storage-state-refactor` 补充后端安全策略：明确 `MySQL` 参数化查询、`MongoDB` 白名单查询 builder、输入校验、限流、防越权、错误脱敏与注入回归测试要求。
- `2026-05-13`：已补充“`MySQL` 侧统一使用 `GORM`”到架构与执行计划，并在 `api-server` 中落地 `MySQL + GORM` 配置与客户端骨架，作为后续 `iam` 迁移起点。
- `2026-05-13`：`api-server` 已继续推进 `0004-storage-state-refactor` 的第一阶段实现：补齐 `mysql_redis` 的 IAM 装配路径、`GORM` 版 MySQL 用户仓储、MySQL 健康探针与对应测试；当前 MySQL 迁移脚本与真实联调仍待下一轮补齐。
- `2026-05-13`：执行 `0003-owner-viewer-differentiation` 计划进展：后端已为所有六个校园生活类型补齐 `user_role`/`is_owner`/`can_edit`/`can_delete` 及特有 `can_xxx` 字段（详情+列表），新增 Market/Resource/LostFound/Carpool 下架接口和 LostFound 标记已找到接口，小程序 Errand/Market/LostFound/Meetup 详情页主态/客态 UI 改造完成，OpenAPI 契约为 6 个类型定义了详情响应 schema 和新增接口路径。剩余 A6（编辑接口）、A8（Carpool 加入拼车）、B4（SDK 生成）、C4/C6（Resource/Carpool 前端）待后续迭代。
- `2026-05-13`：新增 `0003-owner-viewer-differentiation` 计划，统一为拼车、闲置、组局、失物招领、跑腿五个校园生活模块落地主态/客态操作行为差异化，后端返回完整 `user_role` + `can_xxx` 布尔字段，客户端直接消费不做本地推断。
- `2026-05-12`：`apps/mobile-flutter` 已完成最小可运行外壳初始化，落地 `feature-first + MVVM + repository` 基线目录、`Riverpod + GoRouter + Dio` 接线，以及最小 `analyze/test` 校验。
- `2026-05-12`：已补齐 `contracts -> SDK` 最小闭环入口：新增 `make generate-sdk` 与 `packages/contracts/scripts/generate-sdks.sh`，明确 JS/Dart SDK 生成目录和步骤。
- `2026-05-12`：管理员后台 IAM 用户/角色/权限页已收口为真实会话快照视图，不再维护前端可编辑样例数据，并移除未落地后端接口的误导性调用。
- `2026-05-12`：已补齐仓库级最小 CI（后端 `go test`、后台 `typecheck+build`、小程序语法校验）并新增统一本地校验入口 `make check` / `make ci-check`。
- `2026-05-12`：微信小程序已收口门户与通知主线：首页接入门户公告卡片与通知未读角标，消息中心接入真实通知列表与已读回执，`dataCenter` 教务读取路径切到 `/academic/*` 契约。
- `2026-05-12`：已修复 `iam/repo` 的既有测试失败，`services/api-server` 全量 `go test ./...` 恢复通过。
- `2026-05-12`：管理员后台 `admin-web` 已将跑腿、组局、二手、资料、失物招领五个校园生活管理页切到真实审核分页数据，并补齐详情抽屉与上下线/重新发布操作。
- `2026-05-12`：管理员后台登录响应已补齐真实 `permissions`，前端路由守卫与侧边栏已改为按后端权限码控制，不再在登录页硬编码权限集合。
- `2026-05-11`：完成当前未收口事项盘点，新增 `0002-m1-m2-closure` 计划；当前优先级调整为 `P0` 聚焦契约闭环、管理员后台真实联调、小程序门户/通知/消息中心接入与文档状态修正，`P2` 明确保留 `Flutter App` 外壳与 `academic` 真实连接器演进。
- `2026-05-09`：微信小程序进入壳层重构阶段，目录职责收敛为 `api -> services -> stores -> pages`，并同步补齐依赖接口文档。
- `2026-05-10`：后端 `api-server` 完成基础工程初始化，已落地服务启动骨架、统一鉴权上下文、系统级接口与 OpenAPI 基线。
- `2026-05-10`：后端继续落地小程序核心接口，已支持微信登录、教务绑定、首页动态、二手/跑腿/资料/失物招领基础接口，并接入内存种子数据与 mock provider。
- `2026-05-10`：后端补齐 Docker 运行面，已提供 `api-server + postgres + redis` 的 compose 编排，并把 `/readyz` 改为真实依赖探测。
- `2026-05-10`：后端已把 IAM 主链路接入 `PostgreSQL + Redis`，支持自动迁移、持久化用户资料、登录会话与教务验证码。
- `2026-05-10`：微信小程序已完成阶段 A API 对接收口，补齐二手/跑腿/资料/失物招领的详情与发布链路、统一错误提示，并明确上传接口未部署时的前端降级边界。
- `2026-05-10`：后端已接入腾讯 COS 文件管理，提供 STS 临时凭证、下载预签名、对象存储健康探针，并把小程序上传链路改为“只存对象路径、读取时签 URL”。
- `2026-05-10`：后端已将 `campus_life` 切换为可配置的 `memory / postgres` 双后端，现有二手、跑腿、资料、失物招领数据可落入 `PostgreSQL`。
- `2026-05-11`：后端已补齐 `carpool` 拼车基础能力，支持列表、详情、发布，并接入首页动态聚合与小程序真实发布字段。
- `2026-05-11`：后端已为二手、跑腿、资料、失物招领、拼车统一补齐 `review_status` 和最小审核接口；新发布内容默认进入待审，公开列表仅展示已发布内容。
- `2026-05-11`：已补齐 `meetup/组局` 的 OpenAPI 契约、小程序接口文档与后端模块说明，并将审核接口 `content_type` 正式扩展到 `meetup`，消除代码与契约漂移。
- `2026-05-11`：后端已补齐 `portal` 与 `notification` 两个基础域，支持门户首页/公告、站内通知列表/未读数/已读回执，以及对应管理发布接口，`M1` 的门户与通知通路具备联调条件。
- `2026-05-11`：后端已补齐 `analytics` 基础域与统一审计存储，支持后台查看审计日志与基础看板，并已接入登录、教务绑定、内容发布、通知发布、公告发布、审核更新等关键动作留痕。
- `2026-05-11`：后端已补齐 `academic` 基础域，支持已绑定学生查询学期、课程表、考试安排、成绩单，并将教务查询成功动作接入统一审计记录。
- `2026-05-11`：后端已将 `portal`、`notification`、`analytics` 切换为可配置的 `memory / postgres` 双后端，并补齐门户/通知种子数据与审计日志持久化能力。
- `2026-05-11`：管理员后台 `admin-web` 已完成工程骨架与主线真实联调；登录认证、IAM、文章/公告、校园生活管理、待审核、审核历史、数据看板、审计日志已可联调，文章管理仍保持预留入口。

## 规则

- 所有跨模块工作都必须落到计划文件
- 计划状态变化时同步更新
- 完成的计划移入 `docs/exec-plans/completed/`
- 计划应说明目标、工作流、风险、退出条件
