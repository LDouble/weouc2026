# campus\_life 统一集合设计 Spec

## Why

当前 campus\_life 模块为六种社区内容类型各建一个 MongoDB 集合（`campus_life_markets`、`campus_life_errands` 等），导致：

- 六套平行的 repo 方法、service 逻辑、handler 路由，大量重复代码
- 状态字段割裂（`review_status` + `status` 双字段），需要 `mergeErrandStatus` / `mergeMeetupStatus` 等合并函数
- 审计字段不统一（部分缺少 `updated_at` / `updated_by` / `deleted_at`）
- 新增社区类型需要复制整套集合+代码结构

以聚合根组织集合、收敛为"公共字段 + 类型载荷"的统一 `community_content` 集合，可消除上述问题。

## What Changes

- **BREAKING**：将 6 个业务集合合并为 1 个 `community_content` 集合 + 辅助集合
- **BREAKING**：统一状态模型：每个对象只有 `status` 字段，`review_status` 成为 `status` 的枚举值之一
- 统一审计字段：`created_at` / `updated_at` / `created_by` / `updated_by` / `deleted_at`
- 引入 `content_type` 字段区分社区类型（`market` / `errand` / `resource` / `lost_found` / `carpool` / `meetup`）
- 类型特有字段收敛到 `type_payload` 结构化载荷，不使用 `ext_json` 黑盒
- 保留 `ext_json` 用于未来可扩展场景，但严格限制使用范围
- 重构 repo 层为统一泛型接口，消除六套平行方法
- 重构 service 层为统一内容操作 + 类型特化逻辑
- transport 层路由保持现有 API 路径不变，内部路由到统一 handler

## Impact

- Affected specs: `unified-status-design`（状态统一已完成 API 层面，本 spec 在存储层面落地）
- Affected code:
  - `services/api-server/internal/modules/campus_life/types/models.go` — 全面重写
  - `services/api-server/internal/modules/campus_life/repo/mongo.go` — 全面重写
  - `services/api-server/internal/modules/campus_life/repo/repository.go` — 全面重写
  - `services/api-server/internal/modules/campus_life/service/service.go` — 大幅重构
  - `services/api-server/internal/modules/campus_life/transport/handler.go` — 适配新接口
  - `services/api-server/internal/modules/campus_life/transport/routes.go` — 小幅调整
  - `packages/contracts/openapi/api-server.yaml` — 契约更新

***

## ADDED Requirements

### Requirement: 统一 community\_content 集合

系统 SHALL 将所有校园生活社区内容存储在单一 `community_content` MongoDB 集合中，以聚合根组织，不按内容类型拆分集合。

#### Scenario: 存储一条二手交易内容

- **WHEN** 用户发布一条二手交易
- **THEN** 系统将内容写入 `community_content` 集合，`content_type` 为 `"market"`，类型特有字段存入 `type_payload`

#### Scenario: 存储一条跑腿任务

- **WHEN** 用户发布一条跑腿任务
- **THEN** 系统将内容写入 `community_content` 集合，`content_type` 为 `"errand"`，类型特有字段存入 `type_payload`

### Requirement: 公共字段 + 类型载荷模型

系统 SHALL 采用"公共字段 + 类型载荷"的文档模型。每条文档包含所有类型共享的公共字段，以及一个 `type_payload` 字段承载类型特有数据。

#### Scenario: 公共字段定义

- **WHEN** 写入任意类型的社区内容
- **THEN** 文档 MUST 包含以下公共字段：
  - `_id`：文档唯一标识，使用 MongoDB 自动生成的 ObjectID
  - `content_type`：内容类型枚举（`market` / `errand` / `resource` / `lost_found` / `carpool` / `meetup`）
  - `title`：标题
  - `desc`：描述
  - `status`：统一状态（见状态需求）
  - `publisher_user_id`：发布者用户 ID（唯一引用，不冗余存储昵称等可变字段）
  - `contact`：联系方式
  - `images`：图片路径列表
  - `tags`：标签列表
  - `created_at`：创建时间
  - `updated_at`：最后更新时间
  - `created_by`：创建者用户 ID
  - `updated_by`：最后更新者用户 ID
  - `deleted_at`：软删除时间（null 表示未删除）
  - `type_payload`：类型特有载荷（结构化，非黑盒 JSON）
  - `ext_json`：可扩展 JSON（仅用于未来未预见的扩展，不允许滥用为万能字段）

#### Scenario: 类型载荷结构化定义

- **WHEN** 存储类型特有数据
- **THEN** `type_payload` MUST 是结构化的、有明确 schema 的对象，每种 `content_type` 对应一个 Go struct：
  - `market` → `MarketPayload`：`category`, `price`, `original_price`, `condition`, `trade_mode`
  - `errand` → `ErrandPayload`：`category`, `route_start`, `route_end`, `deadline`, `reward`, `urgent`, `acceptor_user_id`
  - `resource` → `ResourcePayload`：`category`, `course_name`, `files`, `file_size`, `file_type`, `download_url`, `likes`, `views`
  - `lost_found` → `LostFoundPayload`：`type`, `category`, `location`, `event_time`, `item_feature`
  - `carpool` → `CarpoolPayload`：`category`, `from`, `to`, `travel_at`, `type`, `seats_text`, `price`, `note`
  - `meetup` → `MeetupPayload`：`category`, `location`, `start_at`, `deadline_at`, `max_participants`, `fee_text`, `participant_user_ids`

#### Scenario: ext\_json 使用限制

- **WHEN** 需要存储未预见的扩展数据
- **THEN** `ext_json` 仅允许用于以下场景：
  - 第三方集成需要传递的非标准字段
  - A/B 测试的临时配置数据
- **THEN** `ext_json` 不允许用于：
  - 存储本应属于 `type_payload` 的业务字段
  - 存储需要查询或索引的字段
  - 存储权限或状态相关数据

### Requirement: 发布者信息实时解析

系统 SHALL 不在内容文档中冗余存储发布者昵称和首字母，而是在读取时从用户上下文实时解析。

#### Scenario: 发布者信息不冗余存储

- **WHEN** 用户发布社区内容
- **THEN** 文档只存储 `publisher_user_id`，不存储 `publisher`（昵称）和 `publisher_initial`（首字母）

#### Scenario: 读取时解析发布者信息

- **WHEN** 返回社区内容给客户端
- **THEN** service 层根据 `publisher_user_id` 从当前请求的 principal 或用户服务中解析昵称和首字母，填充到响应 payload 中

### Requirement: 统一状态模型

系统 SHALL 对每个社区内容对象只维护一个 `status` 字段。`review_status` 的语义合并到 `status` 枚举中，不再作为独立字段。

#### Scenario: 统一状态枚举定义

- **WHEN** 定义社区内容的状态
- **THEN** `status` 字段 MUST 为以下枚举值之一：

| 状态值         | 含义     | 适用类型                    |
| ----------- | ------ | ----------------------- |
| `reviewing` | 审核中    | 全部                      |
| `published` | 已发布/正常 | 全部                      |
| `rejected`  | 审核不通过  | 全部                      |
| `offline`   | 已下线    | 全部                      |
| `cancelled` | 已取消    | errand, meetup, carpool |
| `accepted`  | 已接单    | errand                  |
| `open`      | 开放中    | meetup                  |
| `full`      | 已满     | meetup                  |
| `resolved`  | 已找到    | lost\_found             |

#### Scenario: 状态流转规则

- **WHEN** 内容状态发生变更
- **THEN** 状态流转 MUST 遵循以下规则：
  - 新发布 → `reviewing`
  - 审核通过 → `published`
  - 审核不通过 → `rejected`
  - 发布者下线 → `offline`
  - 跑腿接单 → `published` → `accepted`
  - 跑腿取消接单 → `accepted` → `published`
  - 跑腿取消发布 → `published` → `cancelled`
  - 组局报名满 → `open` → `full`
  - 组局取消 → `open`/`full` → `cancelled`
  - 失物招领标记已找到 → `published` → `resolved`

#### Scenario: 状态查询简化

- **WHEN** 判断内容是否对普通用户可见
- **THEN** 只需检查 `status == "published"`，无需同时判断 `review_status` 和 `status`

### Requirement: 统一审计字段

系统 SHALL 对所有社区内容对象维护统一的审计字段。

#### Scenario: 创建内容时设置审计字段

- **WHEN** 创建新的社区内容
- **THEN** 系统自动设置 `created_at`、`created_by`，`updated_at` 等于 `created_at`，`updated_by` 等于 `created_by`，`deleted_at` 为 null

#### Scenario: 更新内容时更新审计字段

- **WHEN** 更新社区内容
- **THEN** 系统自动更新 `updated_at` 和 `updated_by`

#### Scenario: 软删除内容

- **WHEN** 删除社区内容
- **THEN** 系统设置 `deleted_at` 为当前时间，不从数据库物理删除

### Requirement: 统一 Repository 接口

系统 SHALL 提供统一的 Repository 接口，不再为每种内容类型定义一套平行方法。

#### Scenario: 统一 CRUD 方法

- **WHEN** 对社区内容进行 CRUD 操作
- **THEN** Repository 接口 MUST 提供以下统一方法：
  - `Save(ctx, item CommunityContent) (CommunityContent, error)`
  - `GetByID(ctx, id string) (CommunityContent, error)`
  - `Update(ctx, id string, mutate func(*CommunityContent) error) (CommunityContent, error)`
  - `ListByType(ctx, contentType string, filter ContentFilter) ([]CommunityContent, error)`
  - `ListForFeed(ctx, filter FeedFilter) ([]CommunityContent, error)`

#### Scenario: 按类型查询

- **WHEN** 查询特定类型的社区内容
- **THEN** 通过 `ListByType` 方法按 `content_type` 字段过滤

#### Scenario: Feed 聚合查询

- **WHEN** 查询首页动态流
- **THEN** 通过 `ListForFeed` 方法跨类型查询，按 `created_at` 降序排列

### Requirement: ID 生成策略

系统 SHALL 使用 MongoDB 自动生成的 ObjectID 作为文档唯一标识，不再使用自定义序号 ID。

#### Scenario: 新文档 ID 生成

- **WHEN** 创建新的社区内容文档
- **THEN** 由 MongoDB 驱动自动生成 `_id`（ObjectID），不使用 `campus_life_sequences` 辅助集合

#### Scenario: 旧 ID 兼容

- **WHEN** 客户端使用旧的 `{content_type}-{序号}` 格式 ID 请求内容
- **THEN** 系统通过数据迁移映射或 404 响应处理，不保证旧 ID 继续有效

### Requirement: 收藏/点赞等关联数据

- **WHEN** 存储用户与内容的交互数据（收藏、点赞）
- **THEN** 这些数据存储在 `CommunityContent` 文档内部（如 `liked_by_user_ids`），不单独建集合

### Requirement: MongoDB 索引设计

系统 SHALL 为 `community_content` 集合建立必要的索引以支持常见查询模式。

#### Scenario: 索引定义

- **WHEN** 部署 `community_content` 集合
- **THEN** MUST 建立以下索引：
  - 唯一索引：`_id`
  - 复合索引：`{ content_type: 1, status: 1, created_at: -1 }`（按类型+状态列表查询）
  - 复合索引：`{ status: 1, created_at: -1 }`（Feed 流查询）
  - 复合索引：`{ publisher_user_id: 1, created_at: -1 }`（我的发布查询）
  - 索引：`{ deleted_at: 1 }`（软删除过滤，稀疏索引）

### Requirement: API 路径兼容

系统 SHALL 保持现有 API 路径不变，内部路由到统一 handler。

#### Scenario: 现有客户端无感迁移

- **WHEN** 客户端调用 `/api/market/list` 或 `/api/errand/detail/:id` 等现有接口
- **THEN** 响应结构与现有行为一致，客户端无需修改

***

## MODIFIED Requirements

### Requirement: 审核列表查询

审核列表查询 SHALL 通过 `community_content` 集合的 `content_type` + `status` 过滤实现，不再需要跨集合合并。

- **WHEN** 管理员查询审核列表
- **THEN** 系统通过 `ListByType` 或 `ListForFeed` 方法按 `status` 过滤（`reviewing` / `rejected`），不再需要遍历六个集合

### Requirement: 联系方式可见性

联系方式可见性判断逻辑不变，但 `contact` 字段从各类型的 `Extra` 结构提升到公共字段。

- **WHEN** 判断用户是否可查看联系方式
- **THEN** 后端根据 `canViewContact`（已绑定教务 或 是发布者本人）决定是否返回 `contact` 字段值，前端不自行放开

***

## REMOVED Requirements

### Requirement: 六集合分散存储

**Reason**: 统一为 `community_content` 单集合存储，消除六套平行代码
**Migration**: 数据迁移脚本将六个集合的数据转换并写入 `community_content`，映射规则：

- `campus_life_markets` → `content_type: "market"`
- `campus_life_errands` → `content_type: "errand"`
- `campus_life_resources` → `content_type: "resource"`
- `campus_life_lost_found` → `content_type: "lost_found"`
- `campus_life_carpools` → `content_type: "carpool"`
- `campus_life_meetups` → `content_type: "meetup"`
- 各类型的 `review_status` 映射为 `status` 的对应枚举值
- 各类型的 `Extra` 字段映射为 `type_payload`

### Requirement: review\_status 独立字段

**Reason**: 合并到统一 `status` 枚举，消除双字段判断
**Migration**: `review_status` 值直接映射为 `status` 的对应值（`reviewing` → `reviewing`，`published` → `published` 等）；对于 errand/meetup，`review_status == "published"` 时取原 `status` 字段值

### Requirement: mongoEnvelope 信封结构

**Reason**: 统一集合后不再需要信封包装，文档直接扁平存储公共字段 + `type_payload`
**Migration**: 移除 `mongoEnvelope[T]`，文档结构直接映射 `CommunityContent` struct

### Requirement: campus\_life\_sequences 辅助集合

**Reason**: 改用 MongoDB 自动生成的 ObjectID，不再需要自定义序号集合
**Migration**: 删除 `campus_life_sequences` 集合，移除 `NextID` 方法
