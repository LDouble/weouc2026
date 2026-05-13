# 后端存储与状态驱动重构 PRD

## 背景

当前 `services/api-server` 的实际实现仍以 `PostgreSQL + Redis` 为主，`campus_life`、`portal`、`notification`、`analytics` 也围绕 `memory / postgres` 双后端展开。这套实现能支撑当前联调，但已经暴露出三个结构性问题：

- 核心主数据与社区内容、审计日志、运营配置混放在同一类关系型建模思路里，后续扩展成本高
- 校园生活六类内容虽然业务属性差异明显，但现有结构更多是“按功能补接口”，还不是按领域属性统一建模
- 审核、发布、关闭、报名、接单等状态规则尚未抽象成统一状态驱动机制，后续新增内容类型时会继续复制状态流转代码

用户明确提出以下重构方向：

1. 社区相关内容、审计、日志、配置等使用 `MongoDB`
2. 核心数据（如用户表）使用 `MySQL`
3. `MySQL` 访问统一使用 `GORM`
4. 使用 `BMFS` 做状态驱动
5. 不完全复用当前数据结构，而是基于架构与业务属性重新设计

因此，本次工作不是在现有表结构上继续修补，而是为后端下一阶段演进建立新的目标架构、领域模型和迁移路线。

## 目标

- 建立 `MySQL + MongoDB + Redis + BMFS` 的目标后端架构
- 把“核心主数据”和“社区内容/审计/配置”分开建模
- 基于业务属性重建社区领域模型，而不是继续按当前表结构小修小补
- 使用 `BMFS` 统一承载发布、审核、成交、报名、关闭等状态流转，审核不再单独抽 `review_status`
- 保持“业务规则只在后端 service 层、权限只在后端裁决、多端只消费后端语义”这三条原则不变

## 非目标

- 本轮不直接完成全部代码迁移
- 本轮不拆分微服务，仍保持模块化单体
- 本轮不要求客户端重写业务规则
- 本轮不引入分布式事务
- 本轮不要求保留历史数据迁移方案
- 本轮不要求兼容当前旧链路和旧 DTO
- 本轮不保留历史代码逻辑，允许直接删除旧 `PostgreSQL/review_status/legacy service` 实现

## 关键假设

- `BMFS` 为你指定的状态驱动引擎/状态机框架；仓库当前尚无现成依赖，因此本 PRD 先按“BMFS 适配层”设计接入位，具体 Go SDK 或集成方式需在实施前最终确认
- `MySQL` 仓储统一基于 `GORM` 实现，避免回到手写 SQL 拼装模式
- 现有 `Redis` 继续承担会话、验证码、热点缓存、幂等键等短期数据
- 社区领域读写以 `MongoDB` 为主事实源；核心身份、组织、角色、绑定关系等强事务主数据以 `MySQL` 为主事实源

## 设计约束

- 数据库访问必须只经过 `repo`
- 外部系统调用必须只经过 `providers`
- 业务规则只在后端 `service` 层实现
- 权限判断只在后端完成
- 接口变更必须遵循“先改契约，再改实现，再升级客户端 SDK”
- 多端必须直接消费后端返回的状态语义与 `can_xxx` 布尔值
- 实施重构时不保留 legacy 代码分支；旧状态逻辑、旧仓储实现、旧 DTO 映射在切换后可直接删除

## 目标架构

### 存储职责划分

| 领域 | 主存储 | 原因 |
| --- | --- | --- |
| `iam` 用户、账号、角色、组织、教务绑定索引 | `MySQL + GORM` | 强事务、一致性要求高、关系明确，且需要统一 ORM 与事务约束 |
| 社区内容主文档（跑腿、组局、拼车、二手、资料、失物招领） | `MongoDB` | 结构异构、字段演进快、适合文档聚合 |
| 审核单、状态流转日志、操作审计、业务日志 | `MongoDB` | 写多读多、结构附带上下文、适合追踪链路 |
| 门户公告、轮播、运营配置、字典配置、活动编排 | `MongoDB` | 配置与内容结构灵活、版本演进快 |
| 会话、验证码、热点缓存、幂等键 | `Redis` | 短时、高频、可过期 |

### 状态驱动职责

`BMFS` 负责统一定义和执行以下状态流转。每类内容只有一个权威 `status`，审核态只是 `status` 的一个阶段，不再单独维护 `review_status`：

- 跑腿：`draft -> reviewing -> open -> accepted -> delivering -> completed | rejected | cancelled | closed`
- 二手：`draft -> reviewing -> listed -> reserved -> traded | rejected | archived | closed`
- 组局/拼车：`draft -> reviewing -> open -> full -> started -> finished | rejected | cancelled | closed`
- 资料：`draft -> reviewing -> published | rejected | hidden | closed`
- 失物招领：`draft -> reviewing -> open -> matching -> claimed | rejected | archived | closed`
- 权限派生动作：`can_edit`、`can_delete`、`can_join`、`can_accept`、`can_close`、`can_view_contact`

`service` 层不再手写分散的 `if status == ...` 判断，而是通过 `BMFS` 查询：

- 当前状态
- 允许动作
- 执行动作后的目标状态
- 是否需要写入审计或派发异步事件

## 领域重建设计

### 设计原则

- 不再按“每个功能一套孤立表结构”继续生长
- 先抽公共领域属性，再承载类型特有载荷
- 公开展示、审核、详情分别建立清晰聚合和投影

### 社区内容统一聚合

建议重建统一的 `community_content` 聚合文档，核心结构按业务属性拆为以下部分：

- `identity`
  - `content_id`
  - `kind`: `errand | market | resource | lost_found | carpool | meetup`
  - `campus_id`
  - `author_user_id`
- `publication`
  - `title`
  - `summary`
  - `tags`
  - `media`
  - `visibility`
- `contact_policy`
  - `contact_visible_mode`
  - `requires_academic_binding`
  - `contact_snapshot`
- `state`
  - `machine_key`
  - `machine_version`
  - `status`
  - `status_reason`
  - `last_event`
- `payload`
  - 跑腿：费用、起终点、时限、接单约束
  - 二手：价格、成色、交易方式、库存语义
  - 资料：资料分类、附件引用、下载策略
  - 失物招领：物品特征、拾取/遗失地点、认领规则
  - 拼车：出发地、目的地、时间、座位
  - 组局：活动地点、时间、人数上限、报名策略
- `counters`
  - 浏览、收藏、报名、评论、举报等计数
- `audit`
  - `created_at`
  - `updated_at`
  - `created_by`
  - `updated_by`

### 辅助集合

- `community_actions`
  - 收藏、报名、接单、认领、举报、已读等用户动作
- `moderation_cases`
  - 审核单、审核意见、操作者、命中规则、关联内容快照
- `state_transition_logs`
  - `BMFS` 每次事件执行前后状态、命令、操作者、时间
- `feed_projections`
  - 首页流、管理端列表、搜索结果等读优化投影
- `app_configs`
  - 频道开关、审核策略、字典项、门户位配置、推荐位配置
- `audit_logs`
  - 登录、绑定、审核、发布、配置变更、敏感字段读取等操作留痕

### MySQL 核心主数据

建议把以下对象作为 `MySQL` 主事实：

- `users`
- `user_accounts`
- `user_profiles`
- `roles`
- `permissions`
- `user_role_bindings`
- `org_units`
- `academic_bindings`
- `identity_verifications`

说明：

- 联系方式是否可见的最终裁决必须以 `MySQL` 中的用户/教务绑定事实为准
- `MongoDB` 中如存在展示优化所需的冗余快照，只能作为读优化，不能替代权限事实源

## 接口与服务边界

### 后端服务层

- `service` 负责：
  - 加载当前聚合
  - 读取操作者身份与权限
  - 调用 `BMFS` 判定允许动作
  - 组织 repo 写入
  - 产生日志、审计与异步事件

- `repo` 负责：
  - `MySQL/GORM` 仓储
  - `MongoDB` 仓储
  - 必要的读模型拼装

- `runtime` 负责：
  - 跨存储异步投影
  - 状态补偿任务
  - 审计归档

### 客户端语义

- 客户端不重建状态机
- 客户端只消费后端返回的：
  - `status`
  - `permissions`
  - `can_xxx`

## 一致性策略

- 单存储内优先使用本地事务或单文档原子操作
- `MySQL` 与 `MongoDB` 之间不做分布式事务
- 跨存储同步采用 `outbox + runtime worker`
- 关键状态切换必须记录幂等键与状态迁移日志

## 安全策略

### 1. 防 SQL 注入

- `MySQL` 访问只允许出现在 `repo` 层
- 统一使用 `GORM` 参数绑定、预编译语句或受控 scope，禁止字符串拼接 SQL
- 排序、筛选、分页字段全部走后端白名单映射，不接受前端直接传列名、表达式、原始片段
- 搜索语句统一封装，禁止临时拼接 `%keyword%` 一类模式串
- 应用账号不授予 DDL、高危系统表访问等非必要权限

### 2. 防 NoSQL 注入

- `MongoDB` 查询只能由后端用强类型输入构建，禁止透传任意 JSON 查询对象
- 禁止客户端控制 `$where`、未受控 `$regex`、动态聚合表达式和任意 pipeline stage
- 查询 builder 只暴露允许的操作符与字段，字段名、排序键、投影字段全部白名单化
- 对数组筛选、标签查询、全文搜索、时间范围查询提供固定 builder，而不是透传 map

### 3. 防状态与权限越权

- 客户端只能提交动作命令，不能直接提交目标 `status`
- `status` 迁移必须经 `service + BMFS` 校验
- 权限、角色、审核、联系方式可见性、组织范围全由后端裁决
- `MySQL` 中的身份与绑定事实是唯一权限源，`MongoDB` 冗余数据不能反向覆盖权限结论

### 4. 输入校验与限流

- `transport` 层做统一 DTO 校验：长度、枚举、时间范围、文件大小、分页上限
- 登录、验证码、搜索、发布、审核等接口增加限流与防刷策略
- 高风险写操作要求幂等键、请求 ID 或防重放控制

### 5. 运行与审计

- 数据库访问设置超时、慢查询监控、连接池上限
- 错误响应禁止回显 SQL、Mongo 查询文档、内部堆栈与密钥
- 审计必须覆盖登录、绑定、权限变更、审核、配置变更、敏感字段读取
- 为 SQL/NoSQL 注入和越权补自动化测试样例，作为重构验收项

## 分阶段实施建议

### 阶段 0：设计冻结

- 完成任务、PRD、架构文档、计划文件与模块 README 同步
- 明确 `BMFS` 的具体依赖与接入方式

### 阶段 1：契约与基础设施抽象

- 先调整 `OpenAPI / contracts`
- 在 `config` 与 `repo` 层增加 `mysql`、`mongodb` 后端抽象
- 新契约直接按单一 `status` 设计，不为旧 `review_status` 做兼容层
- 从阶段 1 开始即允许删除旧契约适配逻辑，不建立 `legacy adapter`
- 同步建立统一查询白名单、输入校验和安全错误响应基线

### 阶段 2：MySQL 核心主数据迁移

- 先迁移 `iam`
- 落地 `GORM` 模型、仓储与事务边界
- 建立账号、角色、绑定与权限主数据表
- 保持 `Redis` 会话与验证码不变

### 阶段 3：MongoDB 社区内容迁移

- 重建 `community_content` 聚合和相关辅助集合
- 不直接沿用当前 `campus_life` 的表结构
- 直接按新模型实现发布、详情、列表、审核主链路，不迁移历史数据
- 新链路稳定后，直接删除旧 `postgres campus_life` 逻辑，不做长期共存
- 落地 `MongoDB` 查询 builder，禁止原始查询对象透传

### 阶段 4：BMFS 状态驱动接入

- 为六类社区内容定义状态机
- 将 `service` 中的状态流转收口到 `BMFS`
- 输出统一 `can_xxx` 语义

### 阶段 5：读模型与管理端迁移

- 迁移首页流、审核列表、审计日志、门户内容、配置管理
- 客户端升级生成 SDK，切到新契约
- 补注入与越权回归测试，作为切换前验收门槛

## 验收标准

- 架构文档、计划、模块 README 与任务/PRD 一致
- 目标存储边界明确，且与当前代码现状区分清楚
- `BMFS` 在设计层被定义为唯一状态驱动入口
- 社区模型被重建为统一聚合 + 类型载荷，而不是继续复制现有表结构
- SQL 注入、NoSQL 注入、排序字段穿透、状态越权都有明确防护策略和测试要求
- 后续实施阶段可按计划独立提交和验证

## 风险

- `BMFS` 若最终选型与本文假设不一致，状态适配层接口需要调整
- `MySQL + MongoDB` 双存储会提高运行和迁移复杂度
- 审核、门户、通知、配置若不一起梳理，会再次出现“内容在 Mongo，规则散落 elsewhere”的问题

## 本次结论

本轮先完成“设计重构”，不直接修改实现代码。审批通过后，建议按“新契约定义 -> `iam` 迁移 -> 社区模型直切重建 -> `BMFS` 接入 -> 管理端与投影迁移”的顺序执行；不额外背负历史数据迁移、旧链路兼容或历史代码保留成本。
